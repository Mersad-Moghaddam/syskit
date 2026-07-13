package model

type DiagnosticFinding struct {
	ID             string   `json:"id"`
	Severity       string   `json:"severity"`
	Category       string   `json:"category"`
	Summary        string   `json:"summary"`
	Evidence       string   `json:"evidence"`
	Sources        []string `json:"sources"`
	Recommendation string   `json:"recommendation"`
}
type DiagnosticReport struct {
	Findings []DiagnosticFinding `json:"findings"`
}
