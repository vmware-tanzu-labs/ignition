package api

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/http/session"
	"github.com/pivotalservices/ignition/user"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUserInfoFromContext(t *testing.T) {
	spec.Run(t, "UserInfoFromContext", testUserInfoFromContext, spec.Report(report.Terminal{}))
}

func testUserInfoFromContext(t *testing.T, when spec.G, it spec.S) {
	var ctx context.Context
	it.Before(func() {
		RegisterTestingT(t)
		ctx = context.Background()
	})

	when("there is no user", func() {
		it("errors", func() {
			_, _, err := userInfoFromContext(ctx)
			Expect(err).To(HaveOccurred())
		})
	})

	when("the profile is nil", func() {
		it.Before(func() {
			ctx = user.WithProfile(ctx, nil)
		})
		it("errors", func() {
			_, _, err := userInfoFromContext(ctx)
			Expect(err).To(HaveOccurred())
		})
	})

	when("there is a valid profile", func() {
		it.Before(func() {
			ctx = user.WithProfile(ctx, &user.Profile{
				AccountName: "test-user",
			})
		})

		when("there is no user id", func() {
			it("errors", func() {
				_, _, err := userInfoFromContext(ctx)
				Expect(err).To(HaveOccurred())
			})
		})

		when("there is a user ID", func() {
			it.Before(func() {
				ctx = session.ContextWithUserID(ctx, "test-user-id")
			})

			it("returns the user id and account name", func() {
				userID, accountName, err := userInfoFromContext(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(userID).To(Equal("test-user-id"))
				Expect(accountName).To(Equal("test-user"))
			})
		})
	})
}
