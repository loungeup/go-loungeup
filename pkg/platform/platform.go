// Package platform provides a unified way to represent the platform of the application (development, studio, production).
package platform

type Platform int

const (
	Unknown Platform = iota
	Development
	Studio
	Production
)

func FromString(s string) Platform {
	switch s {
	case "development":
		return Development
	case "studio":
		return Studio
	case "production":
		return Production
	default:
		return Unknown
	}
}

func (p Platform) String() string {
	switch p {
	case Development:
		return "development"
	case Studio:
		return "studio"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}
