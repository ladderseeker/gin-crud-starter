package v1

import (
	stderrors "errors"
	"github.com/ladderseeker/gin-crud-starter/internal/model"
	"github.com/ladderseeker/gin-crud-starter/internal/service"
	apperrors "github.com/ladderseeker/gin-crud-starter/pkg/errors"
	"github.com/ladderseeker/gin-crud-starter/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService service.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register registers the router for the user controller
func (c *UserController) Register(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("", c.GetAllUsers)
		users.GET("/:id", c.GetUserByID)
		users.POST("", c.CreateUser)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
	}
}

// GetAllUsers returns all users
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} entities.UserResponse
// @Failure 500 {object} errors.AppError
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserByID returns a user by ID
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.UserResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apperrors.NewInvalidInputError("Invalid ID format", nil, err))
		return
	}

	user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.UserCreate true "User object"
// @Success 201 {object} entities.UserResponse
// @Failure 400 {object} errors.AppError
// @Failure 409 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var input model.UserCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logger.Error("Invalid input for creating user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, apperrors.NewInvalidInputError("Invalid input", nil, err))
		return
	}

	user, err := c.userService.CreateUser(ctx.Request.Context(), input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// UpdateUser updates a user
// @Summary Update a user
// @Description Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body entities.UserUpdate true "User object"
// @Success 200 {object} entities.UserResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apperrors.NewInvalidInputError("Invalid ID format", nil, err))
		return
	}

	var input model.UserUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logger.Error("Invalid input for updating user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, apperrors.NewInvalidInputError("Invalid input", nil, err))
		return
	}

	user, err := c.userService.UpdateUser(ctx.Request.Context(), id, input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {object} nil
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apperrors.NewInvalidInputError("Invalid ID format", nil, err))
		return
	}

	if err := c.userService.DeleteUser(ctx.Request.Context(), id); err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Helper function to parse ID parameter
func parseIDParam(ctx *gin.Context) (uint, error) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// Helper function to handle errors
func handleError(ctx *gin.Context, err error) {
	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		ctx.JSON(appErr.StatusCode, appErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, apperrors.NewInternalError("An unexpected error occurred", err))
}
