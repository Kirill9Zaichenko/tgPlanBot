package keyboards

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

func InboxTaskActions(taskID int64) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "✅ Принять",
					CallbackData: fmt.Sprintf("accept:%d", taskID),
				},
				{
					Text:         "❌ Отклонить",
					CallbackData: fmt.Sprintf("reject:%d", taskID),
				},
			},
		},
	}
}
