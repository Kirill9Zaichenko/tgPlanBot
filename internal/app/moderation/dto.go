package moderation

type RejectTaskInput struct {
	TaskID         int64
	ReceiverUserID int64
	Comment        string
}
