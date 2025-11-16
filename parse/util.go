package parse

import (
	"fmt"
	"strings"
)

// parseIcalLineWithReusableMap parses a single property line using a reusable parameter map.
// This avoids allocating a new map for every property with parameters.
func parseIcalLineWithReusableMap(line string, reusableParams map[string]string) (propertyName string, params map[string]string, value string, err error) {
	// Find the first colon that is not inside quotes
	colonIndex := findUnquotedColonIndex(line)
	if colonIndex == -1 {
		err = fmt.Errorf("%w: %s", errInvalidPropertyLine, line)
		return "", nil, "", err
	}

	// Split the line at the colon
	beforeColon := line[:colonIndex]

	// The property name is the first part before any semicolon
	propertyName = beforeColon
	if semicolonIndex := strings.Index(beforeColon, ";"); semicolonIndex != -1 {
		propertyName = beforeColon[:semicolonIndex]
		// Extract parameters from the part between property name and colon
		paramString := beforeColon[semicolonIndex+1:]
		if paramString != "" {
			// Use the reusable map (caller has already cleared it)
			params = reusableParams
			splitParametersWithReusableMap(paramString, params)
		}
	}

	return propertyName, params, line[colonIndex+1:], nil
}

// splitParametersWithReusableMap splits parameters using a reusable map and string builder.
// This avoids allocating new objects for every parameter parsing operation.
func splitParametersWithReusableMap(paramString string, params map[string]string) {
	var current strings.Builder
	// Pre-allocate capacity based on typical parameter length
	current.Grow(len(paramString) / 2)

	var currentKey string
	inQuotes := false

	for _, character := range paramString {
		switch character {
		case '"':
			inQuotes = !inQuotes
		case '=':
			if inQuotes {
				current.WriteRune(character)
				continue
			}
			currentKey = current.String()
			current.Reset()
		case ';':
			if inQuotes {
				current.WriteRune(character)
				continue
			}
			// Found a parameter separator, write the parameter.
			if current.Len() > 0 {
				params[currentKey] = current.String()
				current.Reset()
			}
		default:
			current.WriteRune(character)
		}
	}
	// Write the last parameter (it never hit a semicolon).
	if current.Len() > 0 {
		params[currentKey] = current.String()
	}
}

// splitParameters splits a parameter string by semicolons, respecting quoted strings.
func splitParameters(paramString string) map[string]string {
	var params = make(map[string]string, 1)
	var current strings.Builder
	var currentKey string
	inQuotes := false

	for _, character := range paramString {
		switch character {
		case '"':
			inQuotes = !inQuotes
		case '=':
			if inQuotes {
				current.WriteRune(character)
				continue
			}
			currentKey = current.String()
			current.Reset()
		case ';':
			if inQuotes {
				current.WriteRune(character)
				continue
			}
			// Found a parameter separator, write the parameter.
			if current.Len() > 0 {
				params[currentKey] = current.String()
				current.Reset()
			}
		default:
			current.WriteRune(character)
		}
	}
	// Write the last parameter (it never hit a semicolon).
	if current.Len() > 0 {
		params[currentKey] = current.String()
	}
	return params
}

// findUnquotedColonIndex finds the first colon that is not encapsulated in quotations.
func findUnquotedColonIndex(line string) int {
	inQuotes := false
	for i, c := range line {
		if c == '"' {
			inQuotes = !inQuotes
		} else if c == ':' && !inQuotes {
			return i
		}
	}
	return -1
}
