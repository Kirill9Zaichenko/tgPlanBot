package state

type Step string

const (
	StepIdle                      Step = ""
	StepWaitingTaskTitle          Step = "waiting_task_title"
	StepWaitingDescription        Step = "waiting_description"
	StepWaitingAssigneeTelegramID Step = "waiting_assignee_telegram_id"
)

type Flow string

const (
	FlowNewTask    Flow = "newtask"
	FlowNewTaskFor Flow = "newtaskfor"
)

type NewTaskState struct {
	UserID             int64
	Flow               Flow
	Step               Step
	AssigneeTelegramID int64
	AssigneeUserID     int64
	Title              string
	Description        string
}
