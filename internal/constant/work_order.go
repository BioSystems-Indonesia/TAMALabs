package constant

type WorkOrderRunAction string

const (
	WorkOrderRunActionRun            WorkOrderRunAction = "run"
	WorkOrderRunActionCancel         WorkOrderRunAction = "cancel"
	WorkOrderRunActionIncompleteSend WorkOrderRunAction = "incomplete_send"
)
