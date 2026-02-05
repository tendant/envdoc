package envdoc

import "time"

// MapEnvReader implements EnvReader using a map for testing.
type MapEnvReader map[string]string

func (m MapEnvReader) Getenv(key string) string {
	return m[key]
}

func (m MapEnvReader) LookupEnv(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func (m MapEnvReader) Environ() []string {
	var pairs []string
	for k, v := range m {
		pairs = append(pairs, k+"="+v)
	}
	return pairs
}

// fixedClock implements Clock with a fixed time.
type fixedClock struct {
	t time.Time
}

func (c fixedClock) Now() time.Time { return c.t }

func boolPtr(b bool) *bool { return &b }
func intPtr(i int) *int    { return &i }
