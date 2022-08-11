package v201


type AuthorizeReq struct {
	Certificate 				string 				  `json:"certificate,omitempty" validate:"max=5500"`
	IdToken     				IdTokenType 		  `json:"idTokenType" validate:"dive,required"`
	Iso15118CertificateHashData []OCSPRequestDataType `json:"iso15118CertificateHashData,omitempty" validate:"dive,required"`
}

type BootNotificationReq struct {
	Reason 						string 				  `json:"reason" validate:"required,BootReasonEnumType"` // todo: validation register required
	ChargingStation 			ChargingStationType   `json:"chargingStation" validate:"required"`
}

type CancelReservationReq struct {
	ReservationId 				*int 				  `json:"reservationId" validate:"required"`

}

type CertificateSignedReq struct {
	CertificateChain  string 			`json:"certificateChain" validate:"required,max=10000"`
	CertificateType   string 			`json:"certificateType,omitempty" validate:"CertificateSigningUseEnumType"` // todo: validation register required
}  


type ChangeAvailabilityReq struct {
	OperationalStatus string  			`json:"operationalStatus" validate:"required,OperationalStatusEnumType"` // todo: validation register required	
	Evse 			  EVSEType 			`json:"evse,omitempty"` 
}

type ClearCacheReq struct {}

type ClearChargingProfileReq struct {
	ChargingProfileId  		*int 						`json:"chargingProfileId,omitempty"`
	ChargingProfileCriteria ClearChargingProfileType 	`json:"chargingProfileCriteria,omitempty"`
}

type ClearDisplayMessageReq struct {
	Id  					*int 						`json:"id" validate:"required"`
}

type ClearedChargingLimitReq struct {
	ChargingLimitSource     string        	`json:"chargingLimitSource" validate:"required,ChargingLimitSourceEnumType"` // todo: validation register required
	EvseId    		        *int 			`json:"evseId,omitempty"`
}


type ClearVariableMonitoringReq struct {
	Id 					[]int 						`json:"id" validate:"required"`
}


type CostUpdatedReq struct {
	TotalCost 		     float32 					`json:"totalCost" validate:"required"`
	TransactionId		 *int 						`json:"transactionId" validate:"required,max=36"`
}

type CustomerInformationReq struct {
	RequestId 				*int 						`json:"requestId" validate:"required"`
	Report                  bool 						`json:"report" validate:"required"`
	Clear 				 	bool 						`json:"clear" validate:"required"`
	CustomerIdentifier 		string 						`json:"customerIdentifier" validate:"required,max=64"`
	IdToken 				IdTokenType 				`json:"idTokenType,omitempty"`
	CustomerCertificate		CertificateHashDataType 	`json:"customerCertificate,omitempty" `
}


type DataTransferReq struct {
	MessageId 	            string 						`json:"messageId,omitempty" validate:"omitempty,max=50"`
	Data 				    interface{} 				`json:"data,omitempty"`
	VendorId 				string 						`json:"vendorId" validate:"required,max=255"`
}


type DeleteCertificateReq struct {
	CertificateHashData 		CertificateHashDataType 	`json:"certificateHashData" validate:"required"`
}


type FirmwareStatusNotificationReq struct {
	Status 						string 				  `json:"status" validate:"required,FirmwareStatusEnumType"` // todo: validation register required
	RequestId 					*int 				  `json:"requestId,omitempty" `
}

type Get15118EVCertificateReq struct {
	Iso15118SchemaVersion       string                 `json:"iso15118SchemaVersion" validate:"required,max=50"`
	Action                      string                 `json:"action" validate:"required,CertificateActionEnumType"` // todo: validation register required
	ExiRequest                  string                 `json:"exiRequest" validate:"required,max=5600"`
}

type GetBaseReportRequest struct {
	RequestId 					int 					`json:"requestId"`
	ReportBase					string 					`json:"reportBase" validate:"required,ReportBaseEnumType"`  // todo
}

type GetCertificateStatusRequest struct {
	OcspRequestData 	OCSPRequestDataType 			`json:"ocspRequestData" validate:"required"`
}

type GetChargingProfilesRequest struct {
	RequestId 					int 						 `json:"requestId"`
	EvseId 						int 						 `json:"evseId,omitempty"`
	ChargingProfile 			ChargingProfileCriterionType `json:"chargingProfile"`
}



