package internal

import (
	"github.com/labstack/echo/v4"
)

var (
	FileTooBigErr      = echo.NewHTTPError(400, "file is too big")
	WrongPasswordErr   = echo.NewHTTPError(403, "wrong password")
	UserNotFoundErr    = echo.NewHTTPError(404, "user not found")
	ProjectNotFoundErr = echo.NewHTTPError(404, "project not found")
	CommentNotFoundErr = echo.NewHTTPError(404, "comment not found")
	UsernameExistsErr  = echo.NewHTTPError(409, "username already exists")
	EmailExistsErr     = echo.NewHTTPError(409, "email already exists")
	ProjectExistsErr   = echo.NewHTTPError(409, "project already exists")
	ProjectLikedErr    = echo.NewHTTPError(409, "project is already liked")
	ProjectNotLikedErr = echo.NewHTTPError(409, "project is not liked")
)
