// Code generated by counterfeiter. DO NOT EDIT.
package openidfakes

import (
	"context"
	"sync"

	"github.com/coreos/go-oidc"
	"github.com/pivotalservices/ignition/user/openid"
)

type FakeOIDCVerifier struct {
	VerifyStub        func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		ctx        context.Context
		rawIDToken string
	}
	verifyReturns struct {
		result1 *oidc.IDToken
		result2 error
	}
	verifyReturnsOnCall map[int]struct {
		result1 *oidc.IDToken
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOIDCVerifier) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		ctx        context.Context
		rawIDToken string
	}{ctx, rawIDToken})
	fake.recordInvocation("Verify", []interface{}{ctx, rawIDToken})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(ctx, rawIDToken)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.verifyReturns.result1, fake.verifyReturns.result2
}

func (fake *FakeOIDCVerifier) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *FakeOIDCVerifier) VerifyArgsForCall(i int) (context.Context, string) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].ctx, fake.verifyArgsForCall[i].rawIDToken
}

func (fake *FakeOIDCVerifier) VerifyReturns(result1 *oidc.IDToken, result2 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 *oidc.IDToken
		result2 error
	}{result1, result2}
}

func (fake *FakeOIDCVerifier) VerifyReturnsOnCall(i int, result1 *oidc.IDToken, result2 error) {
	fake.VerifyStub = nil
	if fake.verifyReturnsOnCall == nil {
		fake.verifyReturnsOnCall = make(map[int]struct {
			result1 *oidc.IDToken
			result2 error
		})
	}
	fake.verifyReturnsOnCall[i] = struct {
		result1 *oidc.IDToken
		result2 error
	}{result1, result2}
}

func (fake *FakeOIDCVerifier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeOIDCVerifier) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ openid.OIDCVerifier = new(FakeOIDCVerifier)
