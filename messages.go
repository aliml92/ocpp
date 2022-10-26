package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aliml92/ocpp/v16"
	"github.com/aliml92/ocpp/v201"
	"github.com/google/uuid"
)

var reqmapv16, resmapv16, reqmapv201, resmapv201   map[string]func(json.RawMessage) (Payload, error)

func init(){
	reqmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"BootNotification":              ureqV16[v16.BootNotificationReq],
		"Authorize":                     ureqV16[v16.AuthorizeReq],
		"DataTransfer":                  ureqV16[v16.DataTransferReq],
		"DiagnosticsStatusNotification": ureqV16[v16.DiagnosticsStatusNotificationReq],
		"FirmwareStatusNotification":    ureqV16[v16.FirmwareStatusNotificationReq],
		"Heartbeat":                     ureqV16[v16.HeartbeatReq],
		"MeterValues":                   ureqV16[v16.MeterValuesReq],
		"StartTransaction":              ureqV16[v16.StartTransactionReq],
		"StatusNotification":            ureqV16[v16.StatusNotificationReq],
		"StopTransaction":               ureqV16[v16.StopTransactionReq],
		"CanCelReservation":             ureqV16[v16.CancelReservationReq],
		"ChangeAvailability":            ureqV16[v16.ChangeAvailabilityReq],
		"ChangeConfiguration":           ureqV16[v16.ChangeConfigurationReq],
		"ClearCache":                    ureqV16[v16.ClearCacheReq],
		"ClearChargingProfile":          ureqV16[v16.ClearChargingProfileReq],
		"GetCompositeSchedule":          ureqV16[v16.GetCompositeScheduleReq],
		"GetConfiguration":              ureqV16[v16.GetConfigurationReq],
		"GetDiagnostics":                ureqV16[v16.GetDiagnosticsReq],
		"GetLocalListVersion":           ureqV16[v16.GetLocalListVersionReq],
		"RemoteStartTransaction":        ureqV16[v16.RemoteStartTransactionReq],
		"RemoteStopTransaction":         ureqV16[v16.RemoteStopTransactionReq],
		"ReserveNow":                    ureqV16[v16.ReserveNowReq],
		"Reset":                         ureqV16[v16.ResetReq],
		"SendLocalList":                 ureqV16[v16.SendLocalListReq],
		"SetChargingProfile":            ureqV16[v16.SetChargingProfileReq],
		"TriggerMessage":                ureqV16[v16.TriggerMessageReq],
		"UnlockConnector":               ureqV16[v16.UnlockConnectorReq],
		"UpdateFirmware":                ureqV16[v16.UpdateFirmwareReq],
	}

	resmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"BootNotification":              uresV16[v16.BootNotificationConf],
		"Authorize":                     uresV16[v16.AuthorizeConf],
		"DataTransfer":                  uresV16[v16.DataTransferConf],
		"DiagnosticsStatusNotification": uresV16[v16.DiagnosticsStatusNotificationConf],
		"FirmwareStatusNotification":    uresV16[v16.FirmwareStatusNotificationConf],
		"Heartbeat":                     uresV16[v16.HeartbeatConf],
		"MeterValues":                   uresV16[v16.MeterValuesConf],
		"StartTransaction":              uresV16[v16.StartTransactionConf],
		"StatusNotification":            uresV16[v16.StatusNotificationConf],
		"StopTransaction":               uresV16[v16.StopTransactionConf],
		"CancelReservation":             uresV16[v16.CancelReservationConf],
		"ChangeAvailability":            uresV16[v16.ChangeAvailabilityConf],
		"ChangeConfiguration":           uresV16[v16.ChangeConfigurationConf],
		"ClearCache":                    uresV16[v16.ClearCacheConf],
		"ClearChargingProfile":          uresV16[v16.ClearChargingProfileConf],
		"GetCompositeSchedule":          uresV16[v16.GetCompositeScheduleConf],
		"GetConfiguration":              uresV16[v16.GetConfigurationConf],
		"GetDiagnostics":                uresV16[v16.GetDiagnosticsConf],
		"GetLocalListVersion":           uresV16[v16.GetLocalListVersionConf],
		"RemoteStartTransaction":        uresV16[v16.RemoteStartTransactionConf],
		"RemoteStopTransaction":         uresV16[v16.RemoteStopTransactionConf],
		"ReserveNow":                    uresV16[v16.ReserveNowConf],
		"Reset":                         uresV16[v16.ResetConf],
		"SendLocalList":                 uresV16[v16.SendLocalListConf],
		"SetChargingProfile":            uresV16[v16.SetChargingProfileConf],
		"TriggerMessage":                uresV16[v16.TriggerMessageConf],
		"UnlockConnector":               uresV16[v16.UnlockConnectorConf],
		"UpdateFirmware":                uresV16[v16.UpdateFirmwareConf],
	}

	reqmapv201 = map[string]func(json.RawMessage) (Payload, error){
		"Authorize":                     	ureqV201[v201.AuthorizeReq],
		"BootNotification":              	ureqV201[v201.BootNotificationReq],
		"CancelReservation":             	ureqV201[v201.CancelReservationReq],
		"CertificateSigned":			 	ureqV201[v201.CertificateSignedReq],
		"ChangeAvailability":            	ureqV201[v201.ChangeAvailabilityReq],
		"ClearCache":                    	ureqV201[v201.ClearCacheReq],
		"ClearChargingProfile":          	ureqV201[v201.ClearChargingProfileReq],
		"ClearDisplayMessage":           	ureqV201[v201.ClearDisplayMessageReq],
		"ClearedChargingLimit":          	ureqV201[v201.ClearedChargingLimitReq],
		"ClearVariableMonitoring":       	ureqV201[v201.ClearVariableMonitoringReq],
		"CostUpdated":                   	ureqV201[v201.CostUpdatedReq],
		"CustomerInformation":           	ureqV201[v201.CustomerInformationReq],
		"DataTransfer":                  	ureqV201[v201.DataTransferReq],
		"DeleteCertificate":             	ureqV201[v201.DeleteCertificateReq],
		"FirmwareStatusNotification":    	ureqV201[v201.FirmwareStatusNotificationReq],
		"Get15118EVCertificate":         	ureqV201[v201.Get15118EVCertificateReq],
		"GetBaseReport":                 	ureqV201[v201.GetBaseReportReq],
		"GetCertificateStatus":          	ureqV201[v201.GetCertificateStatusReq],
		"GetChargingProfiles":            	ureqV201[v201.GetChargingProfilesReq],
		"GetCompositeSchedule":          	ureqV201[v201.GetCompositeScheduleReq],
		"GetDisplayMessages":            	ureqV201[v201.GetDisplayMessagesReq],
		"GetInstalledCertificateIds":      	ureqV201[v201.GetInstalledCertificateIdsReq],
		"GetLocalListVersion":           	ureqV201[v201.GetLocalListVersionReq],
		"GetLog":							ureqV201[v201.GetLogReq],
		"GetMonitoringReport":           	ureqV201[v201.GetMonitoringReportReq],
		"GetReport":                     	ureqV201[v201.GetReportReq],
		"GetTransactionStatus":          	ureqV201[v201.GetTransactionStatusReq],
		"GetVariables":                  	ureqV201[v201.GetVariablesReq],
		"Heartbeat":                     	ureqV201[v201.HeartbeatReq],
		"InstallCertificate":            	ureqV201[v201.InstallCertificateReq],
		"LogStatusNotification":         	ureqV201[v201.LogStatusNotificationReq],
		"MeterValues":                   	ureqV201[v201.MeterValuesReq],
		"NotifyChargingLimit":           	ureqV201[v201.NotifyChargingLimitReq],
		"NotifyCustomerInformation":     	ureqV201[v201.NotifyCustomerInformationReq],
		"NotifyDisplayMessages":         	ureqV201[v201.NotifyDisplayMessagesReq],
		"NotifyEVChargingNeeds":         	ureqV201[v201.NotifyEVChargingNeedsReq],
		"NotifyEVChargingSchedule":        	ureqV201[v201.NotifyEVChargingScheduleReq],
		"NotifyEvent":                   	ureqV201[v201.NotifyEventReq],
		"NotifyMonitoringReport":        	ureqV201[v201.NotifyMonitoringReportReq],
		"NotifyReport":                  	ureqV201[v201.NotifyReportReq],
		"PublishFirmware":               	ureqV201[v201.PublishFirmwareReq],
		"PublishFirmawareStatusNotification": ureqV201[v201.PublishFirmwareStatusNotificationReq],
		"ReportChargingProfiles":        	ureqV201[v201.ReportChargingProfilesReq],
		"RequestStartTransaction":       	ureqV201[v201.RequestStartTransactionReq],
		"RequestStopTransaction":        	ureqV201[v201.RequestStopTransactionReq],
		"ReservationStatusUpdate":       	ureqV201[v201.ReservationStatusUpdateReq],
		"ReserveNow":                    	ureqV201[v201.ReserveNowReq],
		"Reset":                         	ureqV201[v201.ResetReq],
		"SecurityEventNotification":     	ureqV201[v201.SecurityEventNotificationReq],
		"SendLocalList":                 	ureqV201[v201.SendLocalListReq],
		"SetChargingProfile":            	ureqV201[v201.SetChargingProfileReq],
		"SetDisplayMessage":             	ureqV201[v201.SetDisplayMessageReq],
		"SetMonitoringBase":             	ureqV201[v201.SetMonitoringBaseReq],
		"SetMonitoringLevel":            	ureqV201[v201.SetMonitoringLevelReq],
		"SetNetworkProfile":             	ureqV201[v201.SetNetworkProfileReq],
		"SetVariableMonitoring":         	ureqV201[v201.SetVariableMonitoringReq],
		"SetVariables":                  	ureqV201[v201.SetVariablesReq],
		"SignCertificate":               	ureqV201[v201.SignCertificateReq],
		"StatusNotification":            	ureqV201[v201.StatusNotificationReq],
		"TransactionEvent":			  		ureqV201[v201.TransactionEventReq],
		"TriggerMessage":                	ureqV201[v201.TriggerMessageReq],
		"UnlockConnector":               	ureqV201[v201.UnlockConnectorReq],
		"UnpublishFirmware":             	ureqV201[v201.UnpublishFirmwareReq],
		"UpdateFirmware":                	ureqV201[v201.UpdateFirmwareReq],
	}

	resmapv201 = map[string]func(json.RawMessage) (Payload, error){
		"Authorize":                     	uresV201[v201.AuthorizeRes],
		"BootNotification":              	uresV201[v201.BootNotificationRes],
		"CancelReservation":             	uresV201[v201.CancelReservationRes],
		"CertificateSigned":			 	uresV201[v201.CertificateSignedRes],
		"ChangeAvailability":            	uresV201[v201.ChangeAvailabilityRes],
		"ClearCache":                    	uresV201[v201.ClearCacheRes],
		"ClearChargingProfile":          	uresV201[v201.ClearChargingProfileRes],
		"ClearDisplayMessage":           	uresV201[v201.ClearDisplayMessageRes],
		"ClearedChargingLimit":          	uresV201[v201.ClearedChargingLimitRes],
		"ClearVariableMonitoring":       	uresV201[v201.ClearVariableMonitoringRes],
		"CostUpdated":                   	uresV201[v201.CostUpdatedRes],
		"CustomerInformation":           	uresV201[v201.CustomerInformationRes],
		"DataTransfer":                  	uresV201[v201.DataTransferRes],
		"DeleteCertificate":             	uresV201[v201.DeleteCertificateRes],
		"FirmwareStatusNotification":    	uresV201[v201.FirmwareStatusNotificationRes],
		"Get15118EVCertificate":         	uresV201[v201.Get15118EVCertificateRes],
		"GetBaseReport":                 	uresV201[v201.GetBaseReportRes],
		"GetCertificateStatus":          	uresV201[v201.GetCertificateStatusRes],
		"GetChargingProfiles":            	uresV201[v201.GetChargingProfilesRes],
		"GetCompositeSchedule":          	uresV201[v201.GetCompositeScheduleRes],
		"GetDisplayMessages":            	uresV201[v201.GetDisplayMessagesRes],
		"GetInstalledCertificateIds":      	uresV201[v201.GetInstalledCertificateIdsRes],
		"GetLocalListVersion":           	uresV201[v201.GetLocalListVersionRes],
		"GetLog":							uresV201[v201.GetLogRes],
		"GetMonitoringReport":           	uresV201[v201.GetMonitoringReportRes],
		"GetReport":                     	uresV201[v201.GetReportRes],
		"GetTransactionStatus":          	uresV201[v201.GetTransactionStatusRes],
		"GetVariables":                  	uresV201[v201.GetVariablesRes],
		"Heartbeat":                     	uresV201[v201.HeartbeatRes],
		"InstallCertificate":            	uresV201[v201.InstallCertificateRes],
		"LogStatusNotification":         	uresV201[v201.LogStatusNotificationRes],
		"MeterValues":                   	uresV201[v201.MeterValuesRes],
		"NotifyChargingLimit":           	uresV201[v201.NotifyChargingLimitRes],
		"NotifyCustomerInformation":     	uresV201[v201.NotifyCustomerInformationRes],
		"NotifyDisplayMessages":         	uresV201[v201.NotifyDisplayMessagesRes],
		"NotifyEVChargingNeeds":         	uresV201[v201.NotifyEVChargingNeedsRes],
		"NotifyEVChargingSchedule":        	uresV201[v201.NotifyEVChargingScheduleRes],
		"NotifyEvent":                   	uresV201[v201.NotifyEventRes],
		"NotifyMonitoringReport":        	uresV201[v201.NotifyMonitoringReportRes],
		"NotifyReport":                  	uresV201[v201.NotifyReportRes],
		"PublishFirmware":               	uresV201[v201.PublishFirmwareRes],
		"PublishFirmawareStatusNotification": uresV201[v201.PublishFirmwareStatusNotificationRes],
		"ReportChargingProfiles":        	uresV201[v201.ReportChargingProfilesRes],
		"RequestStartTransaction":       	uresV201[v201.RequestStartTransactionRes],
		"RequestStopTransaction":        	uresV201[v201.RequestStopTransactionRes],
		"ReservationStatusUpdate":       	uresV201[v201.ReservationStatusUpdateRes],
		"ReserveNow":                    	uresV201[v201.ReserveNowRes],
		"Reset":                         	uresV201[v201.ResetRes],
		"SecurityEventNotification":     	uresV201[v201.SecurityEventNotificationRes],
		"SendLocalList":                 	uresV201[v201.SendLocalListRes],
		"SetChargingProfile":            	uresV201[v201.SetChargingProfileRes],
		"SetDisplayMessage":             	uresV201[v201.SetDisplayMessageRes],
		"SetMonitoringBase":             	uresV201[v201.SetMonitoringBaseRes],
		"SetMonitoringLevel":            	uresV201[v201.SetMonitoringLevelRes],
		"SetNetworkProfile":             	uresV201[v201.SetNetworkProfileRes],
		"SetVariableMonitoring":         	uresV201[v201.SetVariableMonitoringRes],
		"SetVariables":                  	uresV201[v201.SetVariablesRes],
		"SignCertificate":               	uresV201[v201.SignCertificateRes],
		"StatusNotification":            	uresV201[v201.StatusNotificationRes],
		"TransactionEvent":			  		uresV201[v201.TransactionEventRes],
		"TriggerMessage":                	uresV201[v201.TriggerMessageRes],
		"UnlockConnector":               	uresV201[v201.UnlockConnectorRes],
		"UnpublishFirmware":             	uresV201[v201.UnpublishFirmwareRes],
		"UpdateFirmware":                	uresV201[v201.UpdateFirmwareRes],
	}
}




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
func (call *Call) createCallResult(r Payload) []byte {
	out := [3]interface{}{
		3,
		call.UniqueId,
		r,
	}
	raw, _ := json.Marshal(out)
	return raw
}

// Creates a CallError from a received Call
// TODO: organize error codes
func (call *Call) createCallError(err error) []byte {
	var id, code, cause string
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
	Payload       json.RawMessage
}

// CallError represents OCPP CallError
type CallError struct {
	MessageTypeId    uint8
	UniqueId         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     interface{}
}

func (ce *CallError) marshal() []byte {
	ed := ce.ErrorDetails.(string)
	out := [5]interface{}{
		4,
		ce.UniqueId,
		ce.ErrorCode,
		ce.ErrorDescription,
		`{"cause":` + ed + `}`,
	}
	raw, _ := json.Marshal(out)
	return raw
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
func unpack(b []byte, proto string) (*Call, *CallResult, *CallError, error) {
	var rm []json.RawMessage //  raw message
	var mti uint8            //  MessageTypeId
	var ui string            //  UniqueId
	var a string             //  Action
	var p Payload            //  Payload
	var c *Call
	var cr *CallResult
	var ce *CallError
	err := json.Unmarshal(b, &rm)
	if err != nil {
		e := &OCPPError{
			id:    "",
			code:  "ProtocolError",
			cause: "Invalid JSON format",
		}

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
		log.Error(err)
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
		p, err = unmarshalReq(a, rm[3], proto)
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
		p := rm[2]
		cr = &CallResult{
			MessageTypeId: mti,
			UniqueId:      ui,
			Payload:       p,
		}
	}
	if mti == 4 {
		var me [5]interface{}
		_ = json.Unmarshal(b, &me)
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
func unmarshalReq(mAction string, rawPayload json.RawMessage, proto string) (Payload, error) {
	var a func(json.RawMessage) (Payload, error)
	var ok bool
	switch proto {
	case ocppV16:
		a, ok = reqmapv16[mAction]
	case ocppV201:
		a, ok = reqmapv201[mAction]	
	}
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
func ureqV16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p *T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		e := &OCPPError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Error(err)
		return nil, e
	}
	err = validateV16.Struct(*p)
	if err != nil {
		// TODO: construct more detailed error
		e := &OCPPError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Error(err)
		return nil, e
	}
	payload = p
	return payload, nil
}


// Unmarshal Payload to a struct of type T, e.g. BootNotificationReq
func ureqV201[T any](rawPayload json.RawMessage) (Payload, error) {
	var p *T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		e := &OCPPError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Error(err)
		return nil, e
	}
	err = validateV201.Struct(*p)
	if err != nil {
		// TODO: construct more detailed error
		e := &OCPPError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Error(err)
		return nil, e
	}
	payload = p
	return payload, nil
}


func unmarshalResV16(mAction string, rawPayload json.RawMessage) (Payload, error) {
	a, ok := resmapv16[mAction]
	if !ok {
		err := errors.New("invalid action")
		return nil, err
	}
	return a(rawPayload)
}


func unmarshalResV201(mAction string, rawPayload json.RawMessage) (Payload, error) {
	a, ok := resmapv201[mAction]
	if !ok {
		err := errors.New("invalid action")
		return nil, err
	}
	return a(rawPayload)
}

// Unmarshal Raw Payload to a struct of type T, e.g. ChangeAvailabilityConf
func uresV16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p *T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = validateV16.Struct(*p)
	if err != nil {
		return nil, err
	}
	payload = p
	return payload, nil
}


// Unmarshal Raw Payload to a struct of type T, e.g. ChangeAvailabilityConf
func uresV201[T any](rawPayload json.RawMessage) (Payload, error) {
	var p *T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = validateV201.Struct(*p)
	if err != nil {
		return nil, err
	}
	payload = p
	return payload, nil
}