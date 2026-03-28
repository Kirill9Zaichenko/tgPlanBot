package messages

import (
	"fmt"
	"strings"

	"tgPlanBot/internal/domain"
)

func InboxList(items []domain.TaskRequest) string {
	if len(items) == 0 {
		return NoInboxItems()
	}

	var sb strings.Builder
	sb.WriteString("Входящие запросы:\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf(
			"Task #%d\nОтправитель: %d\nСтатус: %s\n\n",
			item.TaskID,
			item.SenderUserID,
			item.Status,
		))
	}

	return strings.TrimSpace(sb.String())
}
