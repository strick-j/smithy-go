package http

import (
	"context"
	"testing"

	smithy "github.com/strick-j/smithy-go"
	"github.com/strick-j/smithy-go/auth"
)

func TestIdentity(t *testing.T) {
	var expected auth.Identity = &auth.AnonymousIdentity{}

	resolver := auth.AnonymousIdentityResolver{}
	actual, _ := resolver.GetIdentity(context.TODO(), smithy.Properties{})
	if expected != actual {
		t.Errorf("Anonymous identity resolver does not produce correct identity")
	}
}
