package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	id, ok := userID.(int)
	if !ok {
		return 0, fmt.Errorf("user ID is not of type int")
	}

	return id, nil
}

func GetUsername(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", fmt.Errorf("username not found in context")
	}

	name, ok := username.(string)
	if !ok {
		return "", fmt.Errorf("username is not of type string")
	}

	return name, nil
}






