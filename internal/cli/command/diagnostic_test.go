package command

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type diagnosticServiceStub struct {
	report *model.DiagnosticReport
}

func (s diagnosticServiceStub) Collect(string, string) (*model.DiagnosticReport, error) {
	return s.report, nil
}

func TestDiagnosticCommandStructuredFindingContract(t *testing.T) {
	report := &model.DiagnosticReport{Findings: []model.DiagnosticFinding{{
		ID: "memory-pressure", Severity: "warning", Category: "memory",
		Summary: "Memory pressure is elevated", Evidence: "full PSI avg10 is 12.00%",
		Sources: []string{"/proc/pressure/memory"}, Recommendation: "inspect memory-heavy processes",
	}}}
	cmd := NewDiagnosticCmd(diagnosticServiceStub{report}, func() string { return "json" }, func() bool { return false }, func() bool { return false })
	var output bytes.Buffer
	cmd.SetOut(&output)
	require.NoError(t, cmd.Execute())

	var decoded model.DiagnosticReport
	require.NoError(t, json.Unmarshal(output.Bytes(), &decoded))
	require.Len(t, decoded.Findings, 1)
	assert.Equal(t, report.Findings[0], decoded.Findings[0])
	assert.NotContains(t, output.String(), "\x1b")
}
