package envdoc

import (
	"fmt"
	"io"
)

// LogReport writes a one-line-per-variable summary to w.
func LogReport(w io.Writer, report *Report) {
	for _, r := range report.Results {
		line := fmt.Sprintf("envdoc: key=%s present=%t", r.Key, r.Present)
		if r.Present {
			line += fmt.Sprintf(" len=%d", r.Length)
		}
		if r.Fingerprint != "" {
			line += fmt.Sprintf(" fp=%s", r.Fingerprint)
		}
		if r.Required {
			line += fmt.Sprintf(" required=%t", r.Required)
		}
		line += fmt.Sprintf(" valid=%t", r.Valid)
		if r.Trimmed {
			line += " trimmed=true"
		}
		if len(r.Problems) > 0 {
			for _, p := range r.Problems {
				line += fmt.Sprintf(" problem=%q", p)
			}
		}
		fmt.Fprintln(w, line)
	}
}
