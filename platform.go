package loungeup

// Platform of the application (development, studio, production).
type Platform int

const (
	PlatformUnknown Platform = iota
	PlatformDevelopment
	PlatformStudio
	PlatformProduction
)

func MapStringToPlatform(s string) Platform {
	switch s {
	case "development":
		return PlatformDevelopment
	case "studio":
		return PlatformStudio
	case "production":
		return PlatformProduction
	default:
		return PlatformUnknown
	}
}

func (p Platform) String() string {
	switch p {
	case PlatformDevelopment:
		return "development"
	case PlatformStudio:
		return "studio"
	case PlatformProduction:
		return "production"
	default:
		return "unknown"
	}
}
