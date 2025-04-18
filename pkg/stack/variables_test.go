package stack

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"bytes", []byte("hello"), "hello"},
		{"yaml.Node string", yaml.Node{Kind: yaml.ScalarNode, Value: "world"}, "world"},
		{"int fallback", 123, "123"},
		{"float fallback", 45.67, "45.67"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parseString(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, out)
		})
	}
}

func TestParseInt(t *testing.T) {
	t.Run("int from int", func(t *testing.T) {
		val, err := parseInt(42)
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("int from float", func(t *testing.T) {
		val, err := parseInt(42.0)
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("int from string", func(t *testing.T) {
		val, err := parseInt("42")
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("invalid int", func(t *testing.T) {
		_, err := parseInt("not a number")
		assert.Error(t, err)
	})
}

func TestParseFloat(t *testing.T) {
	t.Run("float from float", func(t *testing.T) {
		val, err := parseFloat(3.14)
		assert.NoError(t, err)
		assert.Equal(t, 3.14, val)
	})

	t.Run("float from string", func(t *testing.T) {
		val, err := parseFloat("3.14")
		assert.NoError(t, err)
		assert.Equal(t, 3.14, val)
	})

	t.Run("invalid float", func(t *testing.T) {
		_, err := parseFloat("abc")
		assert.Error(t, err)
	})
}

func TestParseIP(t *testing.T) {
	t.Run("valid IP", func(t *testing.T) {
		val, err := parseIP("192.168.1.1")
		assert.NoError(t, err)
		assert.Equal(t, net.ParseIP("192.168.1.1"), val)
	})

	t.Run("invalid IP format", func(t *testing.T) {
		_, err := parseIP("not.an.ip")
		assert.Error(t, err)
	})

	t.Run("non-string input", func(t *testing.T) {
		_, err := parseIP(12345)
		assert.Error(t, err)
	})
}
func TestIsAllowed(t *testing.T) {
	t.Run("value allowed", func(t *testing.T) {
		ok, err := isAllowed("dev", []interface{}{"dev", "staging", "prod"})
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("value not allowed", func(t *testing.T) {
		ok, err := isAllowed("qa", []interface{}{"dev", "staging", "prod"})
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("value allowed by regex", func(t *testing.T) {
		ok, err := isAllowed("test-123", []interface{}{`test-[0-9]+`})
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("value not matching regex", func(t *testing.T) {
		ok, err := isAllowed("prod-abc", []interface{}{`test-[0-9]+`})
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("invalid regex falls back to string compare", func(t *testing.T) {
		ok, err := isAllowed("dev", []interface{}{`[`, `dev`, `prod`})
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("empty allowed list", func(t *testing.T) {
		ok, err := isAllowed("anything", []interface{}{})
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("non-string value coerced", func(t *testing.T) {
		ok, err := isAllowed(123, []interface{}{`123`})
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("unicode match", func(t *testing.T) {
		ok, err := isAllowed("環境", []interface{}{"環境"})
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("mixed valid and invalid regex with fallback", func(t *testing.T) {
		ok, err := isAllowed("dev", []interface{}{`dev`, `[`, `prod`})
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
