package v16

type AuthorizeReq struct {
	IdTag string `json:"idTag" validate:"required,max=20"`
}

type BootNotificationReq struct {
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" validate:"max=25"`
	ChargePointModel        string `json:"chargePointModel" validate:"required,max=20"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty" validate:"max=25"`
	ChargePointVendor       string `json:"chargePointVendor" validate:"required,max=20"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty" validate:"max=50"`
	Iccid                   string `json:"iccid,omitempty" validate:"max=20"`
	Imsi                    string `json:"imsi,omitempty" validate:"max=20"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty" validate:"max=25"`
	MeterType               string `json:"meterType,omitempty" validate:"max=25"`
}

type DataTransferReq struct {
	VendorId  string `json:"vendorId" validate:"required,max=255"`
	MessageId string `json:"messageId,omitempty" validate:"max=255"`
	Data      string `json:"data,omitempty"`
}

type DiagnosticsStatusNotificationReq struct {
	Status string `json:"status" validate:"required,DiagnosticsStatus"`
}

type FirmwareStatusNotificationReq struct {
	Status string `json:"status" validate:"required,FirmwareStatus"`
}

type HeartbeatReq struct {
}

type MeterValuesReq struct {
	ConnectorId   *int        `json:"connectorId" validate:"required,gte=0"`
	TransactionId int        `json:"transactionId"`
	MeterValue    MeterValue `json:"meterValue" validate:"required,dive,required"`
}

type StartTransactionReq struct {
	ConnectorId   int    `json:"connectorId" validate:"required,gt=0"`
	IdTag         string `json:"idTag" validate:"required,max=20"`
	MeterStart    *int    `json:"meterStart" validate:"required,gte=0"`
	ReservationId int    `json:"reservationId,omitempty" validate:"omitempty,gt=0"`
	Timestamp     string `json:"timestamp" validate:"required,ISO8601date"`
}

type StatusNotificationReq struct {
	ConnectorId     *int    `json:"connectorId" validate:"required,gte=0"`
	ErrorCode       string `json:"errorCode" validate:"required,ChargePointErrorCode"`
	Info            string `json:"info,omitempty" validate:"max=50"`
	Status          string `json:"status" validate:"required,ChargePointStatus"`
	Timestamp       string `json:"timestamp,omitempty" validate:"ISO8601date"`
	VendorId        string `json:"vendorId,omitempty" validate:"max=255"`
	VendorErrorCode string `json:"vendorErrorCode,omitempty" validate:"max=50"`
}

type StopTransactionReq struct {
	IdTag           string     `json:"idTag" validate:"required,max=20"`
	MeterStop       int        `json:"meterStop" validate:"required"`
	Timestamp       string     `json:"timestamp" validate:"required,ISO8601date"`
	TransactionId   int        `json:"transactionId" validate:"required"`
	Reason          string     `json:"reason,omitempty" validate:"Reason"`
	TransactionData MeterValue `json:"transactionData,omitempty" validate:"dive,required"`
}

type CancelReservationReq struct {
	ReservationId int `json:"reservationId" validate:"required"`
}

// actions by csms

type ChangeAvailabilityReq struct {
	ConnectorId *int    `json:"connectorId" validate:"required,gte=0"`
	Type        string `json:"type" validate:"required,AvailabilityType"`
}

type ChangeConfigurationReq struct {
	Key   string `json:"key" validate:"required,max=50"`
	Value string `json:"value" validate:"required,max=50"`
}

type ClearCacheReq struct{}

type ClearChargingProfileReq struct {
	Id                     string `json:"id,omitempty"`
	ConnectorId            *int    `json:"connectorId,omitempty" validate:"gte=0"`
	ChargingProfilePurpose string `json:"chargingProfilePurpose,omitempty" validate:"ChargingProfilePurpose"`
	StackLevel             *int    `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
}

type GetCompositeScheduleReq struct {
	ConnectorId      *int    `json:"connectorId,omitempty" validate:"gte=0"`
	Duration         *int    `json:"duration" validate:"required,gte=0"`
	ChargingRateUnit string `json:"chargingRateUnit,omitempty" validate:"ChargingRateUnit"`
}

type GetConfigurationReq struct {
	Key []string `json:"key,omitempty" validate:"dive,max=50"`
}

type GetDiagnosticsReq struct {
	Location      string `json:"location" validate:"required"`
	Retries       *int    `json:"retries,omitempty" validate:"gte=0"`
	RetryInterval *int    `json:"retryInterval,omitempty" validate:"gte=0"`
	StartTime     string `json:"startTime,omitempty" validate:"ISO8601date"`
	StopTime      string `json:"stopTime,omitempty" validate:"ISO8601date"`
}

type GetLocalListVersionReq struct{}

type RemoteStartTransactionReq struct {
	ConnectorId     *int             `json:"connectorId,omitempty" validate:"gte=0"`
	IdTag           string          `json:"idTag" validate:"required,max=20"`
	ChargingProfile ChargingProfile `json:"chargingProfile,omitempty" validate:"dive,required"`
}

type RemoteStopTransactionReq struct {
	TransactionId int `json:"transactionId" validate:"required"`
}

type ReserveNowReq struct {
	ConnectorId   *int    `json:"connectorId" validate:"required,gte=0"`
	ExpiryDate    string `json:"expiryDate" validate:"required,ISO8601date"`
	IdTag         string `json:"idTag" validate:"required,max=20"`
	ParentIdTag   string `json:"parentIdTag,omitempty" validate:"max=20"`
	ReservationId int    `json:"reservationId" validate:"required"`
}

type ResetReq struct {
	Type string `json:"type" validate:"required,ResetType"`
}

type SendLocalListReq struct {
	ListVersion            *int                 `json:"listVersion" validate:"required,gte=0"`
	LocalAuthorizationList []AuthorizationData `json:"localAuthorizationList,omitempty" validate:"dive,required"`
	UpdateType             string              `json:"updateType" validate:"required,UpdateType"`
}

type SetChargingProfileReq struct {
	ConnectorId     *int             `json:"connectorId" validate:"required,gte=0"`
	ChargingProfile ChargingProfile `json:"chargingProfile" validate:"required,dive,required"`
}

type TriggerMessageReq struct {
	RequestedMessage string `json:"requestedMessage" validate:"required,MessageTrigger"`
	ConnectorId      *int    `json:"connectorId,omitempty" validate:"gte=0"`
}

type UnlockConnectorReq struct {
	ConnectorId *int `json:"connectorId" validate:"required,gte=0"`
}

type UpdateFirmwareReq struct {
	Location      string `json:"location" validate:"required"`
	Retries       *int    `json:"retries,omitempty" validate:"gte=0"`
	RetrieveDate  string `json:"retrieveDate" validate:"required,ISO8601date"`
	RetryInterval *int    `json:"retryInterval,omitempty" validate:"gte=0"`
}

// OCPP 1.6 security whitepaper edition 2 implementation

type CertificateSignedReq struct {
	CertificateChain string `json:"certificateChain" validate:"required,max=10000"`
}

type DeleteCertificateReq struct {
	CertificateHashData CertificateHashDataType `json:"certificateHashData" validate:"required"`
}

type ExtendedTriggerMessageReq struct {
	RequestedMessage string `json:"requestedMessage" validate:"required,MessageTriggerEnumType"`
	ConnectorId      int    `json:"connectorId,omitempty" validate:"gt=0"`
}

type GetInstalledCertificateIdsReq struct {
	CertificateType string `json:"certificateType" validate:"required,CertificateUseEnumType"`
}

type GetLogReq struct {
	LogType       string `json:"logType" validate:"required,LogEnumType"`
	RequestId     int    `json:"requestId" validate:"required"`
	Retries       *int    `json:"retries,omitempty" validate:"gte=0"`
	RetryInterval *int    `json:"retryInterval,omitempty" validate:"gte=0"`
	Log           string `json:"log" validate:"required,LogParametersType"`
}

type InstallCertificateReq struct {
	CertificateType string `json:"certificateType" validate:"required,CertificateUseEnumType"`
	Certificate     string `json:"certificate" validate:"required,max=5500"`
}

type LogStatusNotificationReq struct {
	Status    string `json:"status" validate:"required,UploadLogStatusEnumType"`
	RequestId int    `json:"requestId,omitempty"`
}

type SecurityEventNotificationReq struct {
	Type      string `json:"type" validate:"required,max=50"`
	Timestamp string `json:"timestamp" validate:"required,ISO8601date"`
	TechInfo  string `json:"techInfo,omitempty" validate:"max=255"`
}

type SignCertificateReq struct {
	Csr    string `json:"csr" validate:"required,max=5500"`
	Status string `json:"status" validate:"required,CertificateStatusEnumType"`
}

type SignedFirmwareStatusNotificationReq struct {
	Status    string `json:"status" validate:"required,FirmwareStatusEnumType"`
	RequestId int    `json:"requestId,omitempty"`
}

type SignedUpdateFirmwareReq struct {
	Retries       *int    `json:"retries,omitempty" validate:"gte=0"`
	RetryInterval *int    `json:"retryInterval,omitempty" validate:"gte=0"`
	RequestId     int    `json:"requestId" validate:"required"`
	Firmware      string `json:"firmware" validate:"required,FirewareType"`
}
