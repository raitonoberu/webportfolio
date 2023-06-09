package httptransport

import (
	"os"
	"path/filepath"
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
// @Failure 400 {object} errorResponse "missing id param"
// @Router /avatar [get]
func (h *handler) getAvatar(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return echo.NewHTTPError(400, "missing id param")
	}

	path := filepath.Join("content", "avatars", id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.File(emptyAvatarPath)
	}
	return c.File(path)
}
