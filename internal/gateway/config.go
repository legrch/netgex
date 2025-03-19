package gateway

// JSONConfig holds JSON marshaling configuration for the gateway
type JSONConfig struct {
	// UseProtoNames determines if proto or json names should be used
	UseProtoNames bool `envconfig:"USE_PROTO_NAMES" default:"true"`
	// EmitUnpopulated determines if unpopulated fields should be included
	EmitUnpopulated bool `envconfig:"EMIT_UNPOPULATED" default:"true"`
	// UseEnumNumbers renders enum values as numbers instead of strings
	UseEnumNumbers bool `envconfig:"USE_ENUM_NUMBERS" default:"true"`
	// AllowPartial allows incomplete proto messages
	AllowPartial bool `envconfig:"ALLOW_PARTIAL" default:"true"`
	// Multiline formats the output in indented form
	Multiline bool `envconfig:"MULTILINE" default:"true"`
	// Indent specifies the set of indentation characters to use in multiline mode
	Indent string `envconfig:"INDENT" default:"  "`
}

// DefaultJSONConfig returns a JSONConfig with default values
func DefaultJSONConfig() *JSONConfig {
	return &JSONConfig{
		UseProtoNames:   true,
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
		AllowPartial:    true,
		Multiline:       true,
		Indent:          "  ",
	}
}
