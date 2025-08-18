package models

type Practitioner struct {
	FhirId     string `bson:"fhirId" json:"fhirId"`
	GivenName  string `bson:"givenName" json:"givenName"`
	FamilyName string `bson:"familyName" json:"familyName"`
}

type PractitionerRespose struct {
	FhirId     *string `json:"fhirId,omitempty"`
	GivenName  *string `json:"givenName,omitempty"`
	FamilyName *string `json:"familyName,omitempty"`
}
