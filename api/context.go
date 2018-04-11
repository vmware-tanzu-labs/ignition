package api

import (
	"context"
	"errors"

	"github.com/pivotalservices/ignition/http/session"
	"github.com/pivotalservices/ignition/user"
)

func userInfoFromContext(ctx context.Context) (userID string, accountName string, err error) {
	var profile *user.Profile
	profile, err = user.ProfileFromContext(ctx)
	if err != nil {
		return "", "", err
	}
	if profile == nil {
		return "", "", errors.New("no profile was found")
	}
	userID, err = session.UserIDFromContext(ctx)
	if err != nil {
		return "", "", err
	}
	return userID, profile.AccountName, nil
}
