package common

import "encoding/json"

const (
	// SubjectOperationWildcardData operation.*.data
	SubjectOperationWildcardData = "operation.*.data"

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

	// SubjectOperationRequestData operation.request.data
	SubjectOperationRequestData = "operation.request.data"
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

type PatientEvent struct {
	CurrentID string `json:"current_id"`
	PatientID string `json:"patient_id"`
	//ParentID  string  `json:"parent_id"`
	Step      string  `json:"step"`
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Patient   Patient `json:"patient"`
}

func (p *PatientEvent) String() string {
	d, _ := json.Marshal(p)
	return string(d)
}

func (p *PatientEvent) Bytes() []byte {
	d, _ := json.Marshal(p)
	return d
}

func (p *PatientEvent) Decode(data []byte) error {
	return json.Unmarshal(data, p)
}

type Patient struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Age         int    `json:"age"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	AddressCode string `json:"address_code"`
}
