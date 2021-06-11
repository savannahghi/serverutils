package go_utils_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	base "github.com/savannahghi/go_utils"
)

func TestMain(m *testing.M) {
	os.Setenv("MESSAGE_KEY", "this-is-a-test-key$$$")
	os.Setenv("ENVIRONMENT", "staging")
	err := os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	if err != nil {
		if base.IsDebug() {
			log.Printf("can't set root collection suffix in env: %s", err)
		}
		os.Exit(-1)
	}
	existingDebug, err := base.GetEnvVar("DEBUG")
	if err != nil {
		existingDebug = "false"
	}

	os.Setenv("DEBUG", "true")

	rc := m.Run()
	// Restore DEBUG envar to original value after running test
	os.Setenv("DEBUG", existingDebug)

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < base.CoverageThreshold {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}

	os.Exit(rc)
}
