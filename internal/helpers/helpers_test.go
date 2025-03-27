package helpers

import "testing"

func TestGeneratePatient(t *testing.T) {
	patientEvent := GeneratePatient()
	if patientEvent.PatientID == "" {
		t.Errorf("PatientID is empty")
	}

	if patientEvent.CurrentID == "" {
		t.Errorf("CurrentID is empty")
	}

	t.Log(patientEvent.String())
}
