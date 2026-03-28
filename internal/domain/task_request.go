package domain

import "time"

type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"
	RequestStatusAccepted RequestStatus = "accepted"
	RequestStatusRejected RequestStatus = "rejected"
)

type TaskRequest struct {
	ID             int64
	TaskID         int64
	SenderUserID   int64
	ReceiverUserID int64
	Status         RequestStatus
	Comment        string
	DecidedAt      *time.Time
	CreatedAt      time.Time
}
