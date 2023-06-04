package service

import (
	"context"
	"database/sql"
	"time"

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

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return nil, err
	}
	return &internal.LoginResponse{ID: user.ID, Token: t}, nil
}
