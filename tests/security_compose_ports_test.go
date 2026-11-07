package tests

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The dev stack ships default credentials (postgres password, webhook
// secret), so every published port must stay on the loopback interface.
func TestComposePublishesPortsOnLoopbackOnly(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("../compose.yaml")
	if err != nil {
		t.Fatal(err)
	}

	published := regexp.MustCompile(`^\s+- "(\$\{[^}]+\}|\d+):\d+"\s*$`)

	for i, line := range strings.Split(string(raw), "\n") {
		if published.MatchString(line) {
			t.Errorf("compose.yaml:%d publishes a port on all interfaces: %s", i+1, strings.TrimSpace(line))
		}
	}
}
