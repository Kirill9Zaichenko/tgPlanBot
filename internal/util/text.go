package util

import (
	"fmt"
	"strings"
	"tgBotPlan/internal/domain"
)

func HelpText() string {
	return strings.TrimSpace(`
		Я бот задач ✅
	
		Команды:
			/add <текст>   — добавить задачу
			/list          — список задач
			/done <id>     — отметить выполненной
			/del <id>      — удалить задачу
			/clear         — очистить список
	`)
}

func RenderList(tasks []*domain.Task) string {
	if len(tasks) == 0 {
		return "Список пуст. Добавьте задачу"
	}
	var sb strings.Builder
	sb.WriteString("Задачи:\n")
	for _, t := range tasks {
		status := "⬜️"
		if t.Done {
			status = "✅"
		}
		sb.WriteString(fmt.Sprintf("%s %d) %s\n", status, t.ID, t.Text))
	}
	return strings.TrimSpace(sb.String())
}
