package v201


type AuthorizeRes struct {
	CertificateStatus string 		`json:"certificateStatus,omitempty" validate:"omitempty,AuthorizeCertificateStatusEnumType"` // todo: validation register required
	IdTokenInfo       IdTokenType   `json:"idTokenInfo" validate:"required"`
}


type BootNotificationRes struct {
	CurrentTime 	  string 			`json:"currentTime" validate:"required,ISO8601date"`  // todo: validation regis
	Interval          *int 				`json:"interval" validate:"required,gte=0"`
	Status            string        	`json:"status" validate:"required,RegistrationStatusEnumType"` // todo: validation register required
	StatusInfo        StatusInfoType 	`json:"statusInfo,omitempty" `
}


type CancelReservationRes struct {
	Status 		      string        	`json:"status" validate:"required,CancelReservationStatusEnumType"` // todo: validation register required
	StatusInfo        StatusInfoType 	`json:"statusInfo,omitempty" `    
}


type CertificateSignedRes struct {
	Status           string        		`json:"status" validate:"required,CertificateSignedStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType 	`json:"statusInfo,omitempty" ` 
}


type ChangeAvailabilityRes struct {
	Status           string        		`json:"status" validate:"required,ChangeAvailabilityStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType 	`json:"statusInfo,omitempty" `
}

type ClearCacheRes struct {
	Status           string        		`json:"status" validate:"required,ClearCacheStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType 	`json:"statusInfo,omitempty" `
}

type ClearChargingProfileRes struct {
	Status           string 	  		`json:"status" validate:"required,ClearChargingProfileStatusEnumType"` // todo: validation register required	
	StatusInfo       StatusInfoType 	`json:"statusInfo,omitempty" `
}

type ClearDisplayMessageRes struct {
	Status           string        		`json:"status" validate:"required,ClearMessageStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType 	`json:"statusInfo,omitempty" `
}


type ClearedChargingLimitRes struct {}


type ClearVariableMonitoringRes struct {
	ClearMonitoringResult  []ClearMonitoringResultType `json:"clearMonitoringResult" validate:"required,dive,required"`
}

type CostUpdatedRes struct {}

type CustomerInformationRes struct {
	Status 		 string        		`json:"status" validate:"required,CustomerInformationStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type DataTransferRes struct {
	Status 		 string        		`json:"status" validate:"required,DataTransferStatusEnumType"` // todo: validation register required
	Data         interface{}   		`json:"data,omitempty" `
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type DeleteCertificateRes struct {
	Status 		 string        		`json:"status" validate:"required,DeleteCertificateStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type FirmwareStatusNotificationRes struct {}

type Get15118EVCertificateRes struct {
	Status		string 			`json:"status" validate:"required,Iso15118EVCertificateStatusEnumType"` // todo 
	ExiResponse string			`json:"exiResponse" validate:"required,max=5600"`
	StatusInfo  StatusInfoType 	`json:"statusInfo,omitempty" `
}


type GetReportBaseRes struct {
	Status 		string 			`json:"status" validate:"required,GenericDeviceModelStatusEnumType"`  // todo
	StatusInfo  StatusInfoType 	`json:"statusInfo,omitempty"`
}

type GetCertificateStatusRes struct {
	Status 		string 			`json:"status" validate:"required,GetCertificateStatusEnumType"`  // todo
	OcspResult 	string 			`json:"ocspResult,omitempty" validate:"omitempty,max=5500"`
	StatusInfo  StatusInfoType 	`json:"statusInfo,omitempty"`
}


type GetChargingProfilesRes struct {
	Status 				string 					`json:"status" validate:"required,GetChargingProfilesStatusEnumType"`  // todo
	StatusInfo 			StatusInfoType 			`json:"statusInfo,omitempty"`
}


type GetCompositeScheduleRes struct {
	Status       string          				`json:"status" validate:"required,GenericStatusEnumType"`
	Schedule     CompositeScheduleType			`json:"schedule,omitempty"`
	StatusInfo   StatusInfoType					`json:"statusInfo,omitempty"`
}


type GetDisplayMessagesRes struct {
	Status       string          			`json:"status" validate:"required,GetDisplayMessagesStatusEnumType"`
	StatusInfo   StatusInfoType				`json:"statusInfo,omitempty"`
}


type GetInstalledCertificateIdsRes struct {
	Status       				string          						`json:"status" validate:"required,GetInstalledCertificateIdsStatusEnumType"`
	CertificateHashDataChain 	[]CertificateHashDataChainType			`json:"certificateHashDataChain,omitempty" validate:"omitempty,dive,required"`
	StatusInfo   				StatusInfoType							`json:"statusInfo,omitempty"`
}

type GetLocalListVersionRes struct {
VersionNumber                *int				`json:"versionNumber" validate:"required"`
}

type GetLogRes struct {
	Status       string          			`json:"status" validate:"required,LogStatusEnumType"`
	Filename     string                     `json:"filename,omitempty" validate:"omitempty,max=255"`
	StatusInfo   StatusInfoType				`json:"statusInfo,omitempty"`   
}


type GetMonitoringReportRes struct {
	Status       string          			`json:"status" validate:"required,GenericDeviceModelStatusEnumType"`
	StatusInfo   StatusInfoType				`json:"statusInfo,omitempty"`  
}


type GetReportRes struct {
	Status       string          			`json:"status" validate:"required,GenericDeviceModelStatusEnumType"`
	StatusInfo   StatusInfoType				`json:"statusInfo,omitempty"`
}


type GetTransactionStatusRes struct {
	OngoingIndicator       bool 			`json:"ongoingIndicator,omitempty"`
	MessagesInQueue        bool 			`json:"messagesInQueue" validate:"required"`
}

type GetVariablesRes struct {
	GetVariableResult      []GetVariableResultType  `json:"getVariableResult" validate:"required,dive,required"`  
}


type HeartbeatRes struct {
	CurrentTime 		 string 			`json:"currentTime" validate:"required,ISO8601date"`
}


type InstallCertificateRes struct {
	Status 		 string        		`json:"status" validate:"required,InstallCertificateStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type LogStatusNotificationRes struct {}


type MeterValuesRes struct {}


type NotifyChargingLimitRes struct {}

type NotifyCustomerInformationRes struct {}


type NotifyDisplayMessageRes struct {}


type NotifyEVChargingNeedsRes struct {
	Status 		 string        		`json:"status" validate:"required,NotifyEVChargingNeedsStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}


type NotifyEVChargingScheduleRes struct {
	Status 		 string        		`json:"status" validate:"required,GenericStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}


type NotifyEventRes struct {}

type NotifyMonitoringReportRes struct {}


type NotifyReportRes struct {}


type PublishFirmwareRes struct {
	Status 		 string        		`json:"status" validate:"required,GenericStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type PublishFirmwareStatusNotificationRes struct {}

type ReportChargingProfilesRes struct {}


type RequestStartTransactionRes struct {
	Status 		   string        		`json:"status" validate:"required,RequestStartStopStatusEnumType"` // todo: validation register required
	TransactionId  string 	  			`json:"transactionId,omitempty" validate:"omitempty,max=36"`
	StatusInfo     StatusInfoType 	`json:"statusInfo,omitempty" `
}

type RequestStopTransactionRes struct {
	Status 		   string        		`json:"status" validate:"required,RequestStartStopStatusEnumType"` // todo: validation register required
	StatusInfo    StatusInfoType 		`json:"statusInfo,omitempty"`
}

type ReservationStatusUpdateRes struct {}


type ReserveNowRes struct {
	Status 		 string        		`json:"status" validate:"required,ReserveNowStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type ResetRes struct {
	Status 		 string        		`json:"status" validate:"required,ResetStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type SecurityEventNotificationRes struct {}


type SendLocalListRes struct {
	Status 		 string        		`json:"status" validate:"required,SendLocalListStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type SetChargingProfileRes struct {
	Status 		 string        		`json:"status" validate:"required,ChargingProfileStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}


type SetDisplayMessageRes struct {
	Status 		 string        		`json:"status" validate:"required,DisplayMessageStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}


type SetMonitoringBaseRes struct {
	Status 		 string        		`json:"status" validate:"required,GenericDeviceModelStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type SetMonigoringLevelRes struct {
	Status 		 string        		`json:"status" validate:"required,GenericStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type SetNetworkProfileRes struct {
	Status 		 string        		`json:"status" validate:"required,SetNetworkProfileStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}


type SetVariableMonitoringRes struct {
	SetMonitoringResult []SetMonitoringResultType `json:"setMonitoringResult" validate:"required,dive,required"`
}


type SetVariableRes struct {
	SetVariableResult []SetVariableResultType `json:"setVariableResult" validate:"required,dive,required"`
}


type SignCertificateRes struct {
	Status 		 string        		`json:"status" validate:"required,GenericStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}


type StatusNotificationRes struct {}


type TransactionEventRes struct {
	TotalCost        			float32 			`json:"totalCost,omitempty" validate:"omitempty,min=0"`
	ChargingPriority 			*int 				`json:"chargingPriority,omitempty" validate:"omitempty,gte=-9,lte=9"`
	IdTokenInfo	  	    		IdTokenInfoType 	`json:"idTokenInfo,omitempty"`
	UpdatedPersonalMessage 		MessageContentType 	`json:"updatedPersonalMessage,omitempty"` 
}


type TriggerMessageRes struct {
	Status 		 string        		`json:"status" validate:"required,TriggerMessageStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type UnlockConnectorRes struct {
	Status 		 string        		`json:"status" validate:"required,UnlockStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type UnpublishFirmwareRes struct {
	Status 		 string        		`json:"status" validate:"required,UnpublishFirmwareStatusEnumType"` // todo: validation register required
}

type UpdateFirmwareRes struct {
	Status 		 string        		`json:"status" validate:"required,UpdateFirmwareStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}