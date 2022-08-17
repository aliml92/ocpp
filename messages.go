package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aliml92/ocpp/v16"
	// "github.com/aliml92/ocpp/v201"
	"github.com/google/uuid"
)


var reqmap = map[string]func(*json.RawMessage) (Payload, error){
	"BootNotification":              ureq[v16.BootNotificationReq],
	"Authorize":                     ureq[v16.AuthorizeReq],
	"DataTransfer":                  ureq[v16.DataTransferReq],
	"DiagnosticsStatusNotification": ureq[v16.DiagnosticsStatusNotificationReq],
	"FirmwareStatusNotification":    ureq[v16.FirmwareStatusNotificationReq],
	"Heartbeat":                     ureq[v16.HeartbeatReq],
	"MeterValues":                   ureq[v16.MeterValuesReq],
	"StartTransaction":              ureq[v16.StartTransactionReq],
	"StatusNotification":            ureq[v16.StatusNotificationReq],
	"StopTransaction":               ureq[v16.StopTransactionReq],
	"CanCelReservation":             ureq[v16.CancelReservationReq],
	"ChangeAvailability":            ureq[v16.ChangeAvailabilityReq],
	"ChangeConfiguration":           ureq[v16.ChangeConfigurationReq],
	"ClearCache":                    ureq[v16.ClearCacheReq],
	"ClearChargingProfile":          ureq[v16.ClearChargingProfileReq],
	"GetCompositeSchedule":          ureq[v16.GetCompositeScheduleReq],
	"GetConfiguration":              ureq[v16.GetConfigurationReq],
	"GetDiagnostics":                ureq[v16.GetDiagnosticsReq],
	"GetLocalListVersion":           ureq[v16.GetLocalListVersionReq],
	"RemoteStartTransaction":        ureq[v16.RemoteStartTransactionReq],
	"RemoteStopTransaction":         ureq[v16.RemoteStopTransactionReq],
	"ReserveNow":                    ureq[v16.ReserveNowReq],
	"Reset":                         ureq[v16.ResetReq],
	"SendLocalList":                 ureq[v16.SendLocalListReq],
	"SetChargingProfile":            ureq[v16.SetChargingProfileReq],
	"TriggerMessage":                ureq[v16.TriggerMessageReq],
	"UnlockConnector":               ureq[v16.UnlockConnectorReq],
	"UpdateFirmware":                ureq[v16.UpdateFirmwareReq],
}



var confmap = map[string]func(*json.RawMessage) (Payload, error){
	"BootNotification":              uconf[v16.BootNotificationConf],
	"Authorize":                     uconf[v16.AuthorizeConf],
	"DataTransfer":                  uconf[v16.DataTransferConf],
	"DiagnosticsStatusNotification": uconf[v16.DiagnosticsStatusNotificationConf],
	"FirmwareStatusNotification":    uconf[v16.FirmwareStatusNotificationConf],
	"Heartbeat":                     uconf[v16.HeartbeatConf],
	"MeterValues":                   uconf[v16.MeterValuesConf],
	"StartTransaction":              uconf[v16.StartTransactionConf],
	"StatusNotification":            uconf[v16.StatusNotificationConf],
	"StopTransaction":               uconf[v16.StopTransactionConf],
	"CanCelReservation":             uconf[v16.CancelReservationConf],
	"ChangeAvailability":            uconf[v16.ChangeAvailabilityConf],
	"ChangeConfiguration":           uconf[v16.ChangeConfigurationConf],
	"ClearCache":                    uconf[v16.ClearCacheConf],
	"ClearChargingProfile":          uconf[v16.ClearChargingProfileConf],
	"GetCompositeSchedule":          uconf[v16.GetCompositeScheduleConf],
	"GetConfiguration":              uconf[v16.GetConfigurationConf],
	"GetDiagnostics":                uconf[v16.GetDiagnosticsConf],
	"GetLocalListVersion":           uconf[v16.GetLocalListVersionConf],
	"RemoteStartTransaction":        uconf[v16.RemoteStartTransactionConf],
	"RemoteStopTransaction":         uconf[v16.RemoteStopTransactionConf],
	"ReserveNow":                    uconf[v16.ReserveNowConf],
	"Reset":                         uconf[v16.ResetConf],
	"SendLocalList":                 uconf[v16.SendLocalListConf],
	"SetChargingProfile":            uconf[v16.SetChargingProfileConf],
	"TriggerMessage":                uconf[v16.TriggerMessageConf],
	"UnlockConnector":               uconf[v16.UnlockConnectorConf],
	"UpdateFirmware":                uconf[v16.UpdateFirmwareConf],
}




// var reqmapv201 = map[string]func(*json.RawMessage) (Payload, error){
// 	"Authorize":                     	ureq[v201.AuthorizeReq],
// 	"BootNotification":              	ureq[v201.BootNotificationReq],
// 	"CancelReservation":             	ureq[v201.CancelReservationReq],
// 	"CertificateSigned":			 	ureq[v201.CertificateSignedReq],
// 	"ChangeAvailability":            	ureq[v201.ChangeAvailabilityReq],
// 	"ClearCache":                    	ureq[v201.ClearCacheReq],
// 	"ClearChargingProfile":          	ureq[v201.ClearChargingProfileReq],
// 	"ClearDisplayMessage":           	ureq[v201.ClearDisplayMessageReq],
// 	"ClearedChargingLimit":          	ureq[v201.ClearedChargingLimitReq],
// 	"ClearVariableMonitoring":       	ureq[v201.ClearVariableMonitoringReq],
// 	"CostUpdated":                   	ureq[v201.CostUpdatedReq],
// 	"CustomerInformation":           	ureq[v201.CustomerInformationReq],
// 	"DataTransfer":                  	ureq[v201.DataTransferReq],
// 	"DeleteCertificate":             	ureq[v201.DeleteCertificateReq],
// 	"FirmwareStatusNotification":    	ureq[v201.FirmwareStatusNotificationReq],
// 	"Get15118EVCertificate":         	ureq[v201.Get15118EVCertificateReq],
// 	"GetBaseReport":                 	ureq[v201.GetBaseReportReq],
// 	"GetCertificateStatus":          	ureq[v201.GetCertificateStatusReq],
// 	"GetChargingProfiles":            	ureq[v201.GetChargingProfilesReq],
// 	"GetCompositeSchedule":          	ureq[v201.GetCompositeScheduleReq],
// 	"GetDisplayMessages":            	ureq[v201.GetDisplayMessagesReq],
// 	"GetInstalledCertificateIds":      	ureq[v201.GetInstalledCertificateIdsReq],
// 	"GetLocalListVersion":           	ureq[v201.GetLocalListVersionReq],
// 	"GetLog":							ureq[v201.GetLogReq],
// 	"GetMonitoringReport":           	ureq[v201.GetMonitoringReportReq],
// 	"GetReport":                     	ureq[v201.GetReportReq],
// 	"GetTransactionStatus":          	ureq[v201.GetTransactionStatusReq],
// 	"GetVariables":                  	ureq[v201.GetVariablesReq],
// 	"Heartbeat":                     	ureq[v201.HeartbeatReq],
// 	"InstallCertificate":            	ureq[v201.InstallCertificateReq],
// 	"LogStatusNotification":         	ureq[v201.LogStatusNotificationReq],
// 	"MeterValues":                   	ureq[v201.MeterValuesReq],
// 	"NotifyChargingLimit":           	ureq[v201.NotifyChargingLimitReq],
// 	"NotifyCustomerInformation":     	ureq[v201.NotifyCustomerInformationReq],
// 	"NotifyDisplayMessages":         	ureq[v201.NotifyDisplayMessagesReq],
// 	"NotifyEVChargingNeeds":         	ureq[v201.NotifyEVChargingNeedsReq],
// 	"NotifyEVChargingSchedule":        	ureq[v201.NotifyEVChargingScheduleReq],
// 	"NotifyEvent":                   	ureq[v201.NotifyEventReq],
// 	"NotifyMonitoringReport":        	ureq[v201.NotifyMonitoringReportReq],
// 	"NotifyReport":                  	ureq[v201.NotifyReportReq],
// 	"PublishFirmware":               	ureq[v201.PublishFirmwareReq],
// 	"PublishFirmawareStatusNotification": ureq[v201.PublishFirmwareStatusNotificationReq],
// 	"ReportChargingProfiles":        	ureq[v201.ReportChargingProfilesReq],
// 	"RequestStartTransaction":       	ureq[v201.RequestStartTransactionReq],
// 	"RequestStopTransaction":        	ureq[v201.RequestStopTransactionReq],
// 	"ReservationStatusUpdate":       	ureq[v201.ReservationStatusUpdateReq],
// 	"ReserveNow":                    	ureq[v201.ReserveNowReq],
// 	"Reset":                         	ureq[v201.ResetReq],
// 	"SecurityEventNotification":     	ureq[v201.SecurityEventNotificationReq],
// 	"SendLocalList":                 	ureq[v201.SendLocalListReq],
// 	"SetChargingProfile":            	ureq[v201.SetChargingProfileReq],
// 	"SetDisplayMessage":             	ureq[v201.SetDisplayMessageReq],
// 	"SetMonitoringBase":             	ureq[v201.SetMonitoringBaseReq],
// 	"SetMonitoringLevel":            	ureq[v201.SetMonitoringLevelReq],
// 	"SetNetworkProfile":             	ureq[v201.SetNetworkProfileReq],
// 	"SetVariableMonitoring":         	ureq[v201.SetVariableMonitoringReq],
// 	"SetVariables":                  	ureq[v201.SetVariablesReq],
// 	"SignCertificate":               	ureq[v201.SignCertificateReq],
// 	"StatusNotification":            	ureq[v201.StatusNotificationReq],
// 	"TransactionEvent":			  		ureq[v201.TransactionEventReq],
// 	"TriggerMessage":                	ureq[v201.TriggerMessageReq],
// 	"UnlockConnector":               	ureq[v201.UnlockConnectorReq],
// 	"UnpublishFirmware":             	ureq[v201.UnpublishFirmwareReq],
// 	"UpdateFirmware":                	ureq[v201.UpdateFirmwareReq],
// }


// var resmapv201 = map[string]func(*json.RawMessage) (Payload, error){
// 	"Authorize":                     	uconf[v201.AuthorizeRes],
// 	"BootNotification":              	uconf[v201.BootNotificationRes],
// 	"CancelReservation":             	uconf[v201.CancelReservationRes],
// 	"CertificateSigned":			 	uconf[v201.CertificateSignedRes],
// 	"ChangeAvailability":            	uconf[v201.ChangeAvailabilityRes],
// 	"ClearCache":                    	uconf[v201.ClearCacheRes],
// 	"ClearChargingProfile":          	uconf[v201.ClearChargingProfileRes],
// 	"ClearDisplayMessage":           	uconf[v201.ClearDisplayMessageRes],
// 	"ClearedChargingLimit":          	uconf[v201.ClearedChargingLimitRes],
// 	"ClearVariableMonitoring":       	uconf[v201.ClearVariableMonitoringRes],
// 	"CostUpdated":                   	uconf[v201.CostUpdatedRes],
// 	"CustomerInformation":           	uconf[v201.CustomerInformationRes],
// 	"DataTransfer":                  	uconf[v201.DataTransferRes],
// 	"DeleteCertificate":             	uconf[v201.DeleteCertificateRes],
// 	"FirmwareStatusNotification":    	uconf[v201.FirmwareStatusNotificationRes],
// 	"Get15118EVCertificate":         	uconf[v201.Get15118EVCertificateRes],
// 	"GetBaseReport":                 	uconf[v201.GetBaseReportRes],
// 	"GetCertificateStatus":          	uconf[v201.GetCertificateStatusRes],
// 	"GetChargingProfiles":            	uconf[v201.GetChargingProfilesRes],
// 	"GetCompositeSchedule":          	uconf[v201.GetCompositeScheduleRes],
// 	"GetDisplayMessages":            	uconf[v201.GetDisplayMessagesRes],
// 	"GetInstalledCertificateIds":      	uconf[v201.GetInstalledCertificateIdsRes],
// 	"GetLocalListVersion":           	uconf[v201.GetLocalListVersionRes],
// 	"GetLog":							uconf[v201.GetLogRes],
// 	"GetMonitoringReport":           	uconf[v201.GetMonitoringReportRes],
// 	"GetReport":                     	uconf[v201.GetReportRes],
// 	"GetTransactionStatus":          	uconf[v201.GetTransactionStatusRes],
// 	"GetVariables":                  	uconf[v201.GetVariablesRes],
// 	"Heartbeat":                     	uconf[v201.HeartbeatRes],
// 	"InstallCertificate":            	uconf[v201.InstallCertificateRes],
// 	"LogStatusNotification":         	uconf[v201.LogStatusNotificationRes],
// 	"MeterValues":                   	uconf[v201.MeterValuesRes],
// 	"NotifyChargingLimit":           	uconf[v201.NotifyChargingLimitRes],
// 	"NotifyCustomerInformation":     	uconf[v201.NotifyCustomerInformationRes],
// 	"NotifyDisplayMessages":         	uconf[v201.NotifyDisplayMessagesRes],
// 	"NotifyEVChargingNeeds":         	uconf[v201.NotifyEVChargingNeedsRes],
// 	"NotifyEVChargingSchedule":        	uconf[v201.NotifyEVChargingScheduleRes],
// 	"NotifyEvent":                   	uconf[v201.NotifyEventRes],
// 	"NotifyMonitoringReport":        	uconf[v201.NotifyMonitoringReportRes],
// 	"NotifyReport":                  	uconf[v201.NotifyReportRes],
// 	"PublishFirmware":               	uconf[v201.PublishFirmwareRes],
// 	"PublishFirmawareStatusNotification": uconf[v201.PublishFirmwareStatusNotificationRes],
// 	"ReportChargingProfiles":        	uconf[v201.ReportChargingProfilesRes],
// 	"RequestStartTransaction":       	uconf[v201.RequestStartTransactionRes],
// 	"RequestStopTransaction":        	uconf[v201.RequestStopTransactionRes],
// 	"ReservationStatusUpdate":       	uconf[v201.ReservationStatusUpdateRes],
// 	"ReserveNow":                    	uconf[v201.ReserveNowRes],
// 	"Reset":                         	uconf[v201.ResetRes],
// 	"SecurityEventNotification":     	uconf[v201.SecurityEventNotificationRes],
// 	"SendLocalList":                 	uconf[v201.SendLocalListRes],
// 	"SetChargingProfile":            	uconf[v201.SetChargingProfileRes],
// 	"SetDisplayMessage":             	uconf[v201.SetDisplayMessageRes],
// 	"SetMonitoringBase":             	uconf[v201.SetMonitoringBaseRes],
// 	"SetMonitoringLevel":            	uconf[v201.SetMonitoringLevelRes],
// 	"SetNetworkProfile":             	uconf[v201.SetNetworkProfileRes],
// 	"SetVariableMonitoring":         	uconf[v201.SetVariableMonitoringRes],
// 	"SetVariables":                  	uconf[v201.SetVariablesRes],
// 	"SignCertificate":               	uconf[v201.SignCertificateRes],
// 	"StatusNotification":            	uconf[v201.StatusNotificationRes],
// 	"TransactionEvent":			  		uconf[v201.TransactionEventRes],
// 	"TriggerMessage":                	uconf[v201.TriggerMessageRes],
// 	"UnlockConnector":               	uconf[v201.UnlockConnectorRes],
// 	"UnpublishFirmware":             	uconf[v201.UnpublishFirmwareRes],
// 	"UpdateFirmware":                	uconf[v201.UpdateFirmwareRes],
// }


type OCPPError struct {
	id    string
	code  string
	cause string
}

func (e *OCPPError) Error() string {
	return e.code + ": " + e.cause
}

// Call represents OCPP Call
type Call struct {
	MessageTypeId uint8
	UniqueId      string
	Action        string
	Payload       Payload
}

// Create CallResult from a received Call
func (call *Call) createCallResult(r Payload) *[]byte {
	out := [3]interface{}{
		3,
		call.UniqueId,
		r,
	}
	raw, _ := json.Marshal(out)
	return &raw
}

// Creates a CallError from a received Call
// TODO: organize error codes
func (call *Call) createCallError(err error) *[]byte {
	var id string
	var code string
	var cause string
	var ocppErr *OCPPError
	if errors.As(err, &ocppErr) {
		id = ocppErr.id
		code = ocppErr.code
		cause = ocppErr.cause
	}
	if id == "" {
		id = uuid.New().String()
	}
	callError := &CallError{
		UniqueId:         id,
		ErrorCode:        code,
		ErrorDescription: "",
		ErrorDetails:     cause,
	}
	switch code {
	case "ProtocolError":
		callError.ErrorDescription = "Payload for Action is incomplete"
	case "PropertyConstraintViolation":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	case "NotImplemented":
		callError.ErrorDescription = "Requested Action is not known by receiver"
	case "TypeConstraintViolationError":
		callError.ErrorDescription = "Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12)"
	case "PropertyConstraintViolationError":
		callError.ErrorDescription = "Payload is syntactically correct but at least one field contains an invalid value"
	default:
		callError.ErrorDescription = "Unknown error"
	}
	return callError.marshal()
}

// CallResult represents OCPP CallResult
type CallResult struct {
	MessageTypeId uint8
	UniqueId      string
	Payload       *json.RawMessage
}

// CallError represents OCPP CallError
type CallError struct {
	MessageTypeId    uint8
	UniqueId         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     interface{}
}

func (ce *CallError) marshal() *[]byte {
	ed := ce.ErrorDetails.(string)
	out := [5]interface{}{
		4,
		ce.UniqueId,
		ce.ErrorCode,
		ce.ErrorDescription,
		`{"cause":` + ed + `}`,
	}
	raw, _ := json.Marshal(out)
	return &raw
}

func (ce *CallError) Error() string {
	return fmt.Sprintf("CallError: UniqueId=%s, ErrorCode=%s, ErrorDescription=%s, ErrorDetails=%s",
		ce.UniqueId, ce.ErrorCode, ce.ErrorDescription, ce.ErrorDetails)
}

type TimeoutError struct {
	Message string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("TimeoutError: %s", e.Message)
}

// Converts raw byte to one of the ocpp messages or an error if the message is not valid
// [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}]
func unpack(b *[]byte) (*Call, *CallResult, *CallError, error) {
	var rm []json.RawMessage   	//  raw message
	var mti uint8 				//  MessageTypeId
	var ui string 				//  UniqueId
	var a string  				//  Action
	var p Payload 				//  Payload
	var c *Call
	var cr *CallResult
	var ce *CallError
	err := json.Unmarshal(*b, &rm)
	if err != nil {
		e := &OCPPError{
			id:    "",
			code:  "ProtocolError",
			cause: "Invalid JSON format",
		}
		log.Println(err)
		return nil, nil, nil, e
	}
	err = json.Unmarshal(rm[1], &ui)
	if err != nil {
		e := &OCPPError{
			id:    ui,
			code:  "ProtocolError",
			cause: "Message does not contain UniqueId",
		}
		return nil, nil, nil, e
	}
	if 3 > len(rm) || len(rm) > 5 {
		e := &OCPPError{
			id:    ui,
			code:  "ProtocolError",
			cause: "JSON must be an array of range [3,5]",
		}
		log.Println(err)
		return nil, nil, nil, e
	}
	err = json.Unmarshal(rm[0], &mti)
	if err != nil {
		e := &OCPPError{
			id:    ui,
			code:  "PropertyConstraintViolation",
			cause: fmt.Sprintf("MessageTypeId: %v is not valid", rm[0]),
		}
		return nil, nil, nil, e
	}
	if 2 > mti || mti > 4 {
		e := &OCPPError{
			id:    ui,
			code:  "ProtocolError",
			cause: "Message does not contain MessageTypeId",
		}
		return nil, nil, nil, e
	}
	if mti == 2 {
		err = json.Unmarshal(rm[2], &a)
		if err != nil {
			e := &OCPPError{
				id:    ui,
				code:  "ProtocolError",
				cause: "Message does not contain Action",
			}
			return nil, nil, nil, e
		}
		// print the rm
		// fmt.Println(rm)
		p, err = unmarshalReq(a, &rm[3])
		var ocppErr *OCPPError
		if err != nil {
			if errors.As(err, &ocppErr) {
				ocppErr.id = ui
			}
			return nil, nil, nil, err
		}
		c = &Call{
			MessageTypeId: mti,
			UniqueId:      ui,
			Action:        a,
			Payload:       p,
		}

	}
	if mti == 3 {
		p := &rm[2]
		cr = &CallResult{
			MessageTypeId: mti,
			UniqueId:      ui,
			Payload:       p,
		}
	}
	if mti == 4 {
		var me [5]interface{}
		_ = json.Unmarshal(*b, &me)
		ce = &CallError{
			MessageTypeId:    mti,
			UniqueId:         ui,
			ErrorCode:        me[2].(string),
			ErrorDescription: me[3].(string),
			ErrorDetails:     me[4],
		}
	}
	return c, cr, ce, nil

}



// Converts raw CP initiated Call Payload to a corresponding Payload struct
func unmarshalReq(mAction string, rawPayload *json.RawMessage) (Payload, error) {
	a, ok := reqmap[mAction]
	if !ok {
		e := &OCPPError{
			code:  "NotImplemented",
			cause: fmt.Sprintf("Action %v is not implemented", mAction),
		}
		return nil, e
	}
	return a(rawPayload)
}

// Unmarshal Payload to a struct of type T, e.g. BootNotificationReq
func ureq[T any](rawPayload *json.RawMessage) (Payload, error) {
	var p *T
	var payload Payload
	err := json.Unmarshal(*rawPayload, &p)
	if err != nil {
		e := &OCPPError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Println(err)
		return nil, e
	}
	err = validate.Struct(*p)
	if err != nil {
		// TODO: construct more detailed error
		e := &OCPPError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Println(err)
		return nil, e
	}
	payload = p
	return payload, nil
}



func unmarshalConf(mAction string, rawPayload *json.RawMessage) (Payload, error) {
	a, ok := confmap[mAction]
	if !ok {
		err := errors.New("invalid action")
		return nil, err
	}
	return a(rawPayload)
}

// Unmarshal Raw Payload to a struct of type T, e.g. ChangeAvailabilityConf
func uconf[T any](rawPayload *json.RawMessage) (Payload, error) {
	var p *T
	var payload Payload
	err := json.Unmarshal(*rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(*p)
	if err != nil {
		return nil, err
	}
	payload = p
	return payload, nil
}
