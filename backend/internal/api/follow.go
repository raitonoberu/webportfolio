package api

import (
	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary	Create follow
// @Tags follow
// @Accept json
// @Param request body internal.CreateFollowRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "user not found"
// @Failure 409 {object} errorResponse "user is already followed"
// @Router /follow [post]
// @Security Bearer
func (h *handler) createFollow(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.CreateFollowRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	if data.ID == data.UserID {
		return echo.NewHTTPError(400, "dude wtf")
	}

	ctx := c.Request().Context()
	err := h.CreateFollow(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Get following
// @Tags follow
// @Accept json
// @Produce json
// @Param request query internal.GetFollowingRequest true "query params"
// @Success 200 {object} internal.GetFollowingResponse "follows"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 404 {object} errorResponse "user not found"
// @Router /following [get]
func (h *handler) getFollowing(c echo.Context) error {
	data := internal.GetFollowingRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := h.GetFollowing(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, project)
}

// @Summary Get followers
// @Tags follow
// @Accept json
// @Produce json
// @Param request query internal.GetFollowersRequest true "query params"
// @Success 200 {object} internal.GetFollowersResponse "follows"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 404 {object} errorResponse "user not found"
// @Router /followers [get]
func (h *handler) getFollowers(c echo.Context) error {
	data := internal.GetFollowersRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := h.GetFollowers(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, project)
}

// @Summary Delete follow
// @Tags follow
// @Accept json
// @Param request body internal.DeleteFollowRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "user not found"
// @Failure 409 {object} errorResponse "user is not followed"
// @Router /follow [delete]
// @Security Bearer
func (h *handler) deleteFollow(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.DeleteFollowRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	err := h.DeleteFollow(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
