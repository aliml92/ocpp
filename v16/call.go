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


type StartTransactionReq struct {
	ConnectorId 			int 	`json:"connectorId" validate:"required,gt=0"`
	IdTag       			string 	`json:"idTag" validate:"required,max=20"`
	MeterStart  			int 	`json:"meterStart" validate:"required,gte=0"`
	ReservationId 			int 	`json:"reservationId,omitempty" validate:"omitempty,gt=0"`
	Timestamp   			string 	`json:"timestamp" validate:"required,ISO8601date"`
}











type ChangeAvailabilityReq struct {
	ConnectorId 		int 				`json:"connectorId" validate:"required,gte=0"`
	Type 				string 				`json:"type" validate:"required"`
}