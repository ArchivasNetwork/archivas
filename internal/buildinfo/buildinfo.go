package buildinfo

// Build information - set via ldflags at compile time
var (
	Version     = "dev"
	Commit      = "unknown"
	BuiltAt     = "unknown"
	PoSpaceRule = "quality<=difficulty" // MUST be this value for v1.1.1+
)

// GetInfo returns build information as a map
func GetInfo() map[string]string {
	return map[string]string{
		"version":     Version,
		"commit":      Commit,
		"builtAt":     BuiltAt,
		"poSpaceRule": PoSpaceRule,
	}
}

