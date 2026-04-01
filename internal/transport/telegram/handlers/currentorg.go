package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	organizationapp "tgPlanBot/internal/app/organization"
	"tgPlanBot/internal/domain"
	tgcontext "tgPlanBot/internal/transport/telegram/context"
	"tgPlanBot/internal/transport/telegram/messages"
)

type CurrentOrganizationHandler struct {
	organizationService *organizationapp.Service
	contextStore        *tgcontext.Store
}

func NewCurrentOrganizationHandler(
	organizationService *organizationapp.Service,
	contextStore *tgcontext.Store,
) *CurrentOrganizationHandler {
	return &CurrentOrganizationHandler{
		organizationService: organizationService,
		contextStore:        contextStore,
	}
}

func (h *CurrentOrganizationHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	orgID, ok := h.contextStore.GetActiveOrganization(user.ID)
	if !ok || orgID <= 0 {
		_, err := bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.CurrentOrganization(nil),
		})
		if err != nil {
			log.Printf("send /currentorg empty response: %v", err)
		}
		return
	}

	org, err := h.organizationService.GetByIDForUser(ctx, orgID, user.ID)
	if err != nil {
		h.contextStore.ClearActiveOrganization(user.ID)

		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.CurrentOrganization(nil),
		})
		log.Printf("get current org %d for user %d: %v", orgID, user.ID, err)
		return
	}

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.CurrentOrganization(org),
	})
	if err != nil {
		log.Printf("send /currentorg response: %v", err)
	}
}
