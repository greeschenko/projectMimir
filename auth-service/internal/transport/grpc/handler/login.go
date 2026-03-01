package handler

import (
	"context"

	app "github.com/greeschenko/projectMimir/auth-service/internal/application/user"
	authv1 "github.com/greeschenko/projectMimir/platform/proto/auth/v1"
)

// Login handles the gRPC LoginRequest and returns AuthResponse
func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthResponse, error) {
	// Execute the LoginUseCase with provided email and password
	user, err := h.loginUC.Execute(ctx, app.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// Return AuthResponse with access & refresh tokens and user info
	return &authv1.AuthResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		User: &authv1.User{
			Id:    user.UserID,
			Email: user.Email,
			Role:  "user", // TODO: replace with real role from domain
		},
	}, nil
}
