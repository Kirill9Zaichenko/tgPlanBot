package domain

type InboxItem struct {
	TaskID           int64
	Title            string
	Description      string
	Status           string
	SenderUserID     int64
	SenderTelegramID int64
	SenderUsername   string
	SenderFirstName  string
	SenderLastName   string
}
