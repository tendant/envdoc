package envdoc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// EnvReader abstracts environment variable access for testability.
type EnvReader interface {
	Getenv(key string) string
	LookupEnv(key string) (string, bool)
	Environ() []string
}

// Clock abstracts time for testability.
type Clock interface {
	Now() time.Time
}

// osEnvReader implements EnvReader using the real OS environment.
type osEnvReader struct{}

func (osEnvReader) Getenv(key string) string           { return os.Getenv(key) }
func (osEnvReader) LookupEnv(key string) (string, bool) { return os.LookupEnv(key) }
func (osEnvReader) Environ() []string                   { return os.Environ() }

// realClock implements Clock using real time.
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// Inspector is the main entry point for environment variable inspection.
type Inspector struct {
	env    EnvReader
	clock  Clock
	rules  []Rule
	config Config
	output io.Writer
}

// Option configures an Inspector.
type Option func(*Inspector)

// WithEnvReader sets a custom EnvReader.
func WithEnvReader(r EnvReader) Option {
	return func(i *Inspector) { i.env = r }
}

// WithClock sets a custom Clock.
func WithClock(c Clock) Option {
	return func(i *Inspector) { i.clock = c }
}

// WithRules sets the validation rules.
func WithRules(rules []Rule) Option {
	return func(i *Inspector) { i.rules = rules }
}

// WithConfig sets the configuration directly.
func WithConfig(cfg Config) Option {
	return func(i *Inspector) { i.config = cfg }
}

// WithOutput sets the writer for log output.
func WithOutput(w io.Writer) Option {
	return func(i *Inspector) { i.output = w }
}

// New creates a new Inspector with the given options.
func New(opts ...Option) *Inspector {
	i := &Inspector{
		env:    osEnvReader{},
		clock:  realClock{},
		output: os.Stderr,
	}
	for _, opt := range opts {
		opt(i)
	}
	// If config not explicitly set, load from environment.
	if i.config == (Config{}) {
		i.config = LoadConfig(i.env)
	}
	return i
}

// Run performs inspection, logs results, checks fail-fast, and optionally starts HTTP.
// It returns the Report and any fail-fast error.
func Run(opts ...Option) (*Report, error) {
	i := New(opts...)
	return i.Run()
}

// Run performs the full inspection lifecycle.
func (i *Inspector) Run() (*Report, error) {
	report := i.Inspect()

	LogReport(i.output, report)

	if err := CheckFailFast(report, i.config.FailFast); err != nil {
		return report, err
	}

	return report, nil
}

// Inspect performs environment inspection and returns a Report.
func (i *Inspector) Inspect() *Report {
	return inspect(i.env, i.clock, i.rules, i.config)
}

// Handler returns an http.Handler for the GET /debug/env endpoint.
func (i *Inspector) Handler() http.Handler {
	return newDebugHandler(i)
}

// Config returns the inspector's configuration.
func (i *Inspector) Config() Config {
	return i.config
}

// ListenAndServe starts the HTTP debug server.
func (i *Inspector) ListenAndServe(addr string) error {
	if addr == "" {
		addr = i.config.ListenAddr
	}
	if addr == "" {
		addr = "127.0.0.1:9090"
	}
	mux := http.NewServeMux()
	mux.Handle("/debug/env", i.Handler())
	return http.ListenAndServe(addr, mux)
}

// CheckFailFast returns an error if fail-fast is enabled and there are invalid required vars.
func CheckFailFast(report *Report, failFast bool) error {
	if !failFast {
		return nil
	}
	var problems []string
	for _, r := range report.Results {
		if r.Required && (!r.Present || !r.Valid) {
			msg := fmt.Sprintf("%s: present=%t valid=%t", r.Key, r.Present, r.Valid)
			if len(r.Problems) > 0 {
				msg += " problems=[" + strings.Join(r.Problems, "; ") + "]"
			}
			problems = append(problems, msg)
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("envdoc: fail-fast: %d required variable(s) invalid:\n  %s",
			len(problems), strings.Join(problems, "\n  "))
	}
	return nil
}
