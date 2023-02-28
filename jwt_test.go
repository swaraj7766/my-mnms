package mnms

import (
	"os"
	"testing"
)

func TestEnvironmentVairable(t *testing.T) {
	token := os.Getenv("MNMS_TOKEN")
	t.Log("token: ", token)
}
