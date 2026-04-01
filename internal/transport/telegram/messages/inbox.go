package messages

import (
	"fmt"
	"strings"

	"tgPlanBot/internal/domain"
)

func InboxItem(item domain.InboxItem) string {
	sender := senderDisplay(item)

	description := strings.TrimSpace(item.Description)
	if description == "" {
		description = "—"
	}

	return fmt.Sprintf(
		"📥 Входящая задача\n\n"+
			"ID: #%d\n"+
			"От: %s\n"+
			"Название: %s\n"+
			"Описание: %s\n"+
			"Статус: %s",
		item.TaskID,
		sender,
		item.Title,
		description,
		item.Status,
	)
}

func senderDisplay(item domain.InboxItem) string {
	if strings.TrimSpace(item.SenderUsername) != "" {
		return "@" + item.SenderUsername
	}

	fullName := strings.TrimSpace(item.SenderFirstName + " " + item.SenderLastName)
	if fullName != "" {
		return fullName
	}

	return fmt.Sprintf("user #%d", item.SenderUserID)
}
