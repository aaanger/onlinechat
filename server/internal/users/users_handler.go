package users

import (
	"net/http"
	"strconv"

	"onlineChat/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service *UserService
	logger  *logrus.Logger
}

func NewUserHandler(service *UserService, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Debug("Invalid registration request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.service.RegisterUser(req)
	if err != nil {
		h.logger.WithError(err).WithField("email", req.Email).Error("Failed to register user")

		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithField("email", req.Email).Info("Registered new user")
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(c *gin.Context) {
	var req UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Debug("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		h.logger.WithError(err).WithField("email", req.Email).Warn("Failed login attempt")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user profile")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Debug("Invalid profile update request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.UpdateUser(userID, req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to update user profile")

		statusCode := http.StatusInternalServerError
		if err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to delete user account")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
