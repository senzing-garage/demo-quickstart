//go:build linux

package cmd

import "github.com/senzing-garage/go-cmdhelping/option"

var ContextVariablesForOsArch = []option.ContextVariable{}

const SenzingToolsDatabaseURL = "sqlite3://na:na@/tmp/sqlite/G2C.db"
