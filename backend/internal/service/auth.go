package service

import (
	"context"
	"database/sql"

	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, req internal.LoginRequest) (*internal.LoginResponse, error) {
	user := new(internal.User)
	err := s.DB.NewSelect().
		Model(user).
		Where("username = ?", req.Username).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.UserNotFoundErr
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, internal.WrongPasswordErr
	}

	token, err := s.newToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &internal.LoginResponse{ID: user.ID, Token: token}, nil
}

func (s *service) newToken(id int64) (string, error) {
	claims := &internal.JwtClaims{
		ID: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Secret))
}
