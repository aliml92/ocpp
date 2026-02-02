package v16

type IdTagInfo struct {
	ExpiryDate  string `json:"expiryDate,omitempty" validate:"omitempty,ISO8601date"`
	ParentIdTag string `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	Status      string `json:"status" validate:"required,AuthorizationStatus"`
}

type MeterValue struct {
	Timestamp    string         `json:"timestamp" validate:"required,ISO8601date"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,dive,required"`
}

type SampledValue struct {
	Value     string `json:"value" validate:"required"`
	Context   string `json:"context,omitempty" validate:"omitempty,ReadingContext"`
	Format    string `json:"format,omitempty" validate:"omitempty,ValueFormat"`
	Measurand string `json:"measurand,omitempty" validate:"omitempty,Measurand"`
	Phase     string `json:"phase,omitempty" validate:"omitempty,Phase"`
	Location  string `json:"location,omitempty" validate:"omitempty,Location"`
	Unit      string `json:"unit,omitempty" validate:"omitempty,UnitOfMeasure"`
}

type ChargingProfile struct {
	ChargingProfileId      int              `json:"chargingProfileId" validate:"required,gte=0"`
	TransactionId          int              `json:"transactionId,omitempty"`
	StackLevel             int              `json:"stackLevel" validate:"required,gte=0"`
	ChargingProfilePurpose string           `json:"chargingProfilePurpose" validate:"required,ChargingProfilePurposeType"`
	ChargingProfileKind    string           `json:"chargingProfileKind" validate:"required,ChargingProfileKindType"`
	Context                string           `json:"context,omitempty" validate:"omitempty,ReadingContext"`
	RecurrencyKind         string           `json:"recurrencyKind,omitempty" validate:"omitempty,RecurrencyKindType"`
	ValidFrom              string           `json:"validFrom,omitempty" validate:"omitempty,ISO8601date"`
	ValidTo                string           `json:"validTo,omitempty" validate:"omitempty,ISO8601date"`
	ChargingSchedule       ChargingSchedule `json:"chargingSchedule" validate:"required,dive,required"`
}

type ChargingSchedule struct {
	Duration               int                      `json:"duration,omitempty"`
	StartSchedule          string                   `json:"startSchedule,omitempty" validate:"omitempty,ISO8601date"`
	ChargingRateUnit       string                   `json:"chargingRateUnit" validate:"required,ChargingRateUnitType"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,dive,required"`
	MinChargingRate        float32                  `json:"minChargingRate,omitempty"`
}

type ChargingSchedulePeriod struct {
	StartPeriod  *int    `json:"startPeriod" validate:"required"`
	Limit        float32 `json:"limit" validate:"required,gte=0"`
	NumberPhases int     `json:"numberPhases,omitempty"`
}

type AuthorizationData struct {
	IdTag     string    `json:"idTag" validate:"required,max=20"`
	IdTagInfo IdTagInfo `json:"idTagInfo,omitempty"`
}

// OCPP 1.6 security whitepaper edition 2 implementation

type CertificateHashDataType struct {
	HashAlgorithm  string `json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"`
	IssuerNameHash string `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string `json:"serialNumber" validate:"required,max=40"`
}

type FirmwareType struct {
	Location           string `json:"location" validate:"required,max=512"`
	RetrieveDateTime   string `json:"retrieveDateTime" validate:"required,ISO8601date"`
	InstallDateTime    string `json:"installDateTime,omitempty" validate:"omitempty,ISO8601date"`
	SigningCertificate string `json:"signingCertificate" validate:"required,max=5500"`
	Signature          string `json:"signature" validate:"required,max=800"`
}

type LogParametersType struct {
	RemoteLocation  string `json:"remoteLocation" validate:"required,max=512"`
	OldestTimestamp string `json:"oldestTimestamp,omitempty" validate:"omitempty,ISO8601date"`
	LatestTimestamp string `json:"latestTimestamp,omitempty" validate:"omitempty,ISO8601date"`
}

type KeyValue struct {
	Key      string `json:"key" validate:"required,max=50"`
	Readonly bool   `json:"readonly" validate:"required"`
	Value    string `json:"value,omitempty" validate:"omitempty,max=500"`
}
