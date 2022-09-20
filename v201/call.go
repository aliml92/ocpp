package v201

type AuthorizeReq struct {
	Certificate                 string                `json:"certificate,omitempty" validate:"omitempty,max=5500"`
	IdToken                     IdTokenType           `json:"idTokenType" validate:"required"`
	Iso15118CertificateHashData []OCSPRequestDataType `json:"iso15118CertificateHashData,omitempty" validate:"omitempty,max=4,dive,required"`
}

type BootNotificationReq struct {
	Reason          string              `json:"reason" validate:"required,BootReasonEnumType"`
	ChargingStation ChargingStationType `json:"chargingStation" validate:"required"`
}

type CancelReservationReq struct {
	ReservationId *int `json:"reservationId" validate:"required"`
}

type CertificateSignedReq struct {
	CertificateChain string `json:"certificateChain" validate:"required,max=10000"`
	CertificateType  string `json:"certificateType,omitempty" validate:"omitempty,CertificateSigningUseEnumType"`
}

type ChangeAvailabilityReq struct {
	OperationalStatus string   `json:"operationalStatus" validate:"required,OperationalStatusEnumType"`
	Evse              EVSEType `json:"evse,omitempty"`
}

type ClearCacheReq struct{}

type ClearChargingProfileReq struct {
	ChargingProfileId       *int                     `json:"chargingProfileId,omitempty"`
	ChargingProfileCriteria ClearChargingProfileType `json:"chargingProfileCriteria,omitempty"`
}

type ClearDisplayMessageReq struct {
	Id *int `json:"id" validate:"required"`
}

type ClearedChargingLimitReq struct {
	ChargingLimitSource string `json:"chargingLimitSource" validate:"required,ChargingLimitSourceEnumType"`
	EvseId              *int   `json:"evseId,omitempty"`
}

type ClearVariableMonitoringReq struct {
	Id []int `json:"id" validate:"required,dive,required"`
}

type CostUpdatedReq struct {
	TotalCost     float32 `json:"totalCost" validate:"required"`
	TransactionId string  `json:"transactionId" validate:"required,max=36"`
}

type CustomerInformationReq struct {
	RequestId           *int                    `json:"requestId" validate:"required"`
	Report              bool                    `json:"report" validate:"required"`
	Clear               bool                    `json:"clear" validate:"required"`
	CustomerIdentifier  string                  `json:"customerIdentifier,omitempty" validate:"omitempty,max=64"`
	IdToken             IdTokenType             `json:"idTokenType,omitempty"`
	CustomerCertificate CertificateHashDataType `json:"customerCertificate,omitempty" `
}

type DataTransferReq struct {
	MessageId string      `json:"messageId,omitempty" validate:"omitempty,max=50"`
	Data      interface{} `json:"data,omitempty"`
	VendorId  string      `json:"vendorId" validate:"required,max=255"`
}

type DeleteCertificateReq struct {
	CertificateHashData CertificateHashDataType `json:"certificateHashData" validate:"required"`
}

type FirmwareStatusNotificationReq struct {
	Status    string `json:"status" validate:"required,FirmwareStatusEnumType"`
	RequestId *int   `json:"requestId,omitempty" `
}

type Get15118EVCertificateReq struct {
	Iso15118SchemaVersion string `json:"iso15118SchemaVersion" validate:"required,max=50"`
	Action                string `json:"action" validate:"required,CertificateActionEnumType"`
	ExiRequest            string `json:"exiRequest" validate:"required,max=5600"`
}

type GetBaseReportReq struct {
	RequestId  int    `json:"requestId"`
	ReportBase string `json:"reportBase" validate:"required,ReportBaseEnumType"`
}

type GetCertificateStatusReq struct {
	OcspRequestData OCSPRequestDataType `json:"ocspRequestData" validate:"required"`
}

type GetChargingProfilesReq struct {
	RequestId       int                          `json:"requestId" validate:"required"`
	EvseId          int                          `json:"evseId,omitempty"`
	ChargingProfile ChargingProfileCriterionType `json:"chargingProfile" validate:"required"`
}

type GetCompositeScheduleReq struct {
	Duration         *int   `json:"duration" validate:"required"`
	ChargingRateUnit string `json:"chargingRateUnit,omitempty" validate:"omitempty,ChargingRateUnitEnumType"`
	EvseId           *int   `json:"evseId" validate:"required"`
}

type GetDisplayMessagesReq struct {
	Id        []int  `json:"id,omitempty" validate:"omitempty,dive,required"`
	RequestId *int   `json:"requestId" validate:"required"`
	Priority  string `json:"priority,omitempty" validate:"omitempty,MessagePriorityEnumType"`
	State     string `json:"state,omitempty" validate:"omitempty,MessageStateEnumType"`
}

type GetInstalledCertificateIdsReq struct {
	CertificateType []string `json:"certificateType,omitempty" validate:"omitempty,dive,GetCertificateIdUseEnumType"`
}

type GetLocalListVersionReq struct{}

type GetLogReq struct {
	LogType       string            `json:"logType" validate:"required,LogEnumType"`
	RequestId     *int              `json:"requestId" validate:"required"`
	Retries       *int              `json:"retries,omitempty"`
	RetryInterval *int              `json:"retryInterval,omitempty"`
	Log           LogParametersType `json:"log" validate:"required"`
}

type GetMonitoringReportReq struct {
	RequestId          *int                    `json:"requestId" validate:"required"`
	MonitoringCriteria []string                `json:"monitoringCriteria,omitempty" validate:"omitempty,max=3,dive,required,MonitoringCriteriaEnumType"`
	ComponentVariable  []ComponentVariableType `json:"componentVariable,omitempty" validate:"omitempty,dive,required"`
}

type GetReportReq struct {
	RequestId         *int                    `json:"requestId" validate:"required"`
	ComponentCriteria []string                `json:"componentCriteria,omitempty" validate:"omitempty,max=4,dive,required,ComponentCriteriaEnumType"`
	ComponentVariable []ComponentVariableType `json:"componentVariable,omitempty" validate:"omitempty,dive,required"`
}

type GetTransactionStatusReq struct {
	TransactionId string `json:"transactionId,omitempty" validate:"omiempty,max=36"`
}

type GetVariablesReq struct {
	GetVariableData []GetVariableDataType `json:"getVariableData" validate:"required,dive,required"`
}

type HeartbeatReq struct{}

type InstallCertificateReq struct {
	CertificateType string `json:"certificateType" validate:"required,InstallCertificateUseEnumType"`
	Certificate     string `json:"certificate" validate:"required,max=5500"`
}

type LogStatusNotificationReq struct {
	Status    string `json:"status" validate:"required,UploadLogStatusEnumType"`
	RequestId *int   `json:"requestId,omitempty"`
}

type MeterValuesReq struct {
	EvseId     *int             `json:"evseId" validate:"required,gte=0"`
	MeterValue []MeterValueType `json:"meterValue" validate:"required,dive,required"`
}

type NotifyChargingLimitReq struct {
	EvseId           int                    `json:"evseId,omitempty" validate:"omitempty,gt=0"`
	ChargingLimit    ChargingLimitType      `json:"chargingLimit" validate:"required"`
	ChargingSchedule []ChargingScheduleType `json:"chargingSchedule,omitempty" validate:"omitempty,dive,required"`
}

type NotifyCustomerInformationReq struct {
	Data        string `json:"data" validate:"required,max=512"`
	Tbc         bool   `json:"tbc,omitempty"`
	SeqNo       *int   `json:"seqNo" validate:"required,gte=0"`
	GeneratedAt string `json:"generatedAt" validate:"required,ISO8601date"`
	RequestId   *int   `json:"requestId" validate:"required"`
}

type NotifyDisplayMessagesReq struct {
	RequestId   *int              `json:"requestId" validate:"required"`
	Tbc         bool              `json:"tbc,omitempty"`
	MessageInfo []MessageInfoType `json:"messageInfo,omitempty" validate:"omitempty,dive,required"`
}

type NotifyEVChargingNeedsReq struct {
	MaxScheduleTuples *int              `json:"maxScheduleTuples,omitempty"`
	EvseId            *int              `json:"evseId" validate:"required,gt=0"`
	ChargingNeeds     ChargingNeedsType `json:"chargingNeeds" validate:"required"`
}

type NotifyEVChargingScheduleReq struct {
	TimeBase         string               `json:"timeBase" validate:"required,ISO8601date"`
	EvseId           int                  `json:"evseId" validate:"required,gt=0"`
	ChargingSchedule ChargingScheduleType `json:"chargingSchedule" validate:"required"`
}

type NotifyEventReq struct {
	GeneratedAt string          `json:"generatedAt" validate:"required,ISO8601date"`
	Tbc         bool            `json:"tbc,omitempty"`
	SeqNo       *int            `json:"seqNo" validate:"required,gte=0"`
	EventData   []EventDataType `json:"eventData" validate:"required,dive,required"`
}

type NotifyMonitoringReportReq struct {
	RequestId   *int                 `json:"requestId" validate:"required"`
	Tbc         bool                 `json:"tbc,omitempty"`
	SeqNo       *int                 `json:"seqNo" validate:"required,gte=0"`
	GeneratedAt string               `json:"generatedAt" validate:"required,ISO8601date"`
	Monitor     []MonitoringDataType `json:"monitor,omitempty" validate:"omitempty,dive,required"`
}

type NotifyReportReq struct {
	RequestId   *int             `json:"requestId" validate:"required"`
	GeneratedAt string           `json:"generatedAt" validate:"required,ISO8601date"`
	Tbc         bool             `json:"tbc,omitempty"`
	SeqNo       *int             `json:"seqNo" validate:"required,gte=0"`
	ReportData  []ReportDataType `json:"monitor,omitempty" validate:"omitempty,dive,required"`
}

type PublishFirmwareReq struct {
	Location      string `json:"location" validate:"required,max=512"`
	Retries       *int   `json:"retries,omitempty"`
	CheckSum      string `json:"checkSum" validate:"required,max=32"`
	RequestId     *int   `json:"requestId" validate:"required"`
	RetryInterval *int   `json:"retryInterval,omitempty"`
}

type PublishFirmwareStatusNotificationReq struct {
	Status    string   `json:"status" validate:"required,PublishFirmwareStatusEnumType"`
	Location  []string `json:"location,omitempty" validate:"omitempty,dive,required,max=512"`
	RequestId *int     `json:"requestId,omitempty"`
}

type ReportChargingProfilesReq struct {
	RequestId           *int                  `json:"requestId" validate:"required"`
	ChargingLimitSource string                `json:"chargingLimitSource" validate:"required,ChargingLimitSourceEnumType"`
	Tbc                 bool                  `json:"tbc,omitempty"`
	EvseId              *int                  `json:"evseId" validate:"required,gte=0"`
	ChargingProfile     []ChargingProfileType `json:"chargingProfile" validate:"required,dive,required"`
}

type RequestStartTransactionReq struct {
	EvseId          *int                `json:"evseId,omitempty" validate:"omitempty,gt=0"`
	RemoteStartId   *int                `json:"remoteStartId" validate:"required"`
	IdToken         IdTokenType         `json:"idToken" validate:"required"`
	ChargingProfile ChargingProfileType `json:"chargingProfile,omitempty"`
	GroupIdToken    IdTokenType         `json:"groupIdToken,omitempty"`
}

type RequestStopTransactionReq struct {
	TransactionId string `json:"transactionId" validate:"required,max=36"`
}

type ReservationStatusUpdateReq struct {
	ReservationId           *int   `json:"reservationId" validate:"required"`
	ReservationUpdateStatus string `json:"reservationUpdateStatus" validate:"required,ReservationUpdateStatusEnumType"`
}

type ReserveNowReq struct {
	Id             *int        `json:"id" validate:"required"`
	ExpiryDateTime string      `json:"expiryDateTime" validate:"required,ISO8601date"`
	ConnectorType  string      `json:"connectorType,omitempty" validate:"omitempty,ConnectorEnumType"`
	EvseId         *int        `json:"evseId,omitempty"`
	IdToken        IdTokenType `json:"idToken" validate:"required"`
	GroupIdToken   IdTokenType `json:"groupIdToken,omitempty"`
}

type ResetReq struct {
	Type   string `json:"type" validate:"required,ResetTypeEnumType"`
	EvseId *int   `json:"evseId,omitempty"`
}

type SecurityEventNotificationReq struct {
	Type      string `json:"type" validate:"required,ResetTypeEnumType"`
	Timestamp string `json:"timestamp" validate:"required,ISO8601date"`
	TechInfo  string `json:"techInfo,omitempty" validate:"omitempty,max=255"`
}

type SendLocalListReq struct {
	VersionNumber          *int                `json:"versionNumber" validate:"required"`
	UpdateType             string              `json:"updateType" validate:"required,UpdateTypeEnumType"`
	LocalAuthorizationList []AuthorizarionData `json:"localAuthorizationList,omitempty" validate:"omitempty,dive,required"`
}

type SetChargingProfileReq struct {
	EvseId          *int                `json:"evseId" validate:"required,gte=0"`
	ChargingProfile ChargingProfileType `json:"chargingProfile" validate:"required"`
}

type SetDisplayMessageReq struct {
	Message MessageInfoType `json:"message" validate:"required"`
}

type SetMonitoringBaseReq struct {
	MonitoringBase string `json:"monitoringBase" validate:"required,MonitoringBaseEnumType"`
}

type SetMonitoringLevelReq struct {
	Severity *int `json:"severity" validate:"required,gte=0,lte=9"`
}

type SetNetworkProfileReq struct {
	ConfigurationSlot *int                         `json:"configurationSlot" validate:"required"`
	ConnectionData    NetworkConnectionProfileType `json:"connectionData" validate:"required"`
}

type SetVariableMonitoringReq struct {
	SetMonitoringData []SetMonitoringDataType `json:"setMonitoringData" validate:"required,dive,required"`
}

type SetVariablesReq struct {
	SetVariableData []SetVariableDataType `json:"setVariableData" validate:"required,dive,required"`
}

type SignCertificateReq struct {
	Csr         string `json:"csr" validate:"required,max=5500"`
	Certificate string `json:"certificate,omitempty" validate:"omitempty,CertificateSigningUseEnumType"`
}

type StatusNotificationReq struct {
	Timestamp       string `json:"timestamp" validate:"required,ISO8601date"`
	ConnectorStatus string `json:"connectorStatus" validate:"required,ConnectorStatusEnumType"`
	EvseId          *int   `json:"evseId" validate:"required"`
	ConnectorId     *int   `json:"connectorId" validate:"required"`
}

type TransactionEventReq struct {
	EventType          string           `json:"eventType" validate:"required,TransactionEventTypeEnumType"`
	Timestamp          string           `json:"timestamp" validate:"required,ISO8601date"`
	TriggerReason      string           `json:"triggerReason" validate:"required,TriggerReasonEnumType"`
	SeqNo              *int             `json:"seqNo" validate:"required"`
	Offline            bool             `json:"offline,omitempty"`
	NumberOfPhasesUsed *int             `json:"numberOfPhasesUsed,omitempty"`
	CableMaxCurrent    *int             `json:"cableMaxCurrent,omitempty"`
	ReservationId      *int             `json:"reservationId,omitempty"`
	TransactionInfo    TransactionType  `json:"transactionInfo" validate:"required"`
	IdToken            IdTokenType      `json:"idToken,omitempty"`
	Evse               EVSEType         `json:"evse,omitempty"`
	MeterValue         []MeterValueType `json:"meterValue,omitempty" validate:"omitempty,dive,required"`
}

type TriggerMessageReq struct {
	RequestedMessage string   `json:"requestedMessage" validate:"required,MessageTriggerEnumType"`
	Evse             EVSEType `json:"evse,omitempty"`
}

type UnlockConnectorReq struct {
	EvseId      *int `json:"evseId" validate:"required,gte=0"`
	ConnectorId *int `json:"connectorId" validate:"required,gte=0"`
}

type UnpublishFirmwareReq struct {
	Checksum string `json:"checksum" validate:"required,max=32"`
}

type UpdateFirmwareReq struct {
	Retries       *int         `json:"retries,omitempty"`
	RetryInterval *int         `json:"retryInterval,omitempty"`
	RequestId     *int         `json:"requestId" validate:"required"`
	Firmware      FirmwareType `json:"firmware" validate:"required"`
}
