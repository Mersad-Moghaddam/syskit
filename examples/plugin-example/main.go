// plugin-example is a minimal SysKit protocol-v1 collector plugin.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type request struct {
	APIVersion string `json:"api_version"`
	Action     string `json:"action"`
}

func main() {
	var req request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if req.APIVersion != "v1" || req.Action != "collect" {
		fmt.Fprintln(os.Stderr, "unsupported request")
		os.Exit(1)
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"status": "ok", "example_value": 1})
}
