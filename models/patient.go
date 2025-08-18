package models

type Patient struct {
	FhirId     string `bson:"fhirId" json:"fhirId"`
	GivenName  string `bson:"givenName" json:"givenName"`
	FamilyName string `bson:"familyName" json:"familyName"`
	BirthDate  string `bson:"birthDate" json:"birthDate"`
	Gender     string `bson:"gender" json:"gender"`
}

type PatientResponse struct {
	FhirId     *string `json:"fhirId,omitempty"`
	GivenName  *string `json:"givenName,omitempty"`
	FamilyName *string `json:"familyName,omitempty"`
	BirthDate  *string `json:"birthDate,omitempty"`
	Gender     *string `json:"gender,omitempty"`
}
