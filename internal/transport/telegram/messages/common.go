package messages

import "fmt"

func Start() string {
	return "Привет! Я бот-планировщик задач.\n\n" +
		"Команды:\n" +
		"/help — помощь\n" +
		"/me — мой Telegram ID\n" +
		"/mytasks — мои задачи\n" +
		"/inbox — входящие запросы\n" +
		"/newtask — создать задачу себе\n" +
		"/accept {task_id} — принять задачу\n" +
		"/reject {task_id} {comment} — отклонить задачу"
}

func Help() string {
	return "Доступные команды:\n\n" +
		"/start — запустить бота\n" +
		"/help — показать помощь\n" +
		"/me — показать мой Telegram ID\n" +
		"/mytasks — показать мои задачи\n" +
		"/inbox — показать входящие запросы\n" +
		"/newtask — создать задачу себе\n" +
		"/accept {task_id} — принять задачу\n" +
		"/reject {task_id} {comment} — отклонить задачу"
}

func Me(userID int64, username, firstName, lastName string) string {
	usernameText := "-"
	if username != "" {
		usernameText = "@" + username
	}

	fullName := firstName
	if lastName != "" {
		fullName += " " + lastName
	}

	return fmt.Sprintf(
		"👤 Твои данные:\n\nID: %d\nUsername: %s\nИмя: %s",
		userID,
		usernameText,
		fullName,
	)
}

func UsageAccept() string {
	return "Использование: /accept {task_id}"
}

func UsageReject() string {
	return "Использование: /reject {task_id} {comment}"
}

func InvalidTaskID() string {
	return "Некорректный task_id."
}

func TasksLoadFailed() string {
	return "Не удалось загрузить задачи."
}

func InboxLoadFailed() string {
	return "Не удалось загрузить входящие запросы."
}

func NoTasks() string {
	return "У тебя пока нет задач."
}

func NoInboxItems() string {
	return "Входящих запросов нет."
}

func TaskAccepted() string {
	return "Задача успешно принята."
}

func TaskRejected() string {
	return "Задача отклонена."
}

func AcceptFailed(err error) string {
	return fmt.Sprintf("Не удалось принять задачу: %s", err.Error())
}

func RejectFailed(err error) string {
	return fmt.Sprintf("Не удалось отклонить задачу: %s", err.Error())
}
