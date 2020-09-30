package fxchain

import (
	"strings"
)

// Pair represents effect name & parameters
type Pair struct {
	effect string
	params map[string]string
}

// Chain is a slice of Pairs
type Chain []Pair

// GetStmts retrieves statements from a string
func GetStmts(src string) []string {
	stmts := []string{}
	last := ""

	for _, c := range src {
		if c == '\n' || c == ';' {
			stmts = append(stmts, strings.TrimSpace(last))
			last = ""
		} else {
			last += string(c)
		}
	}

	if last != "" {
		stmts = append(stmts, strings.TrimSpace(last))
	}

	return stmts
}

// GetStmtValues retrieves effect name & parameters from a statement
func GetStmtValues(stmt string) (string, map[string]string) {
	effect := ""
	params := make(map[string]string)

	var (
		last          string
		prev          byte
		hasEffect     bool
		paramKey      string
		hasParamKey   bool
		paramValue    string
		hasParamValue bool
	)

	// for every character in the statement
	for i, c := range stmt {
		// we don't have effect name
		if !hasEffect {
			// check if we are opening the parameters
			if c == '{' {
				// if so, the effect is everything we saw till now
				effect = last
				// we have effect name now
				hasEffect = true
				// reset memory
				last = ""
			} else if c != '\\' {
				// otherwise, save this character for later
				last += string(c)
			}
			// we have effect name, so we're in parameters
		} else {
			// check if params end here
			if c == '}' && prev != '\\' {
				// check if we have leftover param key and/or value
				if hasParamKey {
					// process key & value for possibly unwanted characters
					paramKey = "fx-" + strings.ToLower(strings.TrimSpace(paramKey))
					paramValue = strings.TrimSpace(last)
					// save in parameters
					params[paramKey] = paramValue
					// get outta here
					break
				}
			}
			// we don't have param key
			if !hasParamKey {
				// check if key ends here
				if c == ':' {
					// if so, the key is everything we saw till now
					paramKey = last
					// we have key now
					hasParamKey = true
					// reset memory
					last = ""
				} else if c != '\\' {
					// otherwise, save this character for later
					last += string(c)
				}
				// we have param key, but not value
			} else if !hasParamValue {
				// get previous character for escapes
				prev = stmt[i-1]
				// check if we're switching to next key
				if c == ',' && prev != '\\' {
					// if so, the value is everything we saw till now
					paramValue = last
					// we have value now
					hasParamValue = true
					// reset memory
					last = ""
				} else if c != '\\' {
					// otherwise, save this character for later
					last += string(c)
				}
				// we have param key & value
			} else {
				// process key & value for possibly unwanted characters
				paramKey = "fx-" + strings.ToLower(strings.TrimSpace(paramKey))
				paramValue = strings.TrimSpace(paramValue)
				// save in parameters
				params[paramKey] = paramValue
				// reset memory
				hasParamKey = false
				hasParamValue = false
				// set last to current char, as it will be ignored otherwise
				last = string(c)
			}
		}
	}

	return effect, params
}

// ParseChain parses source and get effect chain
func ParseChain(src string) Chain {
	stmts := GetStmts(src)
	pairs := []Pair{}
	for _, stmt := range stmts {
		effect, params := GetStmtValues(stmt)
		pairs = append(pairs, Pair{effect, params})
	}
	return Chain(pairs)
}
