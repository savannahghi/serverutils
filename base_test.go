package base_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestMain(m *testing.M) {
	os.Setenv("MESSAGE_KEY", "this-is-a-test-key$$$")
	err := os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	if err != nil {
		log.Printf("can't set root collection suffix in env: %s", err)
		os.Exit(-1)
	}
	rc := m.Run()

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
