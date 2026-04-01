package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	organizationapp "tgPlanBot/internal/app/organization"
	"tgPlanBot/internal/domain"
	tgcontext "tgPlanBot/internal/transport/telegram/context"
	"tgPlanBot/internal/transport/telegram/messages"
)

type UseOrganizationHandler struct {
	organizationService *organizationapp.Service
	contextStore        *tgcontext.Store
}

func NewUseOrganizationHandler(
	organizationService *organizationapp.Service,
	contextStore *tgcontext.Store,
) *UseOrganizationHandler {
	return &UseOrganizationHandler{
		organizationService: organizationService,
		contextStore:        contextStore,
	}
}

func (h *UseOrganizationHandler) Handle(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	user *domain.User,
) {
	if update.Message == nil {
		return
	}

	parts := strings.Fields(strings.TrimSpace(update.Message.Text))
	if len(parts) < 2 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.UseOrgUsage(),
		})
		return
	}

	orgID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || orgID <= 0 {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.InvalidOrganizationID(),
		})
		return
	}

	org, err := h.organizationService.GetByIDForUser(ctx, orgID, user.ID)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.OrganizationNotFoundOrForbidden(),
		})
		log.Printf("use org %d for user %d: %v", orgID, user.ID, err)
		return
	}

	h.contextStore.SetActiveOrganization(user.ID, org.ID)

	_, err = bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.ActiveOrganizationChanged(org),
	})
	if err != nil {
		log.Printf("send /useorg response: %v", err)
	}
}
