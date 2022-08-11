package v201


type IdTokenType struct {
	IdToken 		string 					`json:"idToken" validate:"required,max=36"`
	Type    		string 					`json:"type" validate:"required,IdTokenEnumType"`  // todo: validation register required
	AdditionalInfo  []AdditionalInfoType 	`json:"additionalInfo" validate:"required,dive,required"`
}

type AdditionalInfoType struct {
	AdditionalToken string 					`json:"additionalToken" validate:"required,max=36"`
	Type 			string 					`json:"type" validate:"required,max=50"`  
}


type OCSPRequestDataType struct {
	HashAlgorithm 	string 					`json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"` // todo: validation register required
	IssuerNameHash 	string 					`json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash 	string 					`json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber 	string 					`json:"serialNumber" validate:"required,max=40"`
	ResponderUrl 	string 					`json:"responderUrl" validate:"required,max=512"`
}

type ChargingStationType struct {
	SerialNumber    string 					`json:"serialNumber,omitempty" validate:"omitempy,max=25"`
	Model           string 					`json:"model" validate:"required,max=20"`
	VendorName      string 					`json:"vendorName" validate:"required,max=50"`
	FirmwareVersion string 					`json:"firmwareVersion,omitempty" validate:"omitempty,max=50"`
	Modem 			ModemType 				`json:"modem,omitempty" validate:"dive,required"`		
}

type ModemType struct {
	Iccid 			string 					`json:"iccid,omitempty" validate:"omitempty,max=20"`
	Imsi 			string 					`json:"imsi,omitempty" validate:"omitempty,max=20"`
}

type StatusInfoType struct {
	ReasonCode      string 					`json:"reasonCode" validate:"required,max=20"`
	AdditionalInfo  string 					`json:"additionalInfo" validate:"required,max=512"`
}

type EVSEType struct {
	Id              *int                     `json:"id" validate:"required,gt=0"`
	ConnectorId     *int                     `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
}


type ClearChargingProfileType struct {
	EvseId 					*int             `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfilePurpose  string 		     `json:"chargingProfilePurpose,omitempty" validate:"omitempty,ChargingProfilePurposeEnumType"` // todo: validation register required
	StackLevel 				*int             `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
}

type ClearMonitoringResultType struct {
	Status                  string 			 `json:"status" validate:"required,ClearMonitoringStatusEnumType"` // todo: validation register required 
	Id 						*int             `json:"id" validate:"required"`
	StatusInfo 				StatusInfoType   `json:"statusInfo,omitempty" `
}

type CertificateHashDataType struct {
	HashAlgorithm 	string 					`json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"` // todo: validation register required
	IssuerNameHash 	string 					`json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash 	string 					`json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber 	string 					`json:"serialNumber" validate:"required,max=40"`
}


type ChargingProfileCriterionType struct {
	
}