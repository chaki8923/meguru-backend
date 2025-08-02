package controller

import (
	"net/http"

	"meguru-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PushSubscriptionController struct {
	pushSubscriptionUsecase *usecase.PushSubscriptionUsecase
}

func NewPushSubscriptionController(pushSubscriptionUsecase *usecase.PushSubscriptionUsecase) *PushSubscriptionController {
	return &PushSubscriptionController{
		pushSubscriptionUsecase: pushSubscriptionUsecase,
	}
}

func (psc *PushSubscriptionController) Subscribe(c *gin.Context) {
	var req usecase.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := psc.pushSubscriptionUsecase.CreateSubscription(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (psc *PushSubscriptionController) SendNotification(c *gin.Context) {
	var req usecase.SendNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := psc.pushSubscriptionUsecase.SendNotificationToAll(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications sent successfully"})
}
