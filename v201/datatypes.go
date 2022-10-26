package v201

type ACChargingParametersType struct {
	EnergyAmount int `json:"energyAmount" validate:"required"`
	EvMinCurrent int `json:"evMinCurrent" validate:"required"`
	EvMaxCurrent int `json:"evMaxCurrent" validate:"required"`
	EvMaxVoltage int `json:"evMaxVoltage" validate:"required"`
}

type AdditionalInfoType struct {
	AdditionalToken string `json:"additionalToken" validate:"required,max=36"`
	Type            string `json:"type" validate:"required,max=50"`
}

type APNType struct {
	Apn                     string `json:"apn" validate:"required,max=100"`
	ApnUserName             string `json:"apnUserName,omitempty" validate:"omitempty,max=20"`
	ApnPassword             string `json:"apnPassword,omitempty" validate:"omitempty,max=20"`
	SimPin                  int    `json:"simPin,omitempty"`
	PreferredNetwork        string `json:"preferredNetwork,omitempty" validate:"omitempty,max=6"`
	UseOnlyPreferredNetwork bool   `json:"useOnlyPreferredNetwork,omitempty"`
	ApnAuthentication       string `json:"apnAuthentication" validate:"required,APNAuthenticationEnumType"`
}

type AuthorizarionData struct {
	IdTokenInfo IdTokenInfoType `json:"idTokenInfo,omitempty"`
	IdToken     IdTokenType     `json:"idToken" validate:"required"`
}

type CertificateHashDataChainType struct {
	CertificateType          string                    `json:"certificateType" validate:"required,GetCertificateIdUseEnumType"`
	CertificateHashData      CertificateHashDataType   `json:"certificateHashData" validate:"required"`
	ChildCertificateHashData []CertificateHashDataType `json:"childCertificateHashData,omitempty" validate:"omitempty,max=4,dive,required"`
}

type CertificateHashDataType struct {
	HashAlgorithm  string `json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"` // todo: validation register required
	IssuerNameHash string `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string `json:"serialNumber" validate:"required,max=40"`
}

type ChargingLimitType struct {
	ChargingLimitSource string `json:"chargingLimitSource" validate:"required,ChargingLimitSourceEnumType"`
	IsGridCritical      bool   `json:"isGridCritical,omitempty"`
}

type ChargingNeedsType struct {
	RequestedEnergyTransfer string                   `json:"requestedEnergyTransfer" validate:"required,RequestedEnergyTransferEnumType"`
	DepartureTime           string                   `json:"departureTime" validate:"omitempty,ISO8601date"`
	ACChargingParameters    ACChargingParametersType `json:"acChargingParameters,omitempty"`
	DcChargingParameters    DCChargingParametersType `json:"dcChargingParameters,omitempty"`
}

type ChargingProfileCriterionType struct {
	ChargingProfilePurpose string   `json:"chargingProfilePurpose,omitempty" validate:"omitempty,ChargingProfilePurposeEnumType"`
	StackLevel             *int     `json:"stackLevel,omitempty"`
	ChargingProfileId      []int    `json:"chargingProfileId,omitempty" validate:"omitempty,dive,required"`
	ChargingLimitSource    []string `json:"chargingLimitSource,omitempty" validate:"omitempty,max=4,dive,required,ChargingLimitSourceEnumType"`
}

type ChargingProfileType struct {
	Id                     int                 `json:"id" validate:"required"`
	StackLevel             *int                 `json:"stackLevel" validate:"required"`
	ChargingProfilePurpose string               `json:"chargingProfilePurpose" validate:"required,ChargingProfilePurposeEnumType"`
	ChargingProfileKind    string               `json:"chargingProfileKind" validate:"required,ChargingProfileKindEnumType"`
	RecurrencyKind         string               `json:"recurrencyKind,omitempty" validate:"omitempty,RecurrencyKindEnumType"`
	ValidFrom              string               `json:"validFrom,omitempty" validate:"omitempty,ISO8601date"`
	ValidTo                string               `json:"validTo,omitempty" validate:"omitempty,ISO8601date"`
	TransactionId          string               `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	ChargingSchedule       ChargingScheduleType `json:"chargingSchedule" validate:"required,min=1,max=3,dive,required"`
}

type ChargingSchedulePeriodType struct {
	StartPeriod  *int    `json:"startPeriod" validate:"required"`
	Limit        float32 `json:"limit" validate:"required"`
	NumberPhases int     `json:"numberPhases,omitempty" validate:"omitempty,gte=1,lte=3"`
	PhaseToUse   int     `json:"phaseToUse,omitempty" validate:"omitempty,gte=1,lte=3"`
}

type ChargingScheduleType struct {
	Id                     int                         `json:"id" validate:"required"`
	StartSchedule          string                       `json:"startSchedule,omitempty" validate:"omitempty,ISO8601date"`
	Duration               int                         `json:"duration,omitempty"`
	ChargingRateUnit       string                       `json:"chargingRateUnit" validate:"required,ChargingRateUnitEnumType"`
	MinChargingRate        float32                      `json:"minChargingRate,omitempty"`
	ChargingSchedulePeriod []ChargingSchedulePeriodType `json:"chargingSchedulePeriod" validate:"required,min=1,max=1024,dive,required"`
	SalesTariff            SalesTariffType              `json:"salesTariff,omitempty"`
}

type ChargingStationType struct {
	SerialNumber    string    `json:"serialNumber,omitempty" validate:"omitempy,max=25"`
	Model           string    `json:"model" validate:"required,max=20"`
	VendorName      string    `json:"vendorName" validate:"required,max=50"`
	FirmwareVersion string    `json:"firmwareVersion,omitempty" validate:"omitempty,max=50"`
	Modem           ModemType `json:"modem,omitempty"`
}

type ClearChargingProfileType struct {
	EvseId                 *int   `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfilePurpose string `json:"chargingProfilePurpose,omitempty" validate:"omitempty,ChargingProfilePurposeEnumType"` // todo: validation register required
	StackLevel             *int   `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
}

type ClearMonitoringResultType struct {
	Status     string         `json:"status" validate:"required,ClearMonitoringStatusEnumType"` // todo: validation register required
	Id         *int           `json:"id" validate:"required"`
	StatusInfo StatusInfoType `json:"statusInfo,omitempty" `
}

type ComponentType struct {
	Name     string   `json:"name" validate:"required,max=50"`
	Instance string   `json:"instance,omitempty" validate:"omitempty,max=50"`
	Evse     EVSEType `json:"evse,omitempty"`
}

type ComponentVariableType struct {
	Component ComponentType `json:"component" validate:"required"`
	Variable  VariableType  `json:"variable,omitempty"`
}

type CompositeScheduleType struct {
	EvseId                 *int                         `json:"evseId" validate:"required,gte=0"`
	Duration               int                         `json:"duration" validate:"required"`
	ScheduleStart          string                       `json:"scheduleStart" validate:"required,ISO8601date"`
	ChargingRateUnit       string                       `json:"chargingRateUnit" validate:"required,ChargingRateUnitEnumType"`
	ChargingSchedulePeriod []ChargingSchedulePeriodType `json:"chargingSchedulePeriod" validate:"required,min=1,dive,required"`
}

type ConsumptionCostType struct {
	StartValue float32  `json:"startValue" validate:"required"`
	Cost       CostType `json:"cost" validate:"required,min=1,max=3,dive,required"`
}

type CostType struct {
	CostKind         string `json:"costKind" validate:"required,CostKindEnumType"`
	Amount           int    `json:"amount" validate:"required"`
	AmountMultiplier *int   `json:"amountMultiplier,omitempty" validate:"omitempty,gte=-3,lte=3"`
}

type DCChargingParametersType struct {
	EvMaxCurrent     int  `json:"evMaxCurrent" validate:"required"`
	EvMaxVoltage     int  `json:"evMaxVoltage" validate:"required"`
	EnergyAmount     int  `json:"energyAmount,omitempty"`
	EvMaxPower       int  `json:"evMaxPower,omitempty"`
	StateOfCharge    *int `json:"stateOfCharge,omitempty" validate:"omitempty,gte=0,lte=100"`
	EvEnergyCapacity int  `json:"evEnergyCapacity,omitempty"`
	FullSoC          *int `json:"fullSoC,omitempty" validate:"omitempty,gte=0,lte=100"`
	BulkSoC          *int `json:"bulkSoC,omitempty" validate:"omitempty,gte=0,lte=100"`
}

type EventDataType struct {
	EventId               int          `json:"eventId" validate:"required"`
	Timestamp             string        `json:"timestamp" validate:"required,ISO8601date"`
	Trigger               string        `json:"trigger" validate:"required,TriggerEnumType"`
	Cause                 int          `json:"cause,omitempty"`
	ActualValue           string        `json:"actualValue" validate:"required,max=2500"`
	TechCode              string        `json:"techCode,omitempty" validate:"omitempty,max=50"`
	TechInfo              string        `json:"techInfo,omitempty" validate:"omitempty,max=500"`
	Cleared               bool          `json:"cleared,omitempty"`
	TransactionId         string        `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	VariableMonitoringId  int          `json:"variableMonitoringId,omitempty"`
	EventNotificationType string        `json:"eventNotificationType" validate:"required,EventNotificationTypeEnumType"`
	Component             ComponentType `json:"component" validate:"required"`
	Variable              VariableType  `json:"variable" validate:"required"`
}

type EVSEType struct {
	Id          int  `json:"id" validate:"required,gt=0"`
	ConnectorId *int `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
}

type FirmwareType struct {
	Location           string `json:"location" validate:"required,max=512"`
	RetrieveDateTime   string `json:"retrieveDateTime" validate:"required,ISO8601date"`
	InstallDateTime    string `json:"installDateTime,omitempty" validate:"omitempty,ISO8601date"`
	SigningCertificate string `json:"signingCertificate,omitempty" validate:"omitempty,max=5500"`
	Signature          string `json:"signature,omitempty" validate:"omitempty,max=800"`
}

type GetVariableDataType struct {
	AttributeType string        `json:"attributeType,omitempty" validate:"omitempty,AttributeTypeEnumType"`
	Component     ComponentType `json:"component" validate:"required"`
	Variable      VariableType  `json:"variable" validate:"required"`
}

type GetVariableResultType struct {
	AttributeStatus     string         `json:"attributeStatus" validate:"required,GetVariableStatusEnumType"`
	AttributeType       string         `json:"attributeType,omitempty" validate:"omitempty,AttributeTypeEnumType"`
	AttributeValue      string         `json:"attributeValue,omitempty" validate:"omitempty,max=2500"`
	Component           ComponentType  `json:"component" validate:"required"`
	Variable            VariableType   `json:"variable" validate:"required"`
	AttributeStatusInfo StatusInfoType `json:"attributeStatusInfo,omitempty"`
}

type IdTokenInfoType struct {
	Status              string             `json:"status" validate:"required,AuthorizationStatusEnumType"`
	CacheExpiryDateTime string             `json:"cacheExpiryDateTime,omitempty" validate:"omitempty,ISO8601date"`
	ChangePriority      *int               `json:"changePriority,omitempty" validate:"omitempty,gte=-9,lte=9"`
	Language1           string             `json:"language1,omitempty" validate:"omitempty,max=8"`
	EvseId              []int              `json:"evseId,omitempty" validate:"omitempty,dive,required"`
	Language2           string             `json:"language2,omitempty" validate:"omitempty,max=8"`
	GroupIdToken        IdTokenType        `json:"groupIdToken,omitempty"`
	PersonalMessage     MessageContentType `json:"personalMessage,omitempty"`
}

type IdTokenType struct {
	IdToken        string               `json:"idToken" validate:"required,max=36"`
	Type           string               `json:"type" validate:"required,IdTokenEnumType"` // todo: validation register required
	AdditionalInfo []AdditionalInfoType `json:"additionalInfo" validate:"omitempty,dive,required"`
}

type LogParametersType struct {
	RemoteLocation  string `json:"remoteLocation" validate:"required,max=512"`
	OldestTimestamp string `json:"oldestTimestamp,omitempty" validate:"omitempty,ISO8601date"`
	LatestTimeStamp string `json:"latestTimeStamp,omitempty" validate:"omitempty,ISO8601date"`
}

type MessageContentType struct {
	Format   string `json:"format" validate:"required,MessageFormatEnumType"`
	Language string `json:"language,omitempty" validate:"omitempty,max=8"`
	Content  string `json:"content" validate:"required,max=512"`
}

type MessageInfoType struct {
	Id            int               `json:"id" validate:"required"`
	Priority      string             `json:"priority" validate:"required,MessagePriorityEnumType"`
	State         string             `json:"state,omitempty" validate:"omitempty,MessageStateEnumType"`
	StartDateTime string             `json:"startDateTime,omitempty" validate:"omitempty,ISO8601date"`
	EndDateTime   string             `json:"endDateTime,omitempty" validate:"omitempty,ISO8601date"`
	TransactionId string             `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	Message       MessageContentType `json:"message" validate:"required"`
	Display       ComponentType      `json:"display,omitempty"`
}

type MeterValueType struct {
	Timestamp    string             `json:"timestamp" validate:"required,ISO8601date"`
	SampledValue []SampledValueType `json:"sampledValue" validate:"required,dive,required"`
}

type ModemType struct {
	Iccid string `json:"iccid,omitempty" validate:"omitempty,max=20"`
	Imsi  string `json:"imsi,omitempty" validate:"omitempty,max=20"`
}

type MonitoringDataType struct {
	Component          ComponentType            `json:"component" validate:"required"`
	Variable           VariableType             `json:"variable" validate:"required"`
	VariableMonitoring []VariableMonitoringType `json:"variableMonitoring" validate:"required,dive,required"`
}

type NetworkConnectionProfileType struct {
	OcppVersion     string  `json:"ocppVersion" validate:"required,OCPPVersionEnumType"`
	OcppTransport   string  `json:"ocppTransport" validate:"required,OCPPTransportEnumType"`
	OcppCsmsUrl     string  `json:"ocppCsmsUrl" validate:"required,max=512"`
	MessageTimeout  int     `json:"messageTimeout" validate:"required"`
	SecurityProfile int    `json:"securityProfile" validate:"required"`
	OcppInterface   string  `json:"ocppInterface" validate:"required,OCPPInterfaceEnumType"`
	Vpn             VPNType `json:"vpn,omitempty"`
	Apn             APNType `json:"apn,omitempty"`
}

type OCSPRequestDataType struct {
	HashAlgorithm  string `json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"` // todo: validation register required
	IssuerNameHash string `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string `json:"serialNumber" validate:"required,max=40"`
	ResponderUrl   string `json:"responderUrl" validate:"required,max=512"`
}

type RelativeTimeIntervalType struct {
	Start int `json:"start" validate:"required"`
	End   int `json:"end,omitempty"`
}

type ReportDataType struct {
	Component               ComponentType                 `json:"component" validate:"required"`
	Variable                VariableType                  `json:"variable" validate:"required"`
	VariableAttribute       []VariableAttributeType       `json:"variableAttribute" validate:"required,min=1,max=4,dive,required"`
	VariableCharacteristics []VariableCharacteristicsType `json:"variableCharacteristics,omitempty"`
}

type SalesTariffEntryType struct {
	EPriceLevel          int                     `json:"ePriceLevel,omitempty"`
	RelativeTimeInterval RelativeTimeIntervalType `json:"relativeTimeInterval" validate:"required"`
	ConsumptionCost      ConsumptionCostType      `json:"consumptionCost,omitempty" validate:"omitempty,max=3,dive,required"`
}

type SalesTariffType struct {
	Id                     int                   `json:"id" validate:"required"`
	SalesTariffDescription string                 `json:"salesTariffDescription,omitempty" validate:"omitempty,max=512"`
	NumEPriceLevels        int                   `json:"numEPriceLevels,omitempty"`
	SalesTariffEntry       []SalesTariffEntryType `json:"salesTariffEntry" validate:"required,gte=1,lte=1024,dive,required"`
}

type SampledValueType struct {
	Value            float32              `json:"value" validate:"required"`
	Context          string               `json:"context,omitempty" validate:"omitempty,ReadingContextEnumType"`
	Measurand        string               `json:"measurand,omitempty" validate:"omitempty,MeasurandEnumType"`
	Phase            string               `json:"phase,omitempty" validate:"omitempty,PhaseEnumType"`
	Location         string               `json:"location,omitempty" validate:"omitempty,LocationEnumType"`
	SignedMeterValue SignedMeterValueType `json:"signedMeterValue,omitempty"`
	UnitOfMeasure    UnitOfMeasureType    `json:"unitOfMeasure,omitempty"`
}

type SetMonitoringDataType struct {
	Id          int          `json:"id,omitempty"`
	Transaction bool          `json:"transaction,omitempty"`
	Value       float32       `json:"value" validate:"required"`
	Type        string        `json:"type" validate:"required,MonitoringEnumType"`
	Severity    *int          `json:"severity" validate:"required,gte=0,lte=9"`
	Component   ComponentType `json:"component" validate:"required"`
	Variable    VariableType  `json:"variable" validate:"required"`
}

type SetMonitoringResultType struct {
	Id         int           `json:"id,omitempty"`
	Status     string         `json:"status" validate:"required,SetMonitoringStatusEnumType"`
	Type       string         `json:"type" validate:"required,MonitoringEnumType"`
	Severity   *int           `json:"severity" validate:"required,gte=0,lte=9"`
	Component  ComponentType  `json:"component" validate:"required"`
	Variable   VariableType   `json:"variable" validate:"required"`
	StatusInfo StatusInfoType `json:"statusInfo,omitempty"`
}

type SetVariableDataType struct {
	AttributeType  string        `json:"attributeType,omitempty" validate:"omitempty,AttributeTypeEnumType"`
	AttributeValue string        `json:"attributeValue" validate:"required,max=10000"`
	Component      ComponentType `json:"component" validate:"required"`
	Variable       VariableType  `json:"variable" validate:"required"`
}

type SetVariableResultType struct {
	AttributeType       string         `json:"attributeType,omitempty" validate:"omitempty,AttributeTypeEnumType"`
	AttributeStatus     string         `json:"attributeStatus" validate:"required,SetVariableStatusEnumType"`
	Component           ComponentType  `json:"component" validate:"required"`
	Variable            VariableType   `json:"variable" validate:"required"`
	AttributeStatusInfo StatusInfoType `json:"attributeStatusInfo,omitempty"`
}

type SignedMeterValueType struct {
	SignedMeterData string `json:"signedMeterData" validate:"required,max=2500"`
	SigningMethod   string `json:"signingMethod" validate:"required,max=50"`
	EncodingMethod  string `json:"encodingMethod" validate:"required,max=50"`
	PublicKey       string `json:"publicKey" validate:"required,max=2500"`
}

type StatusInfoType struct {
	ReasonCode     string `json:"reasonCode" validate:"required,max=20"`
	AdditionalInfo string `json:"additionalInfo" validate:"required,max=512"`
}

type TransactionType struct {
	TransactionId     string `json:"transactionId" validate:"required,max=36"`
	ChargingState     string `json:"chargingState,omitempty" validate:"omitempty,ChargingStateEnumType"`
	TimeSpentCharging int    `json:"timeSpentCharging,omitempty"`
	StoppedReason     string `json:"stoppedReason,omitempty" validate:"omitempty,ReasonEnumType"`
	RemoteStartId     int   `json:"remoteStartId,omitempty"`
}

type UnitOfMeasureType struct {
	Unit       string `json:"unit,omitempty" validate:"omitempty,max=50"`
	Multiplier *int   `json:"multiplier,omitempty"`
}

type VariableAttributeType struct {
	Type       string `json:"type,omitempty" validate:"omitempty,AttributeTypeEnumType"`
	Value      string `json:"value,omitempty" validate:"omitempty,max=2500"`
	Mutability string `json:"mutability,omitempty" validate:"omitempty,MutabilityEnumType"`
	Persistent bool   `json:"persistent,omitempty"`
	Constant   bool   `json:"constant,omitempty"`
}

type VariableCharacteristicsType struct {
	Unit               string  `json:"unit,omitempty" validate:"omitempty,max=16"`
	DataType           string  `json:"dataType" validate:"required,DataEnumType"`
	MinLimit           float32 `json:"minLimit,omitempty"`
	MaxLimit           float32 `json:"maxLimit,omitempty"`
	ValuesList         string  `json:"valuesList,omitempty" validate:"omitempty,max=1000"`
	SupportsMonitoring bool    `json:"supportsMonitoring" validate:"required"`
}

type VariableMonitoringType struct {
	Id          int    `json:"id" validate:"required"`
	Transaction bool    `json:"transaction" validate:"required"`
	Value       float32 `json:"value" validate:"required"`
	Type        string  `json:"type" validate:"required,MonitorEnumType"`
	Severity    *int    `json:"severity" validate:"required,gte=0,lte=9"`
}

type VariableType struct {
	Name     string `json:"name" validate:"required,max=50"`
	Instance string `json:"instance,omitempty" validate:"omitempty,max=50"`
}

type VPNType struct {
	Server   string `json:"server" validate:"required,max=512"`
	User     string `json:"user" validate:"required,max=20"`
	Group    string `json:"group,omitempty" validate:"omitempty,max=20"`
	Password string `json:"password" validate:"required,max=20"`
	Key      string `json:"key" validate:"required,max=255"`
	Type     string `json:"type" validate:"required,VPNTypeEnumType"`
}
