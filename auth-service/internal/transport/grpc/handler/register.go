package handler

import (
	"context"

	app "github.com/greeschenko/projectMimir/auth-service/internal/application/user"
	authv1 "github.com/greeschenko/projectMimir/platform/proto/auth/v1"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	registerUC *app.RegisterUseCase
}

func NewAuthHandler(registerUC *app.RegisterUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
	}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.AuthResponse, error) {

	user, err := h.registerUC.Execute(ctx, app.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &authv1.AuthResponse{
		AccessToken:  "",
		RefreshToken: "",
		User: &authv1.User{
			Id:    user.ID().String(),
			Email: user.Email(),
			Role:  "user", //TODO hardcode tmp
		},
	}, nil
}
