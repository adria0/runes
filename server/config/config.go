package config

const (
    AuthGoogle = "google"
    AuthNone = "none"
)

type Config struct {
	Port   int
	Prefix string
	Auth   struct {
		Type           string
		GoogleClientID string
		AllowedEmails  []string
	}
}

