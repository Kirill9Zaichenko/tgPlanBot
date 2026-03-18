package app

import (
	"github.com/go-telegram/bot"
	"tgBotPlan/internal/handlers"
	"tgBotPlan/internal/storage"
)

func Routes(store storage.TaskStore) []bot.Option {
	return []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.Start(store)),
		bot.WithMessageTextHandler("/add", bot.MatchTypePrefix, handlers.Add(store)),
		bot.WithMessageTextHandler("/list", bot.MatchTypeExact, handlers.List(store)),
		bot.WithMessageTextHandler("/done", bot.MatchTypePrefix, handlers.Done(store)),
		bot.WithMessageTextHandler("/del", bot.MatchTypePrefix, handlers.Del(store)),
		bot.WithMessageTextHandler("/clear", bot.MatchTypeExact, handlers.Clear(store)),
	}
}
