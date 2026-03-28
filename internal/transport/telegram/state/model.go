package state

type Step string

const (
	StepIdle               Step = ""
	StepWaitingTaskTitle   Step = "waiting_task_title"
	StepWaitingDescription Step = "waiting_description"
)

type NewTaskState struct {
	UserID      int64
	Step        Step
	Title       string
	Description string
}
