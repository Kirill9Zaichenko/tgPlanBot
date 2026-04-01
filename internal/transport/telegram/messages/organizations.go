package messages

import (
	"fmt"
	"strings"

	"tgPlanBot/internal/domain"
)

func OrganizationsList(items []domain.Organization, activeOrganizationID int64) string {
	if len(items) == 0 {
		return "Ты пока не состоишь ни в одной организации."
	}

	var sb strings.Builder
	sb.WriteString("Твои организации:\n\n")

	for _, item := range items {
		marker := ""
		if item.ID == activeOrganizationID {
			marker = " ✅ active"
		}

		sb.WriteString(fmt.Sprintf(
			"#%d | %s (%s)%s\n",
			item.ID,
			item.Name,
			item.Slug,
			marker,
		))
	}

	sb.WriteString("\nИспользуй /useorg {id}, чтобы переключиться.")
	return strings.TrimSpace(sb.String())
}

func CurrentOrganization(item *domain.Organization) string {
	if item == nil {
		return "Активная организация не выбрана."
	}

	return fmt.Sprintf(
		"Текущая организация:\n\n#%d | %s (%s)",
		item.ID,
		item.Name,
		item.Slug,
	)
}

func ActiveOrganizationChanged(item *domain.Organization) string {
	return fmt.Sprintf(
		"Активная организация переключена.\n\n#%d | %s (%s)",
		item.ID,
		item.Name,
		item.Slug,
	)
}

func UseOrgUsage() string {
	return "Использование: /useorg {organization_id}"
}

func InvalidOrganizationID() string {
	return "Некорректный organization_id."
}

func OrganizationNotFoundOrForbidden() string {
	return "Организация не найдена или ты не состоишь в ней."
}
