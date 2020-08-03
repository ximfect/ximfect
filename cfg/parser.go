package cfg

import "strings"

// Parse reads the source string and outputs the parsed Config and an error if any
func Parse(src string) Config {
	out := make(Config)
	var (
		key   string
		value string
	)
	for _, line := range strings.Split(strings.TrimSpace(src), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(strings.Split(line, "=")[0]))
		value = strings.Join(strings.Split(line, "=")[1:], "=")
		out[key] = value
	}
	return out
}
