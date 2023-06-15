package httptransport

import (
	"net/http"
	"strings"

	"webportfolio/internal"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary Create project
// @Tags project
// @Accept json
// @Param request body internal.CreateProjectRequest true "body params"
// @Success 200 {object} internal.CreateProjectResponse "project"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Router /project [post]
// @Security Bearer
func (h *handler) createProject(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.CreateProjectRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}
	data.Name = strings.ToLower(data.Name)

	ctx := c.Request().Context()
	project, err := h.CreateProject(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, project)
}

// @Summary Get project
// @Tags project
// @Accept json
// @Produce json
// @Param request query internal.GetProjectRequest true "query params"
// @Success 200 {object} internal.GetProjectResponse "project"
// @Failure 400 {object} errorResponse "bad request"
// @Failure 404 {object} errorResponse "user not found"
// @Failure 404 {object} errorResponse "project not found"
// @Router /project [get]
// @Security Bearer
func (h *handler) getProject(c echo.Context) error {
	var userID int64
	if user := c.Get("user"); user != nil {
		userID = user.(*jwt.Token).Claims.(*internal.JwtClaims).ID
	}

	data := internal.GetProjectRequest{
		ReqUserID: userID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}

	if data.ID == nil && (data.Name == nil || (data.UserID == nil && data.Username == nil)) {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	project, err := h.GetProject(ctx, data)
	if err != nil {
		return err
	}
	return c.JSON(200, project)
}

// @Summary Update project
// @Tags project
// @Accept json
// @Param request body internal.UpdateProjectRequest true "body params"
// @Success 204
// @Failure 400 {object} errorResponse "bad request"
// @Failure 400 {object} errorResponse "validation failed"
// @Failure 401 {object} errorResponse "not authorized"
// @Router /project [patch]
// @Security Bearer
func (h *handler) updateProject(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.UpdateProjectRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	if data.Description == nil && data.Readme == nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	err := h.UpdateProject(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Delete project
// @Tags project
// @Accept json
// @Param request body internal.DeleteProjectRequest true "body params"
// @Success 204
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "project not found"
// @Router /project [delete]
// @Security Bearer
func (h *handler) deleteProject(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.DeleteProjectRequest{
		UserID: claims.ID,
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	ctx := c.Request().Context()
	err := h.DeleteProject(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Summary Upload project
// @Tags project
// @Accept x-www-form-urlencoded
// @Param file formData file true "zip-archive"
// @Param id formData int true "project id"
// @Success 204
// @Failure 400 {object} errorResponse "file is too big"
// @Failure 401 {object} errorResponse "not authorized"
// @Failure 404 {object} errorResponse "project not found"
// @Router /upload [post]
// @Security Bearer
func (h *handler) uploadProject(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*internal.JwtClaims)

	data := internal.UploadProjectRequest{
		UserID: claims.ID,
	}
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	data.File = file
	if err := c.Bind(&data); err != nil {
		return err
	}
	if err := c.Validate(&data); err != nil {
		return err
	}

	if file.Size > 1024*1024*5 {
		return internal.FileTooBigErr
	}

	ctx := c.Request().Context()
	err = h.UploadProject(ctx, data)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
