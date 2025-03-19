package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultJSONConfig(t *testing.T) {
	// Act
	cfg := DefaultJSONConfig()

	// Assert
	assert.NotNil(t, cfg)
	assert.True(t, cfg.UseProtoNames)
	assert.True(t, cfg.EmitUnpopulated)
	assert.True(t, cfg.UseEnumNumbers)
	assert.True(t, cfg.AllowPartial)
	assert.True(t, cfg.Multiline)
	assert.Equal(t, "  ", cfg.Indent)
}
