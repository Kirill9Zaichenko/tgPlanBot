package handlers

import (
	"net/http"
	"strconv"

	moderationapp "tgPlanBot/internal/app/moderation"

	"github.com/gin-gonic/gin"
)

type rejectTaskRequest struct {
	ReceiverUserID int64  `json:"receiver_user_id"`
	Comment        string `json:"comment"`
}

func GetInbox(service *moderationapp.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		receiverUserIDStr := c.Query("receiver_user_id")
		if receiverUserIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "receiver_user_id is required",
			})
			return
		}

		receiverUserID, err := strconv.ParseInt(receiverUserIDStr, 10, 64)
		if err != nil || receiverUserID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "receiver_user_id must be a positive integer",
			})
			return
		}

		items, err := service.ListInbox(c.Request.Context(), receiverUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to load inbox",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": items,
		})
	}
}

func AcceptTask(service *moderationapp.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskIDStr := c.Param("id")
		taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
		if err != nil || taskID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
			return
		}

		receiverUserIDStr := c.Query("receiver_user_id")
		receiverUserID, err := strconv.ParseInt(receiverUserIDStr, 10, 64)
		if err != nil || receiverUserID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "receiver_user_id must be a positive integer",
			})
			return
		}

		if err := service.AcceptTask(c.Request.Context(), taskID, receiverUserID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "accepted",
		})
	}
}

func RejectTask(service *moderationapp.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskIDStr := c.Param("id")
		taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
		if err != nil || taskID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid task id",
			})
			return
		}

		var req rejectTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		if req.ReceiverUserID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "receiver_user_id must be a positive integer",
			})
			return
		}

		if err := service.RejectTask(c.Request.Context(), taskID, req.ReceiverUserID, req.Comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "rejected",
		})
	}
}
