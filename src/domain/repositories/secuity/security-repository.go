package repositories

type AuthRepository interface {
	ValidateApiKey(apiKey string, expectedSecret string) (bool, error)
	GetApiKey(apiKey string) (string, error)
	ValidateToken(token string) (bool, error)
}
