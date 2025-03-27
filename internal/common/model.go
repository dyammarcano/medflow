package common

const (
	// SubjectOperationWildcardData operation.*.data
	SubjectOperationWildcardData = "operation.%s.data"

	// SubjectOperationResponseData operation.response.data
	SubjectOperationResponseData = "operation.response.data"

	// SubjectOperationIncomingData operation.incoming.data
	SubjectOperationIncomingData = "operation.incoming.data"
	// SubjectOperationErrorData operation.error.dada
	SubjectOperationErrorData = "operation.error.data"

	// SubjectOperationStage1Data operation.stage1.data
	SubjectOperationStage1Data = "operation.stage1.data"

	// SubjectOperationStage2Data operation.stage2.data
	SubjectOperationStage2Data = "operation.stage2.data"

	// SubjectOperationStage3Data operation.stage3.data
	SubjectOperationStage3Data = "operation.stage3.data"

	// SubjectOperationExamsData operation.exams.data
	SubjectOperationExamsData = "operation.exams.data"

	// SubjectOperationPriority1Data operation.priority1.data
	SubjectOperationPriority1Data = "operation.priority1.data"
)

type ClinicalEvent struct {
	ParentID  string            `json:"parent_id"`
	CurrentID string            `json:"current_id"`
	Step      string            `json:"step"`
	Timestamp string            `json:"timestamp"`
	Status    string            `json:"status"`
	Data      map[string]any    `json:"data"`
	Metadata  map[string]string `json:"metadata"`
}

type ExamEvent struct {
	ParentID  string            `json:"parent_id"`
	ExamID    string            `json:"exam_id"`
	ExamType  string            `json:"exam_type"`
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}
