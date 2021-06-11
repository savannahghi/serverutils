package server_utils

import (
	"strconv"

	"github.com/savannahghi/go_utils"
)

// BoolEnv gets and parses a boolean environment variable
func BoolEnv(envVarName string) bool {
	envVar, err := go_utils.GetEnvVar(envVarName)
	if err != nil {
		return false
	}
	val, err := strconv.ParseBool(envVar)
	if err != nil {
		return false
	}
	return val
}

// IsDebug returns true if debug has been turned on in the environment
func IsDebug() bool {
	return BoolEnv(DebugEnvVarName)
}

// IsRunningTests returns true if debug has been turned on in the environment
func IsRunningTests() bool {
	return BoolEnv(IsRunningTestsEnvVarName)
}
