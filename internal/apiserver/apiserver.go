package apiserver

import (
	"log"
	"net/http"
	"reflect"

	"example.com/arch/user-service/internal/users"
	"example.com/arch/user-service/internal/users/models"
	"example.com/arch/user-service/internal/users/repository"
	"example.com/arch/user-service/internal/users/service"
	"example.com/arch/user-service/pkg/database"
	"github.com/gin-gonic/gin"
)

const validationErrorName = "ValidationErrors"

type Server struct {
	userService users.UserService
	router      *gin.Engine
}

func NewServer() *Server {
	db := database.GetPgxPool()
	userRepository := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepository)
	server := &Server{userService: userService}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.RemoveExtraSlash = true

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

func (s *Server) Start() error {
	return s.router.Run("0.0.0.0:8000")
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid body"},
		)
		return
	}

	user := models.User{
		Username:  req.Username,
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
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
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Failed user's creation"},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": userId})
}

func (s *Server) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid user id"},
		)
		return
	}

	err := s.userService.DeleteUser(models.UserId(req.Id))
	if err != nil {
		if err == users.ErrNotFound {
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"code": 404, "message": "User not found"},
			)
			return
		}

		log.Println(err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Internal server error"},
		)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) UpdateUser(ctx *gin.Context) {
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

	user := models.User{
		Id:        models.UserId(reqUri.Id),
		Username:  reqBody.Username,
		Firstname: reqBody.FirstName,
		Lastname:  reqBody.LastName,
		Email:     reqBody.Email,
		Phone:     reqBody.Phone,
	}
	if err := s.userService.UpdateUser(&user); err != nil {
		if ok := reflect.TypeOf(err).Name() == validationErrorName; ok {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"code": 400, "message": err.Error()},
			)
			return
		}

		log.Println(err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Internal server error"},
		)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}

func (s *Server) GetUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"code": 400, "message": "Invalid user id"},
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

		log.Println(err)
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"code": 500, "message": "Internal server error"},
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
}
