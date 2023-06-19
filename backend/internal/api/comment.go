package api

import (
	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary	Create comment
// @Tags comment
// @Accept json
// @Produce json
// @Param request body internal.CreateCommentRequest true "body params"
// @Success 200 {object} internal.CreateCommentResponse "comment"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "project not found"
// @Router /comment [post]
// @Security Bearer
func (h *handler) createComment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.CreateCommentRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	comment, err := h.CreateComment(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, comment)
}

// @Summary Get comments
// @Tags comment
// @Accept json
// @Produce json
// @Param request query internal.GetCommentsRequest true "query params"
// @Success 200 {object} internal.GetCommentsResponse "comments"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 404 {object} errorResponse "project not found"
// @Router /comment [get]
func (h *handler) getComments(c echo.Context) error {
	data := internal.GetCommentsRequest{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	project, err := h.GetComments(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, project)
}

// @Summary Delete comment
// @Tags comment
// @Accept json
// @Param request body internal.DeleteCommentRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "project not found"
// @Failure 404 {object} errorResponse "comment not found"
// @Router /comment [delete]
// @Security Bearer
func (h *handler) deleteComment(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.DeleteCommentRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	err := h.DeleteComment(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
