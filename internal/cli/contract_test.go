package cli

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/plugin"
)

type cliContract struct {
	Formats           []string            `json:"formats"`
	GlobalFlags       []string            `json:"global_flags"`
	Commands          map[string][]string `json:"commands"`
	ConfigKeys        []string            `json:"config_keys"`
	ExitCodes         map[string]string   `json:"exit_codes"`
	PluginAPIVersions []string            `json:"plugin_api_versions"`
}

func TestV1CLIContract(t *testing.T) {
	data, err := os.ReadFile("../../contracts/v1-cli.json")
	require.NoError(t, err)
	var want cliContract
	require.NoError(t, json.Unmarshal(data, &want))

	root := newRootCmd()
	assert.Equal(t, []string{formatJSON, formatTable, formatYAML}, want.Formats)
	assert.Equal(t, want.GlobalFlags, contractFlags(root))
	assert.Equal(t, want.Commands, commandContract(root))
	assert.Equal(t, configContractKeys(), want.ConfigKeys)
	assert.Equal(t, map[string]string{
		strconv.Itoa(exitOK): "success", strconv.Itoa(exitError): "general error", strconv.Itoa(exitUsage): "usage error",
		strconv.Itoa(exitPermission): "permission", strconv.Itoa(exitUnsupported): "unsupported", strconv.Itoa(exitPartial): "partial failure",
	}, want.ExitCodes)
	assert.Equal(t, []string{plugin.APIVersion}, want.PluginAPIVersions)
}

func configContractKeys() []string {
	keys := make([]string, 0, len(knownGlobalKeys))
	for key := range knownGlobalKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func commandContract(root *cobra.Command) map[string][]string {
	root.InitDefaultCompletionCmd()
	result := map[string][]string{}
	var walk func(*cobra.Command)
	walk = func(parent *cobra.Command) {
		for _, cmd := range parent.Commands() {
			if cmd.Hidden || cmd.Name() == "help" {
				continue
			}
			path := strings.TrimPrefix(cmd.CommandPath(), root.Name()+" ")
			result[path] = contractFlags(cmd)
			walk(cmd)
		}
	}
	walk(root)
	return result
}

func contractFlags(cmd *cobra.Command) []string {
	set := map[string]struct{}{}
	cmd.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "help" {
			set[flag.Name] = struct{}{}
		}
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "help" {
			set[flag.Name] = struct{}{}
		}
	})
	flags := make([]string, 0, len(set))
	for flag := range set {
		flags = append(flags, flag)
	}
	sort.Strings(flags)
	return flags
}
