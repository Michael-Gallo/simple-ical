package parse

import (
	"fmt"
	"strings"
)

// type icalParam struct {
// 	Name  string
// 	Value string
// }

// parseIcalLine parses a single property line and returns the property name, parameters, and value.
// The propertyName is the string before the first colon or semicolon
// paramsm are colon separated values after the propertyName
// value is the string after the first colon that is not encapsulated in parentheses
func parseIcalLine(line string) (propertyName string, params []string, value string, err error) {
	propertyNameAndParams, value, ok := strings.Cut(line, ":")
	if !ok {
		return "", nil, "", fmt.Errorf("%w: %s", errInvalidPropertyLine, line)
	}
	propertyName, paramString, ok := strings.Cut(propertyNameAndParams, ";")
	if !ok {
		return propertyName, nil, value, nil
	}
	return propertyName, strings.Split(paramString, ";"), value, nil

}
