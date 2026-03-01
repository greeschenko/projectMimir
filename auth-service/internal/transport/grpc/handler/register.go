package handler

import (
	"context"

	app "github.com/greeschenko/projectMimir/auth-service/internal/application/user"
	authv1 "github.com/greeschenko/projectMimir/platform/proto/auth/v1"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	registerUC *app.RegisterUseCase
	loginUC    *app.LoginUseCase
}

func NewAuthHandler(registerUC *app.RegisterUseCase, loginUC *app.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {
	user, err := h.registerUC.Execute(ctx, app.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &authv1.AuthResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		User: &authv1.User{
			Id:    user.UserID,
			Email: user.Email,
			Role:  "user", //TODO hardcode tmp
		},
	}, nil
}
