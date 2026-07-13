package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/plugin"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
)

type schemaField struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

func TestV1StructuredOutputContract(t *testing.T) {
	actual := collectSchemas(
		model.SystemInfo{}, model.CPUInfo{}, model.MemoryInfo{}, model.DiskInfo{},
		model.ProcessList{}, service.ProcessTreeNode{}, model.NetworkInfo{},
		model.PortInfo{}, model.ContainerList{}, model.ContainerDetail{},
		model.DiagnosticReport{}, plugin.Manifest{}, plugin.Info{}, plugin.Request{},
	)
	path := "../../contracts/v1-schemas.json"
	if os.Getenv("UPDATE_CONTRACT") == "1" {
		var output bytes.Buffer
		encoder := json.NewEncoder(&output)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		require.NoError(t, encoder.Encode(actual))
		require.NoError(t, os.WriteFile(path, output.Bytes(), 0o644))
	}
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var expected map[string]map[string]schemaField
	require.NoError(t, json.Unmarshal(data, &expected))
	assert.Equal(t, expected, actual)
}

func collectSchemas(values ...any) map[string]map[string]schemaField {
	result := map[string]map[string]schemaField{}
	for _, value := range values {
		collectSchemaType(reflect.TypeOf(value), result)
	}
	return result
}

func collectSchemaType(t reflect.Type, result map[string]map[string]schemaField) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		return "array<" + collectSchemaType(t.Elem(), result) + ">"
	case reflect.Map:
		return "object<" + jsonScalarType(t.Key()) + "," + collectSchemaType(t.Elem(), result) + ">"
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return "string<date-time>"
		}
		name := t.PkgPath()[strings.LastIndex(t.PkgPath(), "/")+1:] + "." + t.Name()
		if _, seen := result[name]; seen {
			return "object<" + name + ">"
		}
		fields := map[string]schemaField{}
		result[name] = fields
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			tag := field.Tag.Get("json")
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				continue
			}
			if field.Anonymous && parts[0] == "" {
				collectSchemaType(field.Type, result)
				embedded := result[schemaTypeName(field.Type)]
				for key, value := range embedded {
					fields[key] = value
				}
				continue
			}
			fieldName := parts[0]
			if fieldName == "" {
				fieldName = field.Name
			}
			optional := contains(parts[1:], "omitempty")
			fieldType := collectSchemaType(field.Type, result)
			if !optional && nullableJSONKind(field.Type) {
				fieldType += "|null"
			}
			fields[fieldName] = schemaField{Type: fieldType, Required: !optional}
		}
		return "object<" + name + ">"
	default:
		return jsonScalarType(t)
	}
}

func nullableJSONKind(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface:
		return true
	default:
		return false
	}
}

func schemaTypeName(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.PkgPath()[strings.LastIndex(t.PkgPath(), "/")+1:] + "." + t.Name()
}

func jsonScalarType(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Interface:
		return "any"
	default:
		return t.String()
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
