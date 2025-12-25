package secret_store

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("secret not found")
)

type Store interface {
	Get(ctx context.Context, name string) (string, error)
}

func BuildSecretStore() Store {
	return NewChain(
		//NewEnvStore(),
		NewDummyStore(),
	)
}
