package httptransport

import (
	"os"
	"path/filepath"
	"strconv"
	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var emptyAvatarPath = filepath.Join("static", "avatar.png")

// @Summary Create avatar
// @Tags avatar
// @Accept x-www-form-urlencoded
// @Param file formData file true "image"
// @Success 204
// @Failure 401 {object} errorResponse "not authorized"
// @Router /avatar [post]
// @Security Bearer
func (h *handler) createAvatar(c echo.Context) error {
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
	err = h.CreateAvatar(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Get avatar
// @Tags avatar
// @Accept json
// @Produce png
// @Param request query internal.GetAvatarRequest true "query params"
// @Success 200
// @Failure 400 {object} errorResponse "validation failed"
// @Router /avatar [get]
func (h *handler) getAvatar(c echo.Context) error {
	data := internal.GetAvatarRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	path := filepath.Join("content", "avatars", strconv.FormatInt(data.ID, 10))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.File(emptyAvatarPath)
	}
	return c.File(path)
}

// @Summary Delete avatar
// @Tags avatar
// @Success 204
// @Router /avatar [delete]
func (h *handler) deleteAvatar(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.DeleteAvatarRequest{
		UserID: claims.ID,
	}

	ctx := c.Request().Context()
	err := h.DeleteAvatar(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
