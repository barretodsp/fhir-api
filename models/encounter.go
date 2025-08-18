package models

import "time"

type Period struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end,omitempty" json:"end,omitempty"`
}

type Encounter struct {
	FhirId         string `bson:"fhirId" json:"fhirId"`
	FullUrl        string `bson:"fullUrl" json:"fullUrl"`
	Status         string `bson:"status" json:"status"`
	Class          string `bson:"class" json:"class"`
	Period         Period `bson:"period" json:"period"`
	PractitionerID string `bson:"practitionerId,omitempty" json:"practitionerId,omitempty"`
	PatientID      string `bson:"patientId,omitempty" json:"patientId,omitempty"`
}

type EncounterUpdate struct {
	Status string `json:"status" binding:"required"`
}

type EncounterResponse struct {
	FhirId         *string `json:"fhirId,omitempty"`
	FullURL        *string `json:"fullUrl,omitempty"`
	Status         *string `json:"status,omitempty"`
	Class          *string `json:"class,omitempty"`
	Period         *Period `json:"period,omitempty"`
	PractitionerID *string `json:"practitionerId,omitempty"`
	PatientID      *string `json:"patientId,omitempty"`
}
