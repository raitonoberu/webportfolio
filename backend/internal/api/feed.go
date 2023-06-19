package api

import (
	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary Get feed
// @Tags feed
// @Produce json
// @Success 200 {object} internal.GetFeedResponse "feed"
// @Failure 401 {object} errorResponse "not authorized"
// @Router /feed [get]
// @Security Bearer
func (h *handler) getFeed(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.GetFeedRequest{
		UserID: claims.ID,
	}

	ctx := c.Request().Context()
	feed, err := h.GetFeed(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, feed)
}

// @Summary Get trending
// @Tags feed
// @Produce json
// @Success 200 {object} internal.GetTrendingResponse "trending"
// @Router /trending [get]
// @Security Bearer
func (h *handler) getTrending(c echo.Context) error {
	var userID int64
	if user := c.Get("user"); user != nil {
		userID = user.(*jwt.Token).Claims.(*internal.JwtClaims).ID
	}

	data := internal.GetTrendingRequest{
		UserID: userID,
	}

	ctx := c.Request().Context()
	trending, err := h.GetTrending(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, trending)
}
