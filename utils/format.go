package utils

const (
	AuthHeader     = "Authorization"
	ClientIdHeader = "Client-Id"
)

func FormatToken(token string) string {
	return "Bearer " + token
}
