package config

import (
	"context"
	"log"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	*oidc.IDTokenVerifier
}

// New instantiates the *Authenticator.
func ConnectOIDC() (*Authenticator, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	provider, err := oidc.NewProvider(
		context.Background(),
		os.Getenv("AUTH0_DOMAIN"),
	)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("AUTH0_CLIENT_ID")})

	return &Authenticator{
		Provider:        provider,
		IDTokenVerifier: verifier,
	}, nil
}
