package httptransport

import (
	"net/http"

	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary Create user
// @Tags user
// @Accept json
// @Param request body internal.CreateUserRequest true "body params"
// @Success 201 {object} internal.CreateUserResponse "user"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 409 {object} errorResponse "username already exists"
// @Failure 409 {object} errorResponse "email already exists"
// @Router /user [post]
func (h *handler) createUser(c echo.Context) error {
	data := internal.CreateUserRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	user, err := h.CreateUser(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(201, user)
}

// @Summary Get user
// @Tags user
// @Accept json
// @Produce json
// @Param request query internal.GetUserRequest true "query params"
// @Success 200 {object} internal.GetUserResponse "user"
// @Failure 400 {object} errorResponse "bad request"
// @Failure 404 {object} errorResponse "user not found"
// @Router /user [get]
func (h *handler) getUser(c echo.Context) error {
	data := internal.GetUserRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}

	if data.ID == nil && data.Name == nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	user, err := h.GetUser(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, user)
}

// @Summary Update user
// @Tags user
// @Accept json
// @Param request body internal.UpdateUserRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "bad request"
// @Failure 401 {object} errorResponse "not authorized"
// @Router /user [patch]
// @Security Bearer
func (h *handler) updateUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.UpdateUserRequest{
		ID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}

	if data.Bio == nil && data.Fullname == nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	err := h.UpdateUser(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Delete user
// @Tags user
// @Success 204
// @Failure 401 {object} errorResponse "not authorized"
// @Router /user [delete]
// @Security Bearer
func (h *handler) deleteUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.DeleteUserRequest{
		ID: claims.ID,
	}

	ctx := c.Request().Context()
	err := h.DeleteUser(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Upload avatar
// @Tags upload
// @Accept x-www-form-urlencoded
// @Param file formData file true "image"
// @Success 204
// @Failure 401 {object} errorResponse "not authorized"
// @Router /upload/avatar [post]
// @Security Bearer
func (h *handler) uploadAvatar(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.UploadAvatarRequest{
		UserID: claims.ID,
	}
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	data.File = file

	ctx := c.Request().Context()
	err = h.UploadAvatar(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
