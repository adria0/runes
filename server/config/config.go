package config

const (
    AuthGoogle = "google"
    AuthNone = "none"
)

type Config struct {
	Port   int
	Prefix string
	DataDir string
    TmpDir string
    CacheDir string
    Auth   struct {
		Type           string
		GoogleClientID string
		AllowedEmails  []string
	}
}

