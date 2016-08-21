package api

import (
	"testing"

	"github.com/vjeantet/goldap/message"
)

func TestResolveUsername(t *testing.T) {
	baseDN := "dc=dario,dc=im"
	tests := []struct {
		dn       string
		expected string
	}{
		{"cn=admin,dc=dario,dc=im", "admin"},
		{"uid=root,dc=dario,dc=im", "root"},
		{"uid=nobody,cn=wrong,dc=dario,dc=im", ""},
		{"uid=root,dc=example,dc=com", ""},
	}
	for _, tst := range tests {
		result, err := resolveUsername(message.LDAPDN(tst.dn), baseDN)
		if err != nil && tst.expected != "" {
			t.Fatal(err)
		}
		if result != tst.expected {
			t.Fatalf("expected '%s', got '%s'", tst.expected, result)
		}
	}
}
