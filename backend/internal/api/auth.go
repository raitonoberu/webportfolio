package api

import (
	"webportfolio/internal"

	"github.com/labstack/echo/v4"
)

// @Summary Login to account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body internal.LoginRequest true "body params"
// @Success 200 {object} internal.LoginResponse
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 403 {object} errorResponse "wrong password"
// @Failure 404 {object} errorResponse "user not found"
// @Router /login [post]
func (h *handler) login(c echo.Context) error {
	data := internal.LoginRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	resp, err := h.Login(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, resp)
}
