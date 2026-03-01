package ports

type TokenService interface {
    GenerateAccess(userID string) (string, error)
    GenerateRefresh(userID string) (string, error)
    ValidateToken(token string) (string, error)
}
