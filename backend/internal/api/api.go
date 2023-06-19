package api

import (
	"fmt"
	"net/http"
	"strings"

	"webportfolio/internal"

	_ "webportfolio/docs"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
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
			return &internal.JwtClaims{}
		},
		SigningKey: []byte(secret),

		ErrorHandler: func(c echo.Context, err error) error {
			if c.Request().Method != "GET" {
				return internal.NotAuthorizedErr
			}
			return nil
		},
		ContinueOnIgnoredError: true,
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri}" + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
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
	e.GET("/api/user", h.getUser, authMiddleware) // optional auth
	e.PATCH("/api/user", h.updateUser, authMiddleware)
	e.DELETE("/api/user", h.deleteUser, authMiddleware)

	// project
	e.POST("/api/project", h.createProject, authMiddleware)
	e.GET("/api/project", h.getProject, authMiddleware) // optional auth
	e.PATCH("/api/project", h.updateProject, authMiddleware)
	e.DELETE("/api/project", h.deleteProject, authMiddleware)
	e.POST("/api/upload", h.uploadProject, authMiddleware)

	// like
	e.POST("/api/like", h.createLike, authMiddleware)
	e.DELETE("/api/like", h.deleteLike, authMiddleware)

	// comment
	e.POST("/api/comment", h.createComment, authMiddleware)
	e.GET("/api/comment", h.getComments)
	e.DELETE("/api/comment", h.deleteComment, authMiddleware)

	// follow
	e.POST("/api/follow", h.createFollow, authMiddleware)
	e.GET("/api/following", h.getFollowing)
	e.GET("/api/followers", h.getFollowers)
	e.DELETE("/api/follow", h.deleteFollow, authMiddleware)

	// feed
	e.GET("/api/feed", h.getFeed, authMiddleware)
	e.GET("/api/trending", h.getTrending, authMiddleware)

	// avatar
	e.POST("/api/avatar", h.createAvatar, authMiddleware)
	e.GET("/api/avatar", h.getAvatar)
	e.DELETE("/api/avatar", h.deleteAvatar, authMiddleware)

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

		return c.File(fmt.Sprintf("content/projects/%s%s", hostParts[0], req.URL.Path))
	}
}

type errorResponse struct {
	Message string `json:"message"`
}
