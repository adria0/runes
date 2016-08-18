package config

// AuthGoogle is google OAuth2 authentication
const AuthGoogle = "google"

// AuthNone is no authentication
const AuthNone = "none"

// Config is the server configurtion
type Config struct {
	Port     int    // Port to listen to
	DataDir  string // Where blog content are served
	TmpDir   string // A temporally directory
	CacheDir string // Cache directory for generated images
	Auth     struct {
		Type           string   // Type of authenticatiom, AuthGoogle or AuthNone
		GoogleClientID string   // Client id iff AuthGoogle
		AllowedEmails  []string // Allowed emails to log in iff Authgoogle
	}
}
