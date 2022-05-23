package v16



type IdTagInfo struct {
	ExpiryDate   string              `json:"expiryDate,omitempty" validate:"omitempty,ISO8601date"`
	ParentIdTag  string              `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	Status       string              `json:"status" validate:"required,authorizationStatus"`
}

type RegistrationStatus string













