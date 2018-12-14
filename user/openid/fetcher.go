package openid

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"github.com/pivotalservices/ignition/user"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Fetcher retrieves the profile for a user when the auth variant is "google"
type Fetcher struct {
	Verifier Verifier
}

// Verifier takes an OpenID ID token and verifies it, returning claims
type Verifier interface {
	Verify(ctx context.Context, rawIDToken string) (*Claims, error)
}

// OIDCVerifier takes an OpenID ID token and verifies it, returning an ID Token
type OIDCVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

// Claims represent metadata from an OpenID ID token
type Claims struct {
	Sub        string `json:"sub"`
	UserName   string `json:"user_name"`
	UserID     string `json:"user_id"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

// Profile retrieves the user's profile with the given context, config, and token
func (g *Fetcher) Profile(ctx context.Context, c *oauth2.Config, t *oauth2.Token) (*user.Profile, error) {
	if g.Verifier == nil {
		return nil, errors.New("unable to verify token")
	}
	if t == nil {
		return nil, errors.New("unable to verify token")
	}

	rawIDToken, ok := t.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("profile: no id_token")
	}
	claims, err := g.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch claims")
	}

	username := claims.UserName
	if strings.TrimSpace(username) == "" {
		username = claims.Email
	}

	return &user.Profile{
		Email:       claims.Email,
		AccountName: username,
		Name:        strings.TrimSpace(fmt.Sprintf("%s %s", claims.GivenName, claims.FamilyName)),
	}, nil
}

// OIDCIDVerifier is an ID token verifier
type OIDCIDVerifier struct {
	Verifier OIDCVerifier
}

// Verify takes the given raw ID token and returns claims
func (o *OIDCIDVerifier) Verify(ctx context.Context, rawIDToken string) (*Claims, error) {
	idToken, err := o.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.Wrap(err, "unable to verify rawIDToken")
	}
	var claims Claims
	if err = idToken.Claims(&claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// NewVerifier returns a Verifier that uses a keySet fetched from the jwksURL
func NewVerifier(issuerURL string, clientID string, jwksURL string, skipTLSValidation bool) Verifier {
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipTLSValidation,
			},
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, hc)
	keySet := oidc.NewRemoteKeySet(ctx, jwksURL)
	return &OIDCIDVerifier{
		Verifier: oidc.NewVerifier(issuerURL, keySet, &oidc.Config{
			ClientID: clientID,
		}),
	}
}
