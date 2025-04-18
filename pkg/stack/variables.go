package stack

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Validates a variable value based on the variable definition
func ValidateVar(input any, definition StackVariable) (any, error) {
	var val any
	var err error

	// Check the variable is in the allowed list
	ok, err := isAllowed(input, definition.Allowed)
	if err != nil {
		return nil, fmt.Errorf("failed validate check allowed list: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("value %q is not an allowed value", input)
	}

	// Parse as the type to validate
	switch definition.Type {
	case VarTypeString:
		val, err = parseString(input)
	case VarTypeInt:
		val, err = parseInt(input)
	case VarTypeFloat:
		val, err = parseFloat(input)
	case VarTypeIP:
		val, err = parseIP(input)
	default:
		err = fmt.Errorf("unknown variable type %q", definition.Type)
	}

	return val, err
}

func parseString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case yaml.Node:
		if v.Kind == yaml.ScalarNode {
			return v.Value, nil
		}
	}
	return fmt.Sprint(val), nil
}

func parseInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case json.Number:
		i, err := v.Int64()
		return int(i), err
	case string:
		return strconv.Atoi(v)
	case yaml.Node:
		if v.Kind == yaml.ScalarNode {
			return strconv.Atoi(v.Value)
		}
	}
	return 0, fmt.Errorf("not an int-compatible type (got %T)", val)
}

func parseFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case json.Number:
		return v.Float64()
	case string:
		return strconv.ParseFloat(v, 64)
	case yaml.Node:
		if v.Kind == yaml.ScalarNode {
			return strconv.ParseFloat(v.Value, 64)
		}
	}
	return 0, fmt.Errorf("not a float-compatible type (got %T)", val)
}

func parseIP(val interface{}) (net.IP, error) {
	switch v := val.(type) {
	case string:
		ip := net.ParseIP(v)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP")
		}
		return ip, nil
	case net.IP:
		return v, nil
	case yaml.Node:
		if v.Kind == yaml.ScalarNode {
			ip := net.ParseIP(v.Value)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP")
			}
			return ip, nil
		}
	}
	return nil, fmt.Errorf("not a string or IP-compatible type (got %T)", val)
}

func isAllowed(value any, allowed []any) (bool, error) {
	valStr := fmt.Sprint(value)
	for _, a := range allowed {
		pattern := fmt.Sprint(a)
		re, err := regexp.Compile(pattern)
		if err != nil {
			// If not valid regex, perform string comparison
			if pattern == valStr {
				return true, nil
			}
			continue
		}
		if re.MatchString(valStr) {
			return true, nil
		}
	}
	return false, nil
}
