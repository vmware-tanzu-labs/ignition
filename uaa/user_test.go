package uaa_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/internal"
	"github.com/pivotalservices/ignition/uaa"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"golang.org/x/oauth2"
)

func TestUserIDForAccountName(t *testing.T) {
	spec.Run(t, "UserIDForAccountName", testUserIDForAccountName, spec.Report(report.Terminal{}))
}

func testUserIDForAccountName(t *testing.T, when spec.G, it spec.S) {
	var a *uaa.Client

	it.Before(func() {
		RegisterTestingT(t)
		a = &uaa.Client{}
	})

	it("cannot find a user id for an empty account name", func() {
		userID, err := a.UserIDForAccountName("")
		Expect(err).To(HaveOccurred())
		Expect(userID).To(BeZero())
		Expect(err.Error()).To(Equal("cannot search for a user with an empty account name"))
	})

	when("authentication is required", func() {
		var s *httptest.Server

		it.Before(func() {
			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			a.URL = s.URL
		})

		it.After(func() {
			s.Close()
		})

		it("returns an error", func() {
			userID, err := a.UserIDForAccountName("test-user")
			Expect(err).To(HaveOccurred())
			Expect(userID).To(BeZero())
			Expect(err.Error()).To(ContainSubstring("uaa: cannot authenticate"))
		})
	})

	when("there is a valid client and token", func() {
		var (
			s      *httptest.Server
			called bool
		)

		it.Before(func() {
			called = false
			a.Client = http.DefaultClient
			a.Token = &oauth2.Token{
				Expiry:      time.Now().Add(24 * time.Hour),
				AccessToken: "test-token",
			}
		})

		it.After(func() {
			s.Close()
		})

		when("a valid user is returned", func() {
			it.Before(func() {
				s = internal.ServeFromTestdata(t, "users.json", func() {
					called = true
				})
				a.URL = s.URL
			})

			it("returns the user id", func() {
				userID, err := a.UserIDForAccountName("tester@pivotal.io")
				Expect(err).NotTo(HaveOccurred())
				Expect(userID).To(Equal("abcdef11-0000-dddd-aaaa-1234567890ab"))
				Expect(called).To(BeTrue())
			})
		})

		when("the users call fails", func() {
			it.Before(func() {
				s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					called = true
					w.WriteHeader(http.StatusInternalServerError)
				}))
				a.URL = s.URL
			})

			it("returns the error", func() {
				userID, err := a.UserIDForAccountName("tester@pivotal.io")
				Expect(err).To(HaveOccurred())
				Expect(userID).To(BeZero())
				Expect(called).To(BeTrue())
			})
		})

		when("an empty user is returned", func() {
			it.Before(func() {
				s = internal.ServeFromTestdata(t, "empty-user.json", func() {
					called = true
				})
				a.URL = s.URL
			})

			it("returns an error", func() {
				userID, err := a.UserIDForAccountName("tester@pivotal.io")
				Expect(err).To(HaveOccurred())
				Expect(userID).To(BeZero())
				Expect(err.Error()).To(Equal("cannot find user with account name: [tester@pivotal.io]"))
				Expect(called).To(BeTrue())
			})
		})
	})
}

func TestCreateUser(t *testing.T) {
	spec.Run(t, "CreateUser", testCreateUser, spec.Report(report.Terminal{}))
}

func testCreateUser(t *testing.T, when spec.G, it spec.S) {
	var (
		a *uaa.Client
		s *httptest.Server
	)

	it.Before(func() {
		RegisterTestingT(t)
		a = &uaa.Client{}
	})

	it.After(func() {
		if s != nil {
			s.Close()
		}
	})

	when("there are not valid credentials, and no token or client", func() {
		it("fails to create the user when there are invalid credentials", func() {
			a.Client = nil
			_, err := a.CreateUser("user", "uaa", "external-user", "user@example.com")
			Expect(err).To(HaveOccurred())
		})
	})

	when("there is a valid token and client", func() {
		it.Before(func() {
			a.Client = http.DefaultClient
			a.Token = &oauth2.Token{
				Expiry:      time.Now().Add(24 * time.Hour),
				AccessToken: "test-token",
			}
		})

		when("the user is created successfully", func() {
			it.Before(func() {
				s = internal.ServeFromTestdata(t, "created-user.json", func() {})
				a.URL = s.URL
			})

			it("returns the userID", func() {
				userID, err := a.CreateUser("user", "uaa", "external-user", "user@example.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(userID).To(Equal("abcdef11-0000-dddd-aaaa-1234567890ab"))
			})
		})

		when("the user cannot be created", func() {
			it.Before(func() {
				s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				a.URL = s.URL
			})

			it("is an error", func() {
				userID, err := a.CreateUser("user", "uaa", "external-user", "user@example.com")
				Expect(err).To(HaveOccurred())
				Expect(userID).To(BeZero())
			})
		})
	})

}
