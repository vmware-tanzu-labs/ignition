package openid_test

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	oidc "github.com/coreos/go-oidc"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/internal"
	"github.com/pivotalservices/ignition/user/openid"
	"github.com/pivotalservices/ignition/user/openid/openidfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"golang.org/x/oauth2"
)

func TestFetcher(t *testing.T) {
	spec.Run(t, "Fetcher", testFetcher, spec.Report(report.Terminal{}))
}

func testFetcher(t *testing.T, when spec.G, it spec.S) {
	var f *openid.Fetcher

	it.Before(func() {
		RegisterTestingT(t)
		f = &openid.Fetcher{}
	})

	when("using the OIDCIDVerifier", func() {
		var (
			o *openid.OIDCIDVerifier
			f *openidfakes.FakeOIDCVerifier
		)

		it.Before(func() {
			f = &openidfakes.FakeOIDCVerifier{}
			o = &openid.OIDCIDVerifier{
				Verifier: f,
			}
		})

		it("returns an error when the verifier fails", func() {
			f.VerifyReturns(nil, errors.New("test error"))
			c, err := o.Verify(context.Background(), "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("test error"))
			Expect(c).To(BeNil())
		})

		it("returns an error when claims are not set", func() {
			f.VerifyReturns(&oidc.IDToken{}, nil)
			c, err := o.Verify(context.Background(), "")
			Expect(err).To(HaveOccurred())
			Expect(c).To(BeNil())
		})
	})

	when("using the real oidc.IDTokenVerifier", func() {
		var (
			o *openid.OIDCIDVerifier
			v *oidc.IDTokenVerifier
		)

		it.Before(func() {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				tokenKeys := internal.StringFromTestdata(t, "token_keys.json")
				publicKey := internal.StringFromTestdata(t, "fake-ignition-public-key.pem")
				block, _ := pem.Decode([]byte(publicKey))
				if block == nil {
					t.Fatal("failed to parse PEM block containing the public key")
				}
				pub, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					t.Fatal("failed to parse DER encoded public key: " + err.Error())
				}
				p, ok := pub.(*rsa.PublicKey)
				if !ok {
					t.Fatal("not a public key")
				}
				modulus := string(base64.RawURLEncoding.EncodeToString(p.N.Bytes()))
				jwks := strings.Replace(tokenKeys, "{{modulus}}", modulus, 1)
				io.Copy(w, bytes.NewReader([]byte(jwks)))
			}))
			config := &oidc.Config{
				ClientID:          "test-client-id",
				SkipClientIDCheck: true,
				SkipExpiryCheck:   true,
			}
			v = oidc.NewVerifier("https://ignition.uaa.run.pcfbeta.io/oauth/token", oidc.NewRemoteKeySet(context.Background(), s.URL), config)
			o = &openid.OIDCIDVerifier{
				Verifier: v,
			}
		})

		it("returns the claims if there are no errors", func() {
			c, err := o.Verify(context.Background(), "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDQiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDEiLCJzY29wZSI6WyJvcGVuaWQiLCJwcm9maWxlIiwidXNlcl9hdHRyaWJ1dGVzIl0sImNsaWVudF9pZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMiIsImNpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMiIsImF6cCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMiIsImdyYW50X3R5cGUiOiJhdXRob3JpemF0aW9uX2NvZGUiLCJ1c2VyX2lkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAxIiwib3JpZ2luIjoib2t0YS1pZGVudGl0eSIsInVzZXJfbmFtZSI6InRlc3RlckBwaXZvdGFsLmlvIiwiZW1haWwiOiJ0ZXN0ZXJAcGl2b3RhbC5pbyIsImF1dGhfdGltZSI6MTUyMjc4MzUzOSwicmV2X3NpZyI6IjU5MjA4NzY5IiwiaWF0IjoxNTIyNzgzNTQwLCJleHAiOjE1MjI4Njk5NDAsImlzcyI6Imh0dHBzOi8vaWduaXRpb24udWFhLnJ1bi5wY2ZiZXRhLmlvL29hdXRoL3Rva2VuIiwiemlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAzIiwiYXVkIjpbIjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMiIsIm9wZW5pZCJdfQ.dKoUGK93QSpOAAoxLUdT1VdfhdTW0FHSEs2jgLv77-1Nx_QLZbwkADLXitvM23GZLnFDMQ-JPjhNLo9UeOJNjtdVdzbpfwAITULvgt9k5jmeCxiIpiGiPE40VQpa7rDkpdVWLLMCbjTzowmSY3cG_FQ9FYwvPUXUzFcDqC3mDBVFhtw6eWygk1wv8lr7s57dNVmOAKY8YRJ6IQRs1pk1r7arq2bLLyXcVGT2dzPh18zsYvbbLc0JU0be6t0XBdsfts-7ZY1Yy30tXMmaaG9jihpYapfWDam1PI5cCV8fwptlx3KHfjKVunXxhK5XnNeHNKH9_UpptzaHsc8ov8Z2pYFNmnH_dbEzYS8wyky-46kLpb1mjokdHKGTAsNzV62ceaVs0quafGpLFmAtnQDxihXTeUNhb7sKQQ66vJI-SnzbiMzpdpuco8DxuGmh5fPcUkYM0rLuz_FlHCKKtjavzcjzgXJ9jEloLKnA_GrnomHk-nkCQdDk60HCih0hgo3odu3Le0vLe4jotZaQF_L1xvW0DseMg61Y44xz91Un-B2GI5l6rU7uC7xjKsLqA72aIu8a6KbdaPl9F0cKpJKZ6sGhEjLNeSI259La7dIvWalzdd7oLECT3q9YfXvpV4hfLR2PwW0YGpT7VTY9RpJwJe6uN5ZqFCxqcJ5RlV4Oc-4")
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(BeNil())
			Expect(c.Email).To(Equal("tester@pivotal.io"))
			Expect(c.UserID).To(Equal("00000000-0000-0000-0000-000000000001"))
			Expect(c.UserName).To(Equal("tester@pivotal.io"))
			Expect(c.Sub).To(Equal("00000000-0000-0000-0000-000000000001"))
		})
	})

	it("creates a valid verifier", func() {
		v := openid.NewVerifier("", "", "", false)
		Expect(v).NotTo(BeNil())
	})

	when("the verifier is nil", func() {
		it("returns an error", func() {
			p, err := f.Profile(context.Background(), nil, nil)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("unable to verify token"))
			Expect(p).To(BeNil())
		})
	})

	when("the verifier is valid", func() {
		it.Before(func() {
			f.Verifier = openid.NewVerifier("", "", "", false)
		})

		when("the token is nil", func() {
			it("returns an error", func() {
				p, err := f.Profile(context.Background(), nil, nil)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("unable to verify token"))
				Expect(p).To(BeNil())
			})
		})

		when("the token has no extra", func() {
			it("returns an error", func() {
				p, err := f.Profile(context.Background(), nil, &oauth2.Token{})
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("profile: no id_token"))
				Expect(p).To(BeNil())
			})
		})

		when("the token has includes id_token content", func() {
			var t *oauth2.Token

			it.Before(func() {
				t = &oauth2.Token{
					AccessToken:  "test-token",
					TokenType:    "bearer",
					RefreshToken: "test-refresh-token",
					Expiry:       time.Now().Add(24 * time.Hour),
				}
				extra := map[string]interface{}{
					"id_token": "test-id_token",
				}
				t = t.WithExtra(extra)
			})

			when("the token cannot be verified", func() {
				it("returns an error", func() {
					v := &openidfakes.FakeVerifier{}
					v.VerifyReturns(nil, errors.New("test error"))
					f.Verifier = v
					p, err := f.Profile(context.Background(), nil, t)
					Expect(err).NotTo(BeNil())
					Expect(err.Error()).To(ContainSubstring("unable to fetch claims: test error"))
					Expect(p).To(BeNil())
				})
			})

			when("the token can be verified", func() {
				it("returns the profile", func() {
					v := &openidfakes.FakeVerifier{}
					v.VerifyReturns(&openid.Claims{
						Email:      "test@example.net",
						UserName:   "tester",
						GivenName:  "Test",
						FamilyName: "User",
					}, nil)
					f.Verifier = v
					p, err := f.Profile(context.Background(), nil, t)
					Expect(err).To(BeNil())
					Expect(p).NotTo(BeNil())
					Expect(p.Name).To(Equal("Test User"))
					Expect(p.AccountName).To(Equal("tester"))
					Expect(p.Email).To(Equal("test@example.net"))
				})

				it("uses the email address as the account name if it is not set", func() {
					v := &openidfakes.FakeVerifier{}
					v.VerifyReturns(&openid.Claims{
						Email:      "test@example.net",
						UserName:   "",
						GivenName:  "Test",
						FamilyName: "User",
					}, nil)
					f.Verifier = v
					p, err := f.Profile(context.Background(), nil, t)
					Expect(err).To(BeNil())
					Expect(p).NotTo(BeNil())
					Expect(p.Name).To(Equal("Test User"))
					Expect(p.AccountName).To(Equal("test@example.net"))
					Expect(p.Email).To(Equal("test@example.net"))
				})
			})
		})
	})
}
