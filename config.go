package main

type Config struct {
	// GitHub API
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`

	// Session secret
	SessionSecret string `json:"sessionSecret"`

	// HTTP server
	Port int `json:"port"`

	// MySQL server
	DBHost string `json:"dbHost"`
	DBName string `json:"dbName"`
	DBUser string `json:"dbUser"`
	DBPass string `json:"dbPass"`
}
