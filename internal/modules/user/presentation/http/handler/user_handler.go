package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/application/dto"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/application/service"
	"github.com/ndxbinh1922001/VNalo-be/pkg/response"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with email, password, and username
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Create user request"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID (soft delete)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=dto.UserListResponse}
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users", err)
		return
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

// PromoteToVIP godoc
// @Summary Promote user to VIP
// @Description Promote user to VIP status
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/promote-vip [post]
func (h *UserHandler) PromoteToVIP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.PromoteUserToVIP(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to promote user to VIP", err)
		return
	}

	response.Success(c, http.StatusOK, "User promoted to VIP successfully", nil)
}

// DemoteFromVIP godoc
// @Summary Demote user from VIP
// @Description Remove VIP status from user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/demote-vip [post]
func (h *UserHandler) DemoteFromVIP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.DemoteUserFromVIP(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to demote user from VIP", err)
		return
	}

	response.Success(c, http.StatusOK, "User demoted from VIP successfully", nil)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user password
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to change password", err)
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

// ActivateUser godoc
// @Summary Activate user
// @Description Activate user account
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.ActivateUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to activate user", err)
		return
	}

	response.Success(c, http.StatusOK, "User activated successfully", nil)
}

// DeactivateUser godoc
// @Summary Deactivate user
// @Description Deactivate user account
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.DeactivateUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to deactivate user", err)
		return
	}

	response.Success(c, http.StatusOK, "User deactivated successfully", nil)
}

