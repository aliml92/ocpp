package v201




type ACChargingParametersType struct {
	EnergyAmount      int     `json:"energyAmount" validate:"required"`
	EvMinCurrent	  int     `json:"evMinCurrent" validate:"required"`
	EvMaxCurrent	  int     `json:"evMaxCurrent" validate:"required"`
	EvMaxVoltage	  int     `json:"evMaxVoltage" validate:"required"`
}


type AdditionalInfoType struct {
	AdditionalToken string 					`json:"additionalToken" validate:"required,max=36"`
	Type 			string 					`json:"type" validate:"required,max=50"`  
}


type APNType struct {
	Apn 					string 				`json:"apn" validate:"required,max=100"`
	ApnUserName 			string 				`json:"apnUserName,omitempty" validate:"omitempty,max=20"`
	ApnPassword 			string 				`json:"apnPassword,omitempty" validate:"omitempty,max=20"`
	SimPin              	int 				`json:"simPin,omitempty"`
	PreferredNetwork    	string 				`json:"preferredNetwork,omitempty" validate:"omitempty,max=6"`
	UseOnlyPreferredNetwork bool 				`json:"useOnlyPreferredNetwork,omitempty"`
	ApnAuthentication    	string 				`json:"apnAuthentication" validate:"required,APNAuthenticationEnumType"`
}


type AuthorizarionData struct {
	IdTokenInfo            IdTokenInfoType 		`json:"idTokenInfo,omitempty"`
	IdToken                IdTokenType 			`json:"idToken" validate:"required"`
}


type CertificateHashDataChainType struct {
	CertificateType    			string   					`json:"certificateType" validate:"required,GetCertificateIdUseEnumType"`
	CertificateHashData			CertificateHashDataType 	`json:"certificateHashData" validate:"required"`
	ChildCertificateHashData 	[]CertificateHashDataType 	`json:"childCertificateHashData,omitempty" validate:"omitempty,max=4,dive,required"`  
}

type CertificateHashDataType struct {
	HashAlgorithm 	string 					`json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"` // todo: validation register required
	IssuerNameHash 	string 					`json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash 	string 					`json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber 	string 					`json:"serialNumber" validate:"required,max=40"`
}

type ChargingLimitType struct {
	ChargingLimitSource    	string 					`json:"chargingLimitSource" validate:"required,ChargingLimitSourceEnumType"`
	IsGridCritical 			bool 					`json:"isGridCritical,omitempty"`
}


type ChargingNeedsType struct {
	RequestedEnergyTransfer      string 			   		`json:"requestedEnergyTransfer" validate:"required,RequestedEnergyTransferEnumType"`
	DepartureTime				 string 			   		`json:"departureTime" validate:"omitempty,ISO8601date"`
	ACChargingParameters	     ACChargingParametersType 	`json:"acChargingParameters,omitempty"`
	DcChargingParameters	     DCChargingParametersType 	`json:"dcChargingParameters,omitempty"`
}


type ChargingProfileCriterionType struct {
	ChargingProfilePurpose 		string 						`json:"chargingProfilePurpose,omitempty" validate:"omitempty,ChargingProfilePurposeEnumType"`
	StackLevel					*int 						`json:"stackLevel,omitempty"`
	ChargingProfileId           []int                   	`json:"chargingProfileId,omitempty" validate:"omitempty,dive,required"`
	ChargingLimitSource         []string 					`json:"chargingLimitSource,omitempty" validate:"omitempty,max=4,dive,required,ChargingLimitSourceEnumType"`
}



type ChargingProfileType struct {
	Id 							*int 						`json:"id" validate:"required"`
	StackLevel                  *int 						`json:"stackLevel" validate:"required"`
	ChargingProfilePurpose      string 						`json:"chargingProfilePurpose" validate:"required,ChargingProfilePurposeEnumType"`
	ChargingProfileKind         string 						`json:"chargingProfileKind" validate:"required,ChargingProfileKindEnumType"`
	RecurrencyKind			    string 						`json:"recurrencyKind,omitempty" validate:"omitempty,RecurrencyKindEnumType"`
	ValidFrom                   string 						`json:"validFrom,omitempty" validate:"omitempty,ISO8601date"`
	ValidTo                     string 						`json:"validTo,omitempty" validate:"omitempty,ISO8601date"`
	TransactionId               string 						`json:"transactionId,omitempty" validate:"omitempty,max=36"`
	ChargingSchedule            ChargingScheduleType 		`json:"chargingSchedule" validate:"required,min=1,max=3,dive,required"` 
}


type ChargingSchedulePeriodType struct {
	StartPeriod 				*int 						`json:"startPeriod" validate:"required"`
	Limit 						float32 					`json:"limit" validate:"required"`
	NumberPhases 				int 						`json:"numberPhases,omitempty"`
	PhaseToUse 					int 						`json:"phaseToUse,omitempty"`	
}


type ChargingScheduleType struct {
	Id                  	*int 						`json:"id" validate:"required"`
	StartSchedule       	string 						`json:"startSchedule,omitempty" validate:"omitempty,ISO8601date"`
	Duration				*int 						`json:"duration,omitempty"`
	ChargingRateUnit    	string 						`json:"chargingRateUnit" validate:"required,ChargingRateUnitEnumType"`
	MinChargingRate     	float32 					`json:"minChargingRate,omitempty"`
	ChargingSchedulePeriod []ChargingSchedulePeriodType `json:"chargingSchedulePeriod" validate:"required,min=1,max=1024,dive,required"`
	SalesTariff             SalesTariffType             `json:"salesTariff,omitempty"`
}


type ChargingStationType struct {
	SerialNumber    string 					`json:"serialNumber,omitempty" validate:"omitempy,max=25"`
	Model           string 					`json:"model" validate:"required,max=20"`
	VendorName      string 					`json:"vendorName" validate:"required,max=50"`
	FirmwareVersion string 					`json:"firmwareVersion,omitempty" validate:"omitempty,max=50"`
	Modem 			ModemType 				`json:"modem,omitempty"`		
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



type ComponentType struct {
	Name 					string 			 `json:"name" validate:"required,max=50"`
	Instance 				string 			 `json:"instance,omitempty" validate:"omitempty,max=50"`
	Evse 					EVSEType 		 `json:"evse,omitempty"`
}


type ComponentVariableType struct {
	Component 				ComponentType 	 `json:"component" validate:"required"`
	Variable 				VariableType 	 `json:"variable,omitempty"`
}


type CompositeScheduleType struct {
	EvseId                	*int             				`json:"evseId" validate:"required,gte=0"`
	Duration              	*int             				`json:"duration" validate:"required`
	ScheduleStart         	string 		   					`json:"scheduleStart" validate:"required,ISO8601date"`
	ChargingRateUnit      	string 		   					`json:"chargingRateUnit" validate:"required,ChargingRateUnitEnumType"`
	ChargingSchedulePeriod	[]ChargingSchedulePeriodType 	`json:"chargingSchedulePeriod" validate:"required,min=1,dive,required"`
}




type ConsumptionCostType struct {
	StartValue           float32 					`json:"startValue" validate:"required"`
	Cost                 CostType 					`json:"cost" validate:"required,min=1,max=3,dive,required"`
}


type CostType struct {
	CostKind			  string 					`json:"costKind" validate:"required,CostKindEnumType"`
	Amount                int 						`json:"amount" validate:"required"`
	AmountMultiplier      int 						`json:"amountMultiplier,omitempty"`
}


type DCChargingParametersType struct {
	EvMaxCurrent          int               `json:"evMaxCurrent" validate:"required"`
	EvMaxVoltage		  int               `json:"evMaxVoltage" validate:"required"`
	EnergyAmount		  int               `json:"energyAmount,omitempty"`
	EvMaxPower            int               `json:"evMaxPower,omitempty"`
	StateOfCharge		  *int              `json:"stateOfCharge,omitempty"`
	EvEnergyCapacity	  int               `json:"evEnergyCapacity,omitempty"`
	FullSoC				  *int              `json:"fullSoC,omitempty"`
	BulkSoC				  *int              `json:"bulkSoC,omitempty"` 
}



type EventDataType struct {
	EventId               *int               `json:"eventId" validate:"required"`
	Timestamp 		      string             `json:"timestamp" validate:"required,ISO8601date"`
	Trigger               string             `json:"trigger" validate:"required,TriggerEnumType"`
	Cause                 *int               `json:"cause,omitempty"`
	ActualValue 		  string 		     `json:"actualValue" validate:"required,max=2500"`
	TechCode              string             `json:"techCode,omitempty" validate:"omitempty,max=50"`
	TechInfo              string             `json:"techInfo,omitempty" validate:"omitempty,max=500"`
	Cleared 			  bool               `json:"cleared,omitempty"`
	TransactionId	      string             `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	VariableMonitoringId  *int               `json:"variableMonitoringId,omitempty"`
	EventNotificationType string             `json:"eventNotificationType" validate:"required,EventNotificationTypeEnumType"`
	Component             ComponentType      `json:"component" validate:"required"`
	Variable              VariableType       `json:"variable" validate:"required"`
}


type EVSEType struct {
	Id              *int                     `json:"id" validate:"required,gt=0"`
	ConnectorId     *int                     `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
}


type FirmwareType  struct {
	Location       		string 					`json:"location" validate:"required,max=512"`
	RetrieveDateTime	string 					`json:"retrieveDateTime" validate:"required,ISO8601date"`
	InstallDateTime     string 					`json:"installDateTime,omitempty" validate:"omitempty,ISO8601date"`
	SigningCertificate  string 					`json:"signingCertificate,omitempty" validate:"omitempty,max=512"`
	Signature           string 					`json:"signature,omitempty" validate:"omitempty,max=800"`
}

type IdTokenType struct {
	IdToken 		string 					`json:"idToken" validate:"required,max=36"`
	Type    		string 					`json:"type" validate:"required,IdTokenEnumType"`  // todo: validation register required
	AdditionalInfo  []AdditionalInfoType 	`json:"additionalInfo" validate:"required,dive,required"`
}




type OCSPRequestDataType struct {
	HashAlgorithm 	string 					`json:"hashAlgorithm" validate:"required,HashAlgorithmEnumType"` // todo: validation register required
	IssuerNameHash 	string 					`json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash 	string 					`json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber 	string 					`json:"serialNumber" validate:"required,max=40"`
	ResponderUrl 	string 					`json:"responderUrl" validate:"required,max=512"`
}



type ModemType struct {
	Iccid 			string 					`json:"iccid,omitempty" validate:"omitempty,max=20"`
	Imsi 			string 					`json:"imsi,omitempty" validate:"omitempty,max=20"`
}

type StatusInfoType struct {
	ReasonCode      string 					`json:"reasonCode" validate:"required,max=20"`
	AdditionalInfo  string 					`json:"additionalInfo" validate:"required,max=512"`
}








