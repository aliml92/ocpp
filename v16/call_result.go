package v16

type AuthorizeConf struct {
	IdTagInfo IdTagInfo `json:"idTagInfo" validate:"required"`
}

type BootNotificationConf struct {
	CurrentTime string `json:"currentTime" validate:"required,ISO8601date"`
	Interval    int    `json:"interval" validate:"required,gte=0"`
	Status      string `json:"status" validate:"required,RegistrationStatus"`
}

type DataTransferConf struct {
	Status string `json:"status" validate:"required,DataTransferStatus"`
	Data   any    `json:"data,omitempty"`
}

type DiagnosticsStatusNotificationConf struct{}

type FirmwareStatusNotificationConf struct{}

type HeartbeatConf struct {
	CurrentTime string `json:"currentTime" validate:"required,ISO8601date"`
}

type MeterValuesConf struct{}

type StartTransactionConf struct {
	IdTagInfo     IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int       `json:"transactionId" validate:"required"`
}

type StatusNotificationConf struct{}

type StopTransactionConf struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo" validate:"omitempty"`
}

type CancelReservationConf struct {
	Status string `json:"status" validate:"required,CancelReservationStatus"`
}

type ChangeAvailabilityConf struct {
	Status string `json:"status" validate:"required,AvailabilityStatus"`
}

type ChangeConfigurationConf struct {
	Status string `json:"status" validate:"required,ConfigurationStatus"`
}

type ClearCacheConf struct {
	Status string `json:"status" validate:"required,ClearCacheStatus"`
}

type ClearChargingProfileConf struct {
	Status string `json:"status" validate:"required,ClearChargingProfileStatus"`
}

type GetCompositeScheduleConf struct {
	Status           string            `json:"status" validate:"required,GetCompositeScheduleStatus"`
	ConnectorId      *int              `json:"connectorId" validate:"omitempty,gte=0"`
	ScheduleStart    string            `json:"scheduleStart,omitempty" validate:"omitempty,ISO8601date"`
	ChargingSchedule *ChargingSchedule `json:"chargingSchedule,omitempty"`
}

type GetConfigurationConf struct {
	ConfigurationKey []KeyValue `json:"configurationKey,omitempty"`
	UnknownKey       []string   `json:"unknownKey,omitempty" validate:"omitempty,dive,max=50"`
}

type GetDiagnosticsConf struct {
	FileName string `json:"fileName,omitempty" validate:"omitempty,max=255"`
}

type GetLocalListVersionConf struct {
	ListVersion int `json:"listVersion" validate:"required,gte=0"`
}

type RemoteStartTransactionConf struct {
	Status string `json:"status" validate:"required,RemoteStartStopStatus"`
}

type RemoteStopTransactionConf struct {
	Status string `json:"status" validate:"required,RemoteStartStopStatus"`
}

type ReserveNowConf struct {
	Status string `json:"status" validate:"required,ReservationStatus"`
}

type ResetConf struct {
	Status string `json:"status" validate:"required,ResetStatus"`
}

type SendLocalListConf struct {
	Status string `json:"status" validate:"required,UpdateStatus"`
}

type SetChargingProfileConf struct {
	Status string `json:"status" validate:"required,ChargingProfileStatus"`
}

type TriggerMessageConf struct {
	Status string `json:"status" validate:"required,TriggerMessageStatus"`
}

type UnlockConnectorConf struct {
	Status string `json:"status" validate:"required,UnlockStatus"`
}

type UpdateFirmwareConf struct{}

// OCPP 1.6 security whitepaper edition 2 implementation

type CertificateSignedConf struct {
	Status string `json:"status" validate:"required,CertificateSignedStatusEnumType"`
}

type DeleteCertificateConf struct {
	Status string `json:"status" validate:"required,DeleteCertificateStatusEnumType"`
}

type ExtendedTriggerMessageConf struct {
	Status string `json:"status" validate:"required,TriggerMessageStatusEnumType"`
}

type GetInstalledCertificateIdsConf struct {
	Status              string                    `json:"status" validate:"required,GetInstalledCertificateStatusEnumType"`
	CertificateHashData []CertificateHashDataType `json:"certificateHashData,omitempty" validate:"omitempty,dive,required"`
}

type GetLogConf struct {
	Status   string `json:"status" validate:"required,LogStatusEnumType"`
	Filename string `json:"filename,omitempty" validate:"omitempty,max=255"`
}

type InstallCertificateConf struct {
	Status string `json:"status" validate:"required,CertificateStatusEnumType"`
}

type LogStatusNotificationConf struct{}

type SecurityEventNotificationConf struct{}

type SignCertificateConf struct {
	Status string `json:"status" validate:"required,GenericStatusEnumType"`
}

type SignedFirmwareStatusNotificationConf struct{}

type SignedUpdateFirmwareConf struct {
	Status string `json:"status" validate:"required,UpdateFirmwareStatusEnumType"`
}
