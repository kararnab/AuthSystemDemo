package service

import (
	"github.com/kararnab/authdemo/pkg/iam/audit"
	"github.com/kararnab/authdemo/pkg/iam/policy"
	"github.com/kararnab/authdemo/pkg/iam/provider"
	"github.com/kararnab/authdemo/pkg/iam/session"
	"github.com/kararnab/authdemo/pkg/iam/token"
	"github.com/kararnab/authdemo/pkg/metrics"
)

// Options defines all dependencies required by the IAM service.
//
// Every dependency is injected explicitly.
// This is what makes the implementation pluggable.
type Options struct {
	Providers map[string]provider.AuthProvider

	SessionManager session.Manager
	SessionStore   session.Store

	TokenIssuer   token.Issuer
	TokenVerifier token.Verifier

	PolicyEngine policy.Engine
	AuditLogger  audit.Logger
	Metrics      metrics.IAMMetrics
}
