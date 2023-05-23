package httptransport

import (
	"fmt"
	"net/http"
	"strings"

	"webportfolio/internal"

	_ "webportfolio/docs"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type handler struct {
	internal.Service
}

// @title WebPortfolio API
// @version 1.0
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func Handler(service internal.Service, secret string) *echo.Echo {
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.HTTPErrorHandler = httpErrorHandler
	e.Logger.SetLevel(log.ERROR)

	h := &handler{service}

	authMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(internal.JwtClaims)
		},
		SigningKey: []byte(secret),
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri}" + "\n",
	}))
	e.Use(middleware.Recover())
	// e.Use(middleware.Secure())
	e.Use(contentMiddleware)
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "static",
		HTML5: true,
		Skipper: func(c echo.Context) bool {
			return c.Path() != ""
		},
	}))

	// docs
	e.GET("/api/docs/*", echoSwagger.WrapHandler)

	// auth
	e.POST("/api/login", h.login)

	// user
	e.POST("/api/user", h.createUser)
	e.GET("/api/user", h.getUser)
	e.PATCH("/api/user", h.updateUser, authMiddleware)
	e.DELETE("/api/user", h.deleteUser, authMiddleware)

	// project
	e.POST("/api/project", h.createProject, authMiddleware)
	e.GET("/api/project", h.getProject)
	e.PATCH("/api/project", h.updateProject, authMiddleware)
	e.DELETE("/api/project", h.deleteProject, authMiddleware)

	// upload
	e.POST("/api/upload/avatar", h.uploadAvatar, authMiddleware)
	e.POST("/api/upload/project", h.uploadProject, authMiddleware)

	// avatars
	e.Static("/avatars", "content/avatars")

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}
	return e
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func httpErrorHandler(err error, c echo.Context) {
	var code int
	var msg string
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if he.Internal != nil {
			msg = fmt.Sprintf("%s: %s", he.Message, he.Internal)
		} else {
			msg = fmt.Sprint(he.Message)
		}
	} else {
		code = http.StatusInternalServerError
		msg = err.Error()
	}
	c.Logger().Error(fmt.Sprintf("[%v] %s", code, msg))
	c.JSON(code, errorResponse{Message: msg})
}

func contentMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		host := req.Host
		hostParts := strings.Split(host, ".")
		if len(hostParts) < 3 {
			// no subdomain
			return next(c)
		}

		user := hostParts[0]

		path := req.URL.Path
		pathParts := strings.Split(path, "/") // len >= 2
		if pathParts[1] == "" {
			// empty project name, redirecting to user page
			return c.Redirect(302, fmt.Sprintf("http://web-portfolio.tech/%s", user))
		}

		referer := req.Referer()
		refererParts := strings.Split(referer, "/")
		if len(refererParts) > 3 {
			return c.File(fmt.Sprintf("content/projects/%s/%s%s", user, refererParts[3], req.URL.Path))
		}
		return c.File(fmt.Sprintf("content/projects/%s%s", user, req.URL.Path))
	}
}

type errorResponse struct {
	Message string `json:"message"`
}
