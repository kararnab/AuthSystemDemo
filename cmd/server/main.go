package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	//"log"

	"github.com/kararnab/authdemo/internal/api"
	"github.com/kararnab/authdemo/internal/books"
	"github.com/kararnab/authdemo/internal/users"
	"github.com/kararnab/authdemo/pkg/iam/audit/stdout"
	"github.com/kararnab/authdemo/pkg/secret_store"

	"github.com/kararnab/authdemo/pkg/iam"
	"github.com/kararnab/authdemo/pkg/iam/policy"
	"github.com/kararnab/authdemo/pkg/iam/provider"
	googleprov "github.com/kararnab/authdemo/pkg/iam/provider/google"
	internalprov "github.com/kararnab/authdemo/pkg/iam/provider/inhouse"
	"github.com/kararnab/authdemo/pkg/iam/service"
	"github.com/kararnab/authdemo/pkg/iam/session"
	"github.com/kararnab/authdemo/pkg/iam/token"
	"github.com/kararnab/authdemo/pkg/iam/token/jwt"
	"github.com/kararnab/authdemo/pkg/iam/token/keys"
	"github.com/kararnab/authdemo/pkg/iam/token/paseto"

	"github.com/kararnab/authdemo/pkg/log"
	zlog "github.com/kararnab/authdemo/pkg/log/zerolog"

	"github.com/kararnab/authdemo/pkg/metrics"
	prom "github.com/kararnab/authdemo/pkg/metrics/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

//
// ================================
// CONFIG (MVP CONSTANTS)
// ================================
//
// TODO (prod):
//   - Move to env/config files
//   - Secrets via Vault / SSM / KMS
//

const (
	googleOAuthIssuerUrl = "https://accounts.google.com"
	jwtIssuer            = "auth-monolith"
	jwtAccessTTL         = 15 * time.Minute
	sessionTTL           = 24 * time.Hour
)

func main() {
	// -------------------------------
	// Logger
	// -------------------------------
	log.Init(zlog.New())

	log.Info("starting application")

	// -------------------------------
	// Metrics
	// -------------------------------
	registry := prom.NewRegistry()
	iamMetrics := prom.NewIAMMetrics(registry)

	// -------------------------------
	// IAM + infra wiring
	// -------------------------------
	iamService, userStore, keyProvider, err := buildIAMService(iamMetrics)
	if err != nil {
		log.Error(
			"failed to start IAM",
			log.F("error", err, log.RedactNone),
		)
		os.Exit(1)
	}

	// -------------------------------
	// HTTP API
	// -------------------------------
	authHandlers := api.NewHandlers(iamService, userStore)
	bookStore := books.NewMemoryStore()
	bookHandlers := api.NewBookHandlers(bookStore)
	keyRotationHandler := api.NewKeyRotationHandler(keyProvider)
	metricsHandler := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{},
	)

	router := api.NewRouter(authHandlers, bookHandlers, keyRotationHandler, metricsHandler)

	portAddr := ":" + getPort()
	log.Info(
		"🚀 Server running",
		log.F("addr", portAddr, log.RedactNone),
	)
	if err := http.ListenAndServe(portAddr, router); err != nil {
		//log.Fatalf("server failed: %v", err)
		log.Error(
			"server failed",
			log.F("error", err, log.RedactNone),
		)
		os.Exit(1)
	}
}

//
// ================================
// IAM WIRING
// ================================
//

func buildIAMService(
	iamMetrics metrics.IAMMetrics,
) (
	iam.Service,
	internalprov.UserStore,
	*keys.MemoryProvider,
	error,
) {

	ctx := context.Background()

	// Secret Loader
	store := secret_store.BuildSecretStore()
	secretUserPassword, _ := store.Get(ctx, "SECRET_USER_PASSWORD")
	secretUserId, _ := store.Get(ctx, "SECRET_USER_ID")
	secretUserName, _ := store.Get(ctx, "SECRET_USERNAME")
	secretJWTSigningKey, _ := store.Get(ctx, "SECRET_JWT_SIGNING_KEY")
	secretPasetoSigningKey, _ := store.Get(ctx, "SECRET_PASETO_SIGNING_KEY")
	googleOAuthClientID, _ := store.Get(ctx, "GOOGLE_OAUTH_CLIENTID")

	// -------------------------------
	// User store (application-owned)
	// -------------------------------
	userStore := users.NewMemoryUserStore()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(secretUserPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := userStore.Create(ctx, &internalprov.User{
		ID:           secretUserId,
		Email:        secretUserName,
		PasswordHash: string(hash),
		Roles:        []string{policy.Admin},
	}); err != nil {
		return nil, nil, nil, err
	}

	// -------------------------------
	// Providers
	// -------------------------------
	internalProvider := internalprov.New(userStore)
	googleProvider := googleprov.New(googleOAuthClientID)
	//or oidcProvider, _ := oidcprov.New(ctx, googleOAuthIssuerUrl, googleOAuthClientID)

	providers := map[string]provider.AuthProvider{
		internalProvider.Name(): internalProvider,
		googleProvider.Name():   googleProvider,
	}

	// -------------------------------
	// Sessions
	// -------------------------------
	sessionStore := session.NewMemoryStore()
	sessionManager := session.NewManager(
		sessionStore,
		sessionTTL,
	)

	// -------------------------------
	// Tokens + keys
	// -------------------------------
	var (
		issuer      token.Issuer
		verifier    token.Verifier
		keyProvider *keys.MemoryProvider
	)

	pasetoKeyB64 := secretPasetoSigningKey

	if pasetoKeyB64 != "" {
		rawKey, err := base64.StdEncoding.DecodeString(pasetoKeyB64)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("invalid PASETO_KEY: %w", err)
		}
		if len(rawKey) != 32 {
			return nil, nil, nil, fmt.Errorf("PASETO_KEY must be 32 bytes")
		}

		keyProvider = keys.NewMemoryProvider(keys.Key{
			ID:  "paseto-1",
			Key: rawKey,
		})

		issuer, err = paseto.NewIssuer(
			keyProvider.ActiveKey().Key,
			keyProvider.ActiveKey().ID,
			jwtIssuer,
			jwtAccessTTL,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		verifier = &token.MultiVerifier{
			Verifier:    paseto.NewVerifier(jwtIssuer),
			KeyProvider: keyProvider,
		}

	} else {
		// ================================
		// JWT (fallback)
		// ================================
		signingKey := []byte(secretJWTSigningKey)

		keyProvider = keys.NewMemoryProvider(keys.Key{
			ID:  "jwt-1",
			Key: signingKey,
		})

		issuer = jwt.NewIssuer(
			keyProvider.ActiveKey().Key,
			keyProvider.ActiveKey().ID,
			jwtIssuer,
			jwtAccessTTL,
		)

		verifier = &token.MultiVerifier{
			Verifier:    jwt.NewVerifier(jwtIssuer),
			KeyProvider: keyProvider,
		}
	}

	if keyProvider == nil {
		return nil, nil, nil, fmt.Errorf("no signing key configured")
	}

	// -------------------------------
	// IAM service
	// -------------------------------
	iamService, err := service.New(service.Options{
		Providers:      providers,
		SessionManager: sessionManager,
		SessionStore:   sessionStore,
		TokenIssuer:    issuer,
		TokenVerifier:  verifier,
		PolicyEngine:   &policy.DefaultPolicy{},
		AuditLogger:    &stdout.AuditLogger{},
		Metrics:        iamMetrics,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return iamService, userStore, keyProvider, nil
}

func getPort() string {
	const defaultPort = 8080

	raw := os.Getenv("PORT")
	if raw == "" {
		return fmt.Sprintf("%d", defaultPort)
	}

	port, err := strconv.Atoi(raw)
	if err != nil {
		log.Error(
			"invalid PORT (not a number)",
			log.F("PORT", raw, log.RedactNone),
		)
		os.Exit(1)
	}

	if port < 1 || port > 65535 {
		log.Error(
			"invalid PORT (out of range)",
			log.F("PORT", raw, log.RedactNone),
		)
		os.Exit(1)
	}

	return raw
}
