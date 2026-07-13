package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type DiagnosticService interface {
	Collect(string, string) (*model.DiagnosticReport, error)
}

func NewDiagnosticCmd(s DiagnosticService, format func() string, noHeader func() bool) *cobra.Command {
	var category, severity string
	cmd := &cobra.Command{Use: "diagnostics", Short: "Show read-only health findings", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		report, err := s.Collect(category, severity)
		if err != nil {
			return fmt.Errorf("collecting diagnostics: %w", err)
		}
		r, err := render.New(format(), render.WithNoHeader(noHeader()))
		if err != nil {
			return err
		}
		if format() != "table" {
			return r.Render(c.OutOrStdout(), report)
		}
		t := render.Table{Headers: []string{"SEVERITY", "CATEGORY", "FINDING", "EVIDENCE"}}
		for _, f := range report.Findings {
			t.Rows = append(t.Rows, []string{f.Severity, f.Category, f.Summary, f.Evidence})
		}
		return r.Render(c.OutOrStdout(), t)
	}}
	cmd.Flags().StringVar(&category, "category", "", "filter by category")
	cmd.Flags().StringVar(&severity, "severity", "", "filter by severity")
	return cmd
}
