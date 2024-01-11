//go:build darwin

package cmd

import "github.com/senzing-garage/go-cmdhelping/option"

var ContextVariablesForOsArch = []option.ContextVariable{
	option.SenzingDirectory,
	option.ConfigPath,
	option.ResourcePath,
	option.SupportPath,
}
