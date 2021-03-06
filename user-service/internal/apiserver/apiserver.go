package apiserver

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"example.com/arch/user-service/internal/users"
	"example.com/arch/user-service/internal/users/models"
	"example.com/arch/user-service/internal/users/repository"
	"example.com/arch/user-service/internal/users/service"
	"example.com/arch/user-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

const validationErrorName = "ValidationErrors"

type Server struct {
	userService users.UserService
	router      *gin.Engine
}

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"method", "url"})

func NewServer() *Server {
	db := database.GetPgxPool()
	userRepository := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepository)
	server := &Server{userService: userService}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.RemoveExtraSlash = true

	prometheus := createPrometheus()
	prometheus.Use(router)

	router.GET("/health", server.Health)
	superGroup := router.Group("/api/v1")
	{
		userGroup := superGroup.Group("/user")
		{
			userGroup.POST("", server.CreateUser)
			userGroup.GET(":id", server.GetUser)
			userGroup.PUT(":id", server.UpdateUser)
			userGroup.DELETE(":id", server.DeleteUser)
		}
	}

	server.router = router
	return server
}

func createPrometheus() *ginprometheus.Prometheus {
	prom := ginprometheus.NewPrometheus("api")

	prom.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "id" {
				url = "/api/v1/user/:id"
				break
			}
		}
		return url
	}

	return prom
}

func (s *Server) Start() error {
	return s.router.Run("0.0.0.0:8000")
}

func (s *Server) Health(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusOK)
}

func (s *Server) CreateUser(ctx *gin.Context) {
	timer := prometheus.NewTimer(httpDuration.WithLabelValues("POST", "/api/v1/user/"))
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid body"},
		)
		return
	}

	var h header
	if err := ctx.ShouldBindHeader(&h); err != nil {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"code": 403, "message": "Invalid user id"},
		)
		return
	}

	user := models.User{
		Username:  req.Username,
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		OwnerId:   h.UserId,
	}
	userId, err := s.userService.CreateUser(&user)
	if err != nil {
		if ok := reflect.TypeOf(err).Name() == validationErrorName; ok {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"code": 400, "message": err.Error()},
			)
			return
		}

		fmt.Fprintf(os.Stderr, "Internal server error: %v\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Failed user's creation"},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": userId})
	timer.ObserveDuration()
}

func (s *Server) DeleteUser(ctx *gin.Context) {
	timer := prometheus.NewTimer(httpDuration.WithLabelValues("DELETE", "/api/v1/user/:id"))
	var req deleteUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid user id"},
		)
		return
	}

	var h header
	if err := ctx.ShouldBindHeader(&h); err != nil {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"code": 403, "message": "Invalid user id"},
		)
		return
	}

	err := s.userService.DeleteUser(models.UserId(req.Id), h.UserId)
	if err != nil {
		if err == users.ErrNotFound {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"code": 404, "message": "User not found"},
			)
			return
		}

		fmt.Fprintf(os.Stderr, "Internal server error: %v\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Internal server error"},
		)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
	timer.ObserveDuration()
}

func (s *Server) UpdateUser(ctx *gin.Context) {
	timer := prometheus.NewTimer(httpDuration.WithLabelValues("PUT", "/api/v1/user/:id"))
	var reqBody updateUserRequestBody

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid body content"},
		)
		return
	}

	var reqUri updateUserRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid user's id"},
		)
		return
	}

	var h header
	if err := ctx.ShouldBindHeader(&h); err != nil {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"code": 403, "message": "Invalid user id"},
		)
		return
	}

	user := models.User{
		Id:        models.UserId(reqUri.Id),
		Username:  reqBody.Username,
		Firstname: reqBody.FirstName,
		Lastname:  reqBody.LastName,
		Email:     reqBody.Email,
		Phone:     reqBody.Phone,
		OwnerId:   h.UserId,
	}
	if err := s.userService.UpdateUser(&user); err != nil {
		if ok := reflect.TypeOf(err).Name() == validationErrorName; ok {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"code": 400, "message": err.Error()},
			)
			return
		}

		if err == users.ErrNotFound {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"code": 404, "message": "User profile for this owner isn't found"},
			)
			return
		}

		fmt.Fprintf(os.Stderr, "Internal server error: %v\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Internal server error"},
		)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	timer.ObserveDuration()
}

func (s *Server) GetUser(ctx *gin.Context) {
	timer := prometheus.NewTimer(httpDuration.WithLabelValues("GET", "/api/v1/user/:id"))
	var req getUserRequest
	var h header

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid user id"},
		)
		return
	}

	if err := ctx.ShouldBindHeader(&h); err != nil {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"code": 403, "message": "Invalid user id"},
		)
		return
	}

	user, err := s.userService.GetUser(models.UserId(req.Id))
	if err != nil {
		if err == users.ErrNotFound {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"code": 404, "message": "User not found"},
			)
			return
		}

		fmt.Fprintf(os.Stderr, "Internal server error: %v\n", err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Internal server error"},
		)
		return
	}

	if user.OwnerId != h.UserId {
		ctx.JSON(
			http.StatusForbidden,
			gin.H{"code": 404, "message": "Action is not permitted"},
		)
		return
	}

	response := getUserByIdResponse{
		Id:        int64(user.Id),
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Phone:     user.Phone,
	}

	ctx.JSON(http.StatusOK, response)
	timer.ObserveDuration()
}
