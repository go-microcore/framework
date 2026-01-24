package client // import "go.microcore.dev/framework/db/postgres/client"

import (
	"net/url"
	"regexp"
)

// List of sensitive parameter keys whose values should be masked in logs.
var sensitiveKeys = []string{
	// User authentication
	"user", "username", "password", "passfile", "require_auth",

	// SSL/TLS
	"sslkey", "sslpassword", "sslcert", "sslkeylogfile", "sslrootcert", "sslcrl",

	// OAuth
	"oauth_client_id", "oauth_client_secret",

	// SCRAM
	"scram_client_key", "scram_server_key",

	// Kerberos / GSSAPI
	"krbsrvname", "gsslib", "gssdelegation",
}

// Mask placeholder for sensitive values.
const mask = "xxxxx"

// MaskDSN masks sensitive information in a PostgreSQL Data Source Name (DSN) or connection string.
//
// The function handles two formats:
// 1. Standard URL format: "postgres://user:password@host:port/dbname?param=value"
// 2. PostgreSQL-style key-value connection string: "host=localhost port=5432 user=postgres password=secret"
//
// It ensures that sensitive fields such as passwords, keys, and client secrets are replaced with asterisks
// before being logged or displayed.
func MaskDSN(raw string) string {
	// Attempt to parse the input as a URL.
	u, err := url.Parse(raw)
	if err == nil && u.Scheme != "" {
		// If it is a valid URL with a scheme, mask user info (username/password).
		if u.User != nil {
			if _, hasPassword := u.User.Password(); hasPassword || u.User.Username() != "" {
				u.User = url.UserPassword(mask, mask)
			}
		}

		// Mask any sensitive query parameters.
		q := u.Query()
		for _, key := range sensitiveKeys {
			if _, ok := q[key]; ok {
				q.Set(key, mask)
			}
		}
		u.RawQuery = q.Encode()
		return u.String()
	}

	// If it is not a URL with a scheme, treat it as a PostgreSQL-style key=value connection string.
	return maskConnectionString(raw)
}

// maskConnectionString masks sensitive parameters in a PostgreSQL-style connection string.
//
// Example:
// Input:  "user=postgres password=secret sslkey=mykey.pem"
// Output: "user=***** password=***** sslkey=*****"
func maskConnectionString(connStr string) string {
	for _, key := range sensitiveKeys {
		// Use a case-insensitive regex to match "key=value" patterns and replace the value with mask.
		re := regexp.MustCompile(`(?i)(` + key + `)=([^\s]+)`)
		connStr = re.ReplaceAllString(connStr, "$1="+mask)
	}
	return connStr
}