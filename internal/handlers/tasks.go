package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
	"tgBotPlan/internal/storage"
	"tgBotPlan/internal/telegram"
	"tgBotPlan/internal/util"
)

func Add(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		text := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/add"))
		if text == "" {
			telegram.Reply(ctx, b, chatID, "Укажите текст задачи: /add сделать отчет")
			return
		}
		t := store.Add(chatID, text)
		telegram.Reply(ctx, b, chatID, fmt.Sprintf("Добавлено: %d) %s", t.ID, t.Text))
	}
}

func List(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		telegram.Reply(ctx, b, chatID, util.RenderList(store.List(chatID)))
	}
}

func Done(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		arg := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/done"))
		id, err := strconv.Atoi(arg)
		if err != nil || id <= 0 {
			telegram.Reply(ctx, b, chatID, "Укажи id: /done 2")
			return
		}
		if store.Done(chatID, id) {
			telegram.Reply(ctx, b, chatID, "Готово ✅")
		} else {
			telegram.Reply(ctx, b, chatID, "Не нашёл задачу с таким id")
		}
	}
}

func Del(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		arg := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/del"))
		id, err := strconv.Atoi(arg)
		if err != nil || id <= 0 {
			telegram.Reply(ctx, b, chatID, "Укажи id: /del 2")
			return
		}
		if store.Delete(chatID, id) {
			telegram.Reply(ctx, b, chatID, "Удалено 🗑️")
		} else {
			telegram.Reply(ctx, b, chatID, "Не нашёл задачу с таким id")
		}
	}
}

func Clear(store storage.TaskStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		store.Clear(chatID)
		telegram.Reply(ctx, b, chatID, "Список очищен.")
	}
}
