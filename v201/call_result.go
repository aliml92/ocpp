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
	Status           string        	`json:"status" validate:"required,CertificateSignedStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType `json:"statusInfo,omitempty" ` 
}

type ClearCacheRes struct {
	Status           string        	`json:"status" validate:"required,ClearCacheStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType `json:"statusInfo,omitempty" `
}

type ClearChargingProfileRes struct {
	Status           string 	  	`json:"status" validate:"required,ClearChargingProfileStatusEnumType"` // todo: validation register required	
	StatusInfo       StatusInfoType `json:"statusInfo,omitempty" `
}

type ClearDisplayMessageRes struct {
	Status           string        	`json:"status" validate:"required,ClearMessageStatusEnumType"` // todo: validation register required
	StatusInfo       StatusInfoType `json:"statusInfo,omitempty" `
}


type ClearedChargingLimitRes struct {}


type ClearVariableMonitoringRes struct {
	ClearMonitoringResult  []ClearMonitoringResultType `json:"clearMonitoringResult" validate:"required,dive,required"`
}

type CostUpdatedRes struct {}

type CustomerInformationRes struct {
	Status 		 string        	`json:"status" validate:"required,CustomerInformationStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type DataTransferRes struct {
	Status 		 string        	`json:"status" validate:"required,DataTransferStatusEnumType"` // todo: validation register required
	Data         interface{}   	`json:"data,omitempty" `
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type DeleteCertificateRes struct {
	Status 		 string        	`json:"status" validate:"required,DeleteCertificateStatusEnumType"` // todo: validation register required
	StatusInfo   StatusInfoType 	`json:"statusInfo,omitempty" `
}

type FirmwareStatusNotificationRes struct {}

type Get15118EVCertificateRes struct {
	Status		string 			`json:"status" validate:"required,Iso15118EVCertificateStatusEnumType"` // todo 
	ExiResponse string			`json:"exiResponse" validate:"required,max=5600"`
	StatusInfo  StatusInfoType 	`json:"statusInfo,omitempty" `
}


type GetReportBaseResponse struct {
	Status 		string 			`json:"status" validate:"required,GenericDeviceModelStatusEnumType"`  // todo
	StatusInfo  StatusInfoType 	`json:"statusInfo,omitempty"`
}

type GetCertificateStatusResponse struct {
	Status 		string 			`json:"status" validate:"required,GetCertificateStatusEnumType"`  // todo
	OcspResult 	string 			`json:"ocspResult,omitempty" validate:"omitempty,max=5500"`
	StatusInfo  StatusInfoType 	`json:"statusInfo,omitempty"`
}


type GetChargingProfilesResponse struct {
	
}





//////////////

package temp


type GetCompositeScheduleRe struct {
	Status       string          			`json:"status" validate:"required,GenericStatusEnumType"`
	Schedule     CompositeScheduleType		`json:"schedule,omitempty"`
	StatusInfo   StatusInfoType				`json:"statusInfo,omitempty"`
}

type GetDisplayMessagesRes struct {
	Status       string          			`json:"status" validate:"required,GetDisplayMessagesStatusEnumType"`
	StatusInfo   StatusInfoType				`json:"statusInfo,omitempty"`
}

type GetInstalledCertificateIdsRes struct {
	Status       				string          						`json:"status" validate:"required,GetInstalledCertificateIdsStatusEnumType"`
	CertificateHashDataChain 	[]CertificateHashDataChainType			`json:"certificateHashDataChain,omitempty"`
	StatusInfo   				StatusInfoType							`json:"statusInfo,omitempty"`
}

type GetLocalListVersionRes struct {
	VersionNumber                *int                                    `json:"versionNumber" validate:"required"`
}

type GetLogRes struct {
	Status       string          			`json:"status" validate:"required,LogStatusEnumType"`
	Filename     string                     `json:"filename,omitempty"`
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