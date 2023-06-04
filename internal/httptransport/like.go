package httptransport

import (
	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary	Create like
// @Tags like
// @Accept json
// @Param request body internal.CreateLikeRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "project not found"
// @Failure 409 {object} errorResponse "project is already liked"
// @Router /like [post]
// @Security Bearer
func (h *handler) createLike(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	data := internal.CreateLikeRequest{
		UserID: int64(claims["id"].(float64)),
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	err := h.CreateLike(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Delete like
// @Tags like
// @Accept json
// @Param request body internal.DeleteLikeRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "project not found"
// @Failure 409 {object} errorResponse "project is not liked"
// @Router /like [delete]
// @Security Bearer
func (h *handler) deleteLike(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	data := internal.DeleteLikeRequest{
		UserID: int64(claims["id"].(float64)),
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	err := h.DeleteLike(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
