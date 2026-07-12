package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

// yamlRenderer mirrors the JSON contract by first applying encoding/json's
// tags, then encoding the resulting generic tree as YAML.
type yamlRenderer struct{}

func (yamlRenderer) Render(w io.Writer, v any) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encoding JSON shape for YAML: %w", err)
	}
	var shaped any
	if err := json.Unmarshal(raw, &shaped); err != nil {
		return fmt.Errorf("decoding JSON shape for YAML: %w", err)
	}
	data, err := yaml.Marshal(shaped)
	if err != nil {
		return fmt.Errorf("encoding YAML: %w", err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("writing YAML: %w", err)
	}
	return nil
}
