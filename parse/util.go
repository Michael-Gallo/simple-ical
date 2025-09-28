package parse

import (
	"fmt"
	"strings"
)

// parseIcalLine parses a single property line and returns the property name, parameters, and value.
// The propertyName is the string before the first colon or semicolon
// paramsm are colon separated values after the propertyName
// value is the string after the first colon that is not encapsulated in parentheses
func parseIcalLine(line string) (propertyName string, params []string, value string, err error) {
	propertyNameAndParams, value, ok := strings.Cut(line, ":")
	if !ok {
		err = fmt.Errorf("%w: %s", errInvalidPropertyLine, line)
		return
	}
	propertyName, paramString, hasParams := strings.Cut(propertyNameAndParams, ";")
	if !hasParams {
		return
	}
	params = strings.Split(paramString, ";")
	return
}
