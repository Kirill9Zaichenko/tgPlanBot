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

type OrganizationsHandler struct {
	organizationService *organizationapp.Service
	contextStore        *tgcontext.Store
}

func NewOrganizationsHandler(
	organizationService *organizationapp.Service,
	contextStore *tgcontext.Store,
) *OrganizationsHandler {
	return &OrganizationsHandler{
		organizationService: organizationService,
		contextStore:        contextStore,
	}
}

func (h *OrganizationsHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	items, err := h.organizationService.ListByUserID(ctx, user.ID)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось загрузить организации.",
		})
		log.Printf("list organizations for user %d: %v", user.ID, err)
		return
	}

	activeOrgID, _ := h.contextStore.GetActiveOrganization(user.ID)

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.OrganizationsList(items, activeOrgID),
	})
	if err != nil {
		log.Printf("send /orgs response: %v", err)
	}
}
