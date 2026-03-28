package messages

import (
	"fmt"
	"strings"

	"tgPlanBot/internal/domain"
)

func TasksList(tasks []domain.Task) string {
	if len(tasks) == 0 {
		return NoTasks()
	}

	var sb strings.Builder
	sb.WriteString("Твои задачи:\n\n")

	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf(
			"#%d | %s\nСтатус: %s\n\n",
			task.ID,
			task.Title,
			task.Status,
		))
	}

	return strings.TrimSpace(sb.String())
}
