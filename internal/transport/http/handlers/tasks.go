package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	taskapp "tgPlanBot/internal/app/task"

	"github.com/gin-gonic/gin"
)

type createTaskRequest struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	CreatorUserID  int64   `json:"creator_user_id"`
	AssigneeUserID int64   `json:"assignee_user_id"`
	DueAt          *string `json:"due_at"`
}

func GetTasks(taskService *taskapp.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		assigneeUserIDStr := c.Query("assignee_user_id")
		if assigneeUserIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "assignee_user_id is required",
			})
			return
		}

		assigneeUserID, err := strconv.ParseInt(assigneeUserIDStr, 10, 64)
		if err != nil || assigneeUserID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "assignee_user_id must be a positive integer",
			})
			return
		}

		tasks, err := taskService.ListByAssignee(c.Request.Context(), assigneeUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to list tasks",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": tasks,
		})
	}
}

func CreateTask(taskService *taskapp.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		req.Title = strings.TrimSpace(req.Title)
		req.Description = strings.TrimSpace(req.Description)

		var dueAt *time.Time
		if req.DueAt != nil && strings.TrimSpace(*req.DueAt) != "" {
			parsedDueAt, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.DueAt))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "due_at must be RFC3339, example: 2026-03-24T18:00:00Z",
				})
				return
			}
			dueAt = &parsedDueAt
		}

		task, err := taskService.Create(c.Request.Context(), taskapp.CreateTaskInput{
			Title:          req.Title,
			Description:    req.Description,
			CreatorUserID:  req.CreatorUserID,
			AssigneeUserID: req.AssigneeUserID,
			DueAt:          dueAt,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"item": task,
		})
	}
}
