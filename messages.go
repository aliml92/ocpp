package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aliml92/ocpp/v16"
	"github.com/aliml92/ocpp/v201"
)

const (
	MessageTypeIdCall 		=  2
	MessageTypeIdCallResult =  3
	MessageTypeIdCallError  =  4
	
)

var errInvalidAction =   errors.New("invalid action") 

var reqmapv16, resmapv16, reqmapv201, resmapv201   map[string]func(json.RawMessage) (Payload, error)

type ocppError struct {
	id    string
	code  string
	cause string
}

func (e *ocppError) Error() string {
	return e.code + ": " + e.cause
}

// Call represents OCPP Call
type Call struct {
	MessageTypeId uint8
	UniqueId      string
	Action        string
	Payload       Payload
}


func (c *Call) getID() string {
	return c.UniqueId
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
	var ocppErr *ocppError
	if errors.As(err, &ocppErr) {
		id = ocppErr.id
		code = ocppErr.code
		cause = ocppErr.cause
	}
	if id == "" {
		id = "-1"
	}
	callError := &CallError{
		UniqueId:         id,
		ErrorCode:        code,
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
	case "MessageTypeNotSupported":
		callError.ErrorDescription = "A message with an Message Type Number received that is not supported by this implementation"
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
 
func (cr *CallResult) getID() string {
	return cr.UniqueId
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

func (ce *CallError) getID() string {
	return ce.UniqueId
}

type OcppMessage interface {
	getID() string
}



// umpack converts json byte to one of the ocpp messages, if not successful 
// returns an error
// unpack expects ocpp messages in the below forms:
//    -   [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}] -> Call
//    -   [<MessageTypeId>, "<UniqueId>", {<Payload>}]             -> CallResult
//    -   [<MessageTypeId>, "<UniqueId>", "<ErrorCode>", "<ErrorDescription>" , {<ErrorDetails>}] -> CallError
// if json byte does not conform with these three formats, OcppError is created
// and is returned as error
//
// TODO: improve error construction depending on protocol
func unpack(b []byte, proto string) (OcppMessage, error) {
	var rm []json.RawMessage //  raw message
	var mti uint8            //  MessageTypeId
	var ui string            //  UniqueId
	var a string             //  Action
	var p Payload            //  Payload
	var ocppMsg OcppMessage  //  OcppMessage
	var e *ocppError

	// unmarshal []byte to []json.RawMessage
	err := json.Unmarshal(b, &rm)
	if err != nil {
		e = &ocppError{
			id:    "-1",
			code:  "ProtocolError",
			cause: "Invalid JSON format",
		}
		return nil, e
	}


	// unmarshal [0]json.RawMessage to MessageTypeId
	err1 := json.Unmarshal(rm[0], &mti)
	// unmarshal [1]json.RawMessage to UniqueId
	err2 := json.Unmarshal(rm[1], &ui)
	if err1 != nil {	
		if err2 != nil {
			e.id = "-1"
			e.cause = e.cause + "," + fmt.Sprintf("UniqueId: %v is not valid", rm[1])
		} else {
			e.id = ui
		}
		return nil, e	
	}
	if len(ui) > 36 {
		e = &ocppError{
			code:  "ProtocolError",
			cause: fmt.Sprintf("UniqueId: %v is too long", ui),
		}
	}
	switch mti {
	case MessageTypeIdCall:
		call := &Call{
			MessageTypeId: mti,
			UniqueId: ui,
		}
		if e != nil {
			return call, e
		}
		err = json.Unmarshal(rm[2], &a)
		if err != nil {
			e := &ocppError{
				id:    ui,
				code:  "ProtocolError",
				cause: "Message does not contain Action",
			}
			return call, e
		}
		call.Action = a
		p, err = unmarshalRequestPayload(a, rm[3], proto)
		var ocppErr *ocppError
		if err != nil {
			if errors.As(err, &ocppErr) {
				ocppErr.id = ui
			}
			return call, err
		}
		call.Payload = p
		ocppMsg = call
	case MessageTypeIdCallResult:
		p := rm[2]
		ocppMsg = &CallResult{
			MessageTypeId: mti,
			UniqueId:      ui,
			Payload:       p,
		}
	case MessageTypeIdCallError:
		var me [5]interface{}
		_ = json.Unmarshal(b, &me)
		ocppMsg = &CallError{
			MessageTypeId:    mti,
			UniqueId:         ui,
			ErrorCode:        me[2].(string),
			ErrorDescription: me[3].(string),
			ErrorDetails:     me[4],
		}
	default:
		e := &ocppError{
			id:    ui,
			code:  "MessageTypeNotSupported",
			cause: fmt.Sprintf("A message with: %v is not supported by this implementation", mti),
		}
		return nil, e									
	}
	return ocppMsg, nil
}

// unmarshalRequestPayload unmarshals raw bytes of request type payload 
// to a corresponding struct depending on Action and ocpp protocol
func unmarshalRequestPayload(actionName string, rawPayload json.RawMessage, proto string) (Payload, error) {
	var uf func(json.RawMessage) (Payload, error)  // uf unmarshal function for a specific action request
	var ok bool
	switch proto {
	case ocppV16:
		uf, ok = reqmapv16[actionName]
	case ocppV201:
		uf, ok = reqmapv201[actionName]	
	}
	if !ok {
		e := &ocppError{
			code:  "NotImplemented",
			cause: fmt.Sprintf("Action %v is not implemented", actionName),
		}
		return nil, e
	}
	return uf(rawPayload)
}



// unmarshalRequestPayloadv16 unmarshals raw request type payload to a ***Req type struct of ocppv16 
func unmarshalRequestPayloadv16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		e := &ocppError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		return nil, e
	}
	err = validateV16.Struct(p)
	if err != nil {
		// TODO: construct more detailed error
		e := &ocppError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		return nil, e
	}
	payload = &p
	return payload, nil
}


// unmarshalRequestPayloadv201 unmarshals raw request type payload to a ***Req type struct of ocppv201 
func unmarshalRequestPayloadv201[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		e := &ocppError{
			code:  "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Error(err)
		return nil, e
	}
	err = validateV201.Struct(p)
	if err != nil {
		// TODO: construct more detailed error
		e := &ocppError{
			code:  "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Error(err)
		return nil, e
	}
	payload = &p
	return payload, nil
}


// unmarshalResponsePv16 unmarshals raw bytes of request type payload 
// to a corresponding struct depending on actionName 
func unmarshalResponsePv16(actionName string, rawPayload json.RawMessage) (Payload, error) {
	uf, ok := resmapv16[actionName]       // uf unmarshal function for a specific action request 
	if !ok {
		return nil, errInvalidAction
	}
	return uf(rawPayload)
}

// unmarshalResponsePv201 unmarshals raw bytes of request type payload 
// to a corresponding struct depending on actionName
func unmarshalResponsePv201(mAction string, rawPayload json.RawMessage) (Payload, error) {
	uf, ok := resmapv201[mAction]          // uf unmarshal function for a specific action request 
	if !ok {
		return nil, errInvalidAction
	}
	return uf(rawPayload)
}

// unmarshalResponsePayloadv16 unmarshals raw response type payload to a ***Conf type struct of ocppv16 
func unmarshalResponsePayloadv16[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = validateV16.Struct(p)
	if err != nil {
		return nil, err
	}
	payload = &p
	return payload, nil
}


// unmarshalResponsePayloadv201 unmarshals raw response type payload to a ***Res type struct of ocppv201 
func unmarshalResponsePayloadv201[T any](rawPayload json.RawMessage) (Payload, error) {
	var p T
	var payload Payload
	err := json.Unmarshal(rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = validateV201.Struct(p)
	if err != nil {
		return nil, err
	}
	payload = &p
	return payload, nil
}



func init(){
	reqmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"BootNotification":              unmarshalRequestPayloadv16[v16.BootNotificationReq],
		"Authorize":                     unmarshalRequestPayloadv16[v16.AuthorizeReq],
		"DataTransfer":                  unmarshalRequestPayloadv16[v16.DataTransferReq],
		"DiagnosticsStatusNotification": unmarshalRequestPayloadv16[v16.DiagnosticsStatusNotificationReq],
		"FirmwareStatusNotification":    unmarshalRequestPayloadv16[v16.FirmwareStatusNotificationReq],
		"Heartbeat":                     unmarshalRequestPayloadv16[v16.HeartbeatReq],
		"MeterValues":                   unmarshalRequestPayloadv16[v16.MeterValuesReq],
		"StartTransaction":              unmarshalRequestPayloadv16[v16.StartTransactionReq],
		"StatusNotification":            unmarshalRequestPayloadv16[v16.StatusNotificationReq],
		"StopTransaction":               unmarshalRequestPayloadv16[v16.StopTransactionReq],
		"CanCelReservation":             unmarshalRequestPayloadv16[v16.CancelReservationReq],
		"ChangeAvailability":            unmarshalRequestPayloadv16[v16.ChangeAvailabilityReq],
		"ChangeConfiguration":           unmarshalRequestPayloadv16[v16.ChangeConfigurationReq],
		"ClearCache":                    unmarshalRequestPayloadv16[v16.ClearCacheReq],
		"ClearChargingProfile":          unmarshalRequestPayloadv16[v16.ClearChargingProfileReq],
		"GetCompositeSchedule":          unmarshalRequestPayloadv16[v16.GetCompositeScheduleReq],
		"GetConfiguration":              unmarshalRequestPayloadv16[v16.GetConfigurationReq],
		"GetDiagnostics":                unmarshalRequestPayloadv16[v16.GetDiagnosticsReq],
		"GetLocalListVersion":           unmarshalRequestPayloadv16[v16.GetLocalListVersionReq],
		"RemoteStartTransaction":        unmarshalRequestPayloadv16[v16.RemoteStartTransactionReq],
		"RemoteStopTransaction":         unmarshalRequestPayloadv16[v16.RemoteStopTransactionReq],
		"ReserveNow":                    unmarshalRequestPayloadv16[v16.ReserveNowReq],
		"Reset":                         unmarshalRequestPayloadv16[v16.ResetReq],
		"SendLocalList":                 unmarshalRequestPayloadv16[v16.SendLocalListReq],
		"SetChargingProfile":            unmarshalRequestPayloadv16[v16.SetChargingProfileReq],
		"TriggerMessage":                unmarshalRequestPayloadv16[v16.TriggerMessageReq],
		"UnlockConnector":               unmarshalRequestPayloadv16[v16.UnlockConnectorReq],
		"UpdateFirmware":                unmarshalRequestPayloadv16[v16.UpdateFirmwareReq],
	}

	resmapv16 = map[string]func(json.RawMessage) (Payload, error){
		"BootNotification":              unmarshalResponsePayloadv16[v16.BootNotificationConf],
		"Authorize":                     unmarshalResponsePayloadv16[v16.AuthorizeConf],
		"DataTransfer":                  unmarshalResponsePayloadv16[v16.DataTransferConf],
		"DiagnosticsStatusNotification": unmarshalResponsePayloadv16[v16.DiagnosticsStatusNotificationConf],
		"FirmwareStatusNotification":    unmarshalResponsePayloadv16[v16.FirmwareStatusNotificationConf],
		"Heartbeat":                     unmarshalResponsePayloadv16[v16.HeartbeatConf],
		"MeterValues":                   unmarshalResponsePayloadv16[v16.MeterValuesConf],
		"StartTransaction":              unmarshalResponsePayloadv16[v16.StartTransactionConf],
		"StatusNotification":            unmarshalResponsePayloadv16[v16.StatusNotificationConf],
		"StopTransaction":               unmarshalResponsePayloadv16[v16.StopTransactionConf],
		"CancelReservation":             unmarshalResponsePayloadv16[v16.CancelReservationConf],
		"ChangeAvailability":            unmarshalResponsePayloadv16[v16.ChangeAvailabilityConf],
		"ChangeConfiguration":           unmarshalResponsePayloadv16[v16.ChangeConfigurationConf],
		"ClearCache":                    unmarshalResponsePayloadv16[v16.ClearCacheConf],
		"ClearChargingProfile":          unmarshalResponsePayloadv16[v16.ClearChargingProfileConf],
		"GetCompositeSchedule":          unmarshalResponsePayloadv16[v16.GetCompositeScheduleConf],
		"GetConfiguration":              unmarshalResponsePayloadv16[v16.GetConfigurationConf],
		"GetDiagnostics":                unmarshalResponsePayloadv16[v16.GetDiagnosticsConf],
		"GetLocalListVersion":           unmarshalResponsePayloadv16[v16.GetLocalListVersionConf],
		"RemoteStartTransaction":        unmarshalResponsePayloadv16[v16.RemoteStartTransactionConf],
		"RemoteStopTransaction":         unmarshalResponsePayloadv16[v16.RemoteStopTransactionConf],
		"ReserveNow":                    unmarshalResponsePayloadv16[v16.ReserveNowConf],
		"Reset":                         unmarshalResponsePayloadv16[v16.ResetConf],
		"SendLocalList":                 unmarshalResponsePayloadv16[v16.SendLocalListConf],
		"SetChargingProfile":            unmarshalResponsePayloadv16[v16.SetChargingProfileConf],
		"TriggerMessage":                unmarshalResponsePayloadv16[v16.TriggerMessageConf],
		"UnlockConnector":               unmarshalResponsePayloadv16[v16.UnlockConnectorConf],
		"UpdateFirmware":                unmarshalResponsePayloadv16[v16.UpdateFirmwareConf],
	}

	reqmapv201 = map[string]func(json.RawMessage) (Payload, error){
		"Authorize":                     	unmarshalRequestPayloadv201[v201.AuthorizeReq],
		"BootNotification":              	unmarshalRequestPayloadv201[v201.BootNotificationReq],
		"CancelReservation":             	unmarshalRequestPayloadv201[v201.CancelReservationReq],
		"CertificateSigned":			 	unmarshalRequestPayloadv201[v201.CertificateSignedReq],
		"ChangeAvailability":            	unmarshalRequestPayloadv201[v201.ChangeAvailabilityReq],
		"ClearCache":                    	unmarshalRequestPayloadv201[v201.ClearCacheReq],
		"ClearChargingProfile":          	unmarshalRequestPayloadv201[v201.ClearChargingProfileReq],
		"ClearDisplayMessage":           	unmarshalRequestPayloadv201[v201.ClearDisplayMessageReq],
		"ClearedChargingLimit":          	unmarshalRequestPayloadv201[v201.ClearedChargingLimitReq],
		"ClearVariableMonitoring":       	unmarshalRequestPayloadv201[v201.ClearVariableMonitoringReq],
		"CostUpdated":                   	unmarshalRequestPayloadv201[v201.CostUpdatedReq],
		"CustomerInformation":           	unmarshalRequestPayloadv201[v201.CustomerInformationReq],
		"DataTransfer":                  	unmarshalRequestPayloadv201[v201.DataTransferReq],
		"DeleteCertificate":             	unmarshalRequestPayloadv201[v201.DeleteCertificateReq],
		"FirmwareStatusNotification":    	unmarshalRequestPayloadv201[v201.FirmwareStatusNotificationReq],
		"Get15118EVCertificate":         	unmarshalRequestPayloadv201[v201.Get15118EVCertificateReq],
		"GetBaseReport":                 	unmarshalRequestPayloadv201[v201.GetBaseReportReq],
		"GetCertificateStatus":          	unmarshalRequestPayloadv201[v201.GetCertificateStatusReq],
		"GetChargingProfiles":            	unmarshalRequestPayloadv201[v201.GetChargingProfilesReq],
		"GetCompositeSchedule":          	unmarshalRequestPayloadv201[v201.GetCompositeScheduleReq],
		"GetDisplayMessages":            	unmarshalRequestPayloadv201[v201.GetDisplayMessagesReq],
		"GetInstalledCertificateIds":      	unmarshalRequestPayloadv201[v201.GetInstalledCertificateIdsReq],
		"GetLocalListVersion":           	unmarshalRequestPayloadv201[v201.GetLocalListVersionReq],
		"GetLog":							unmarshalRequestPayloadv201[v201.GetLogReq],
		"GetMonitoringReport":           	unmarshalRequestPayloadv201[v201.GetMonitoringReportReq],
		"GetReport":                     	unmarshalRequestPayloadv201[v201.GetReportReq],
		"GetTransactionStatus":          	unmarshalRequestPayloadv201[v201.GetTransactionStatusReq],
		"GetVariables":                  	unmarshalRequestPayloadv201[v201.GetVariablesReq],
		"Heartbeat":                     	unmarshalRequestPayloadv201[v201.HeartbeatReq],
		"InstallCertificate":            	unmarshalRequestPayloadv201[v201.InstallCertificateReq],
		"LogStatusNotification":         	unmarshalRequestPayloadv201[v201.LogStatusNotificationReq],
		"MeterValues":                   	unmarshalRequestPayloadv201[v201.MeterValuesReq],
		"NotifyChargingLimit":           	unmarshalRequestPayloadv201[v201.NotifyChargingLimitReq],
		"NotifyCustomerInformation":     	unmarshalRequestPayloadv201[v201.NotifyCustomerInformationReq],
		"NotifyDisplayMessages":         	unmarshalRequestPayloadv201[v201.NotifyDisplayMessagesReq],
		"NotifyEVChargingNeeds":         	unmarshalRequestPayloadv201[v201.NotifyEVChargingNeedsReq],
		"NotifyEVChargingSchedule":        	unmarshalRequestPayloadv201[v201.NotifyEVChargingScheduleReq],
		"NotifyEvent":                   	unmarshalRequestPayloadv201[v201.NotifyEventReq],
		"NotifyMonitoringReport":        	unmarshalRequestPayloadv201[v201.NotifyMonitoringReportReq],
		"NotifyReport":                  	unmarshalRequestPayloadv201[v201.NotifyReportReq],
		"PublishFirmware":               	unmarshalRequestPayloadv201[v201.PublishFirmwareReq],
		"PublishFirmawareStatusNotification": unmarshalRequestPayloadv201[v201.PublishFirmwareStatusNotificationReq],
		"ReportChargingProfiles":        	unmarshalRequestPayloadv201[v201.ReportChargingProfilesReq],
		"RequestStartTransaction":       	unmarshalRequestPayloadv201[v201.RequestStartTransactionReq],
		"RequestStopTransaction":        	unmarshalRequestPayloadv201[v201.RequestStopTransactionReq],
		"ReservationStatusUpdate":       	unmarshalRequestPayloadv201[v201.ReservationStatusUpdateReq],
		"ReserveNow":                    	unmarshalRequestPayloadv201[v201.ReserveNowReq],
		"Reset":                         	unmarshalRequestPayloadv201[v201.ResetReq],
		"SecurityEventNotification":     	unmarshalRequestPayloadv201[v201.SecurityEventNotificationReq],
		"SendLocalList":                 	unmarshalRequestPayloadv201[v201.SendLocalListReq],
		"SetChargingProfile":            	unmarshalRequestPayloadv201[v201.SetChargingProfileReq],
		"SetDisplayMessage":             	unmarshalRequestPayloadv201[v201.SetDisplayMessageReq],
		"SetMonitoringBase":             	unmarshalRequestPayloadv201[v201.SetMonitoringBaseReq],
		"SetMonitoringLevel":            	unmarshalRequestPayloadv201[v201.SetMonitoringLevelReq],
		"SetNetworkProfile":             	unmarshalRequestPayloadv201[v201.SetNetworkProfileReq],
		"SetVariableMonitoring":         	unmarshalRequestPayloadv201[v201.SetVariableMonitoringReq],
		"SetVariables":                  	unmarshalRequestPayloadv201[v201.SetVariablesReq],
		"SignCertificate":               	unmarshalRequestPayloadv201[v201.SignCertificateReq],
		"StatusNotification":            	unmarshalRequestPayloadv201[v201.StatusNotificationReq],
		"TransactionEvent":			  		unmarshalRequestPayloadv201[v201.TransactionEventReq],
		"TriggerMessage":                	unmarshalRequestPayloadv201[v201.TriggerMessageReq],
		"UnlockConnector":               	unmarshalRequestPayloadv201[v201.UnlockConnectorReq],
		"UnpublishFirmware":             	unmarshalRequestPayloadv201[v201.UnpublishFirmwareReq],
		"UpdateFirmware":                	unmarshalRequestPayloadv201[v201.UpdateFirmwareReq],
	}

	resmapv201 = map[string]func(json.RawMessage) (Payload, error){
		"Authorize":                     	unmarshalResponsePayloadv201[v201.AuthorizeRes],
		"BootNotification":              	unmarshalResponsePayloadv201[v201.BootNotificationRes],
		"CancelReservation":             	unmarshalResponsePayloadv201[v201.CancelReservationRes],
		"CertificateSigned":			 	unmarshalResponsePayloadv201[v201.CertificateSignedRes],
		"ChangeAvailability":            	unmarshalResponsePayloadv201[v201.ChangeAvailabilityRes],
		"ClearCache":                    	unmarshalResponsePayloadv201[v201.ClearCacheRes],
		"ClearChargingProfile":          	unmarshalResponsePayloadv201[v201.ClearChargingProfileRes],
		"ClearDisplayMessage":           	unmarshalResponsePayloadv201[v201.ClearDisplayMessageRes],
		"ClearedChargingLimit":          	unmarshalResponsePayloadv201[v201.ClearedChargingLimitRes],
		"ClearVariableMonitoring":       	unmarshalResponsePayloadv201[v201.ClearVariableMonitoringRes],
		"CostUpdated":                   	unmarshalResponsePayloadv201[v201.CostUpdatedRes],
		"CustomerInformation":           	unmarshalResponsePayloadv201[v201.CustomerInformationRes],
		"DataTransfer":                  	unmarshalResponsePayloadv201[v201.DataTransferRes],
		"DeleteCertificate":             	unmarshalResponsePayloadv201[v201.DeleteCertificateRes],
		"FirmwareStatusNotification":    	unmarshalResponsePayloadv201[v201.FirmwareStatusNotificationRes],
		"Get15118EVCertificate":         	unmarshalResponsePayloadv201[v201.Get15118EVCertificateRes],
		"GetBaseReport":                 	unmarshalResponsePayloadv201[v201.GetBaseReportRes],
		"GetCertificateStatus":          	unmarshalResponsePayloadv201[v201.GetCertificateStatusRes],
		"GetChargingProfiles":            	unmarshalResponsePayloadv201[v201.GetChargingProfilesRes],
		"GetCompositeSchedule":          	unmarshalResponsePayloadv201[v201.GetCompositeScheduleRes],
		"GetDisplayMessages":            	unmarshalResponsePayloadv201[v201.GetDisplayMessagesRes],
		"GetInstalledCertificateIds":      	unmarshalResponsePayloadv201[v201.GetInstalledCertificateIdsRes],
		"GetLocalListVersion":           	unmarshalResponsePayloadv201[v201.GetLocalListVersionRes],
		"GetLog":							unmarshalResponsePayloadv201[v201.GetLogRes],
		"GetMonitoringReport":           	unmarshalResponsePayloadv201[v201.GetMonitoringReportRes],
		"GetReport":                     	unmarshalResponsePayloadv201[v201.GetReportRes],
		"GetTransactionStatus":          	unmarshalResponsePayloadv201[v201.GetTransactionStatusRes],
		"GetVariables":                  	unmarshalResponsePayloadv201[v201.GetVariablesRes],
		"Heartbeat":                     	unmarshalResponsePayloadv201[v201.HeartbeatRes],
		"InstallCertificate":            	unmarshalResponsePayloadv201[v201.InstallCertificateRes],
		"LogStatusNotification":         	unmarshalResponsePayloadv201[v201.LogStatusNotificationRes],
		"MeterValues":                   	unmarshalResponsePayloadv201[v201.MeterValuesRes],
		"NotifyChargingLimit":           	unmarshalResponsePayloadv201[v201.NotifyChargingLimitRes],
		"NotifyCustomerInformation":     	unmarshalResponsePayloadv201[v201.NotifyCustomerInformationRes],
		"NotifyDisplayMessages":         	unmarshalResponsePayloadv201[v201.NotifyDisplayMessagesRes],
		"NotifyEVChargingNeeds":         	unmarshalResponsePayloadv201[v201.NotifyEVChargingNeedsRes],
		"NotifyEVChargingSchedule":        	unmarshalResponsePayloadv201[v201.NotifyEVChargingScheduleRes],
		"NotifyEvent":                   	unmarshalResponsePayloadv201[v201.NotifyEventRes],
		"NotifyMonitoringReport":        	unmarshalResponsePayloadv201[v201.NotifyMonitoringReportRes],
		"NotifyReport":                  	unmarshalResponsePayloadv201[v201.NotifyReportRes],
		"PublishFirmware":               	unmarshalResponsePayloadv201[v201.PublishFirmwareRes],
		"PublishFirmawareStatusNotification": unmarshalResponsePayloadv201[v201.PublishFirmwareStatusNotificationRes],
		"ReportChargingProfiles":        	unmarshalResponsePayloadv201[v201.ReportChargingProfilesRes],
		"RequestStartTransaction":       	unmarshalResponsePayloadv201[v201.RequestStartTransactionRes],
		"RequestStopTransaction":        	unmarshalResponsePayloadv201[v201.RequestStopTransactionRes],
		"ReservationStatusUpdate":       	unmarshalResponsePayloadv201[v201.ReservationStatusUpdateRes],
		"ReserveNow":                    	unmarshalResponsePayloadv201[v201.ReserveNowRes],
		"Reset":                         	unmarshalResponsePayloadv201[v201.ResetRes],
		"SecurityEventNotification":     	unmarshalResponsePayloadv201[v201.SecurityEventNotificationRes],
		"SendLocalList":                 	unmarshalResponsePayloadv201[v201.SendLocalListRes],
		"SetChargingProfile":            	unmarshalResponsePayloadv201[v201.SetChargingProfileRes],
		"SetDisplayMessage":             	unmarshalResponsePayloadv201[v201.SetDisplayMessageRes],
		"SetMonitoringBase":             	unmarshalResponsePayloadv201[v201.SetMonitoringBaseRes],
		"SetMonitoringLevel":            	unmarshalResponsePayloadv201[v201.SetMonitoringLevelRes],
		"SetNetworkProfile":             	unmarshalResponsePayloadv201[v201.SetNetworkProfileRes],
		"SetVariableMonitoring":         	unmarshalResponsePayloadv201[v201.SetVariableMonitoringRes],
		"SetVariables":                  	unmarshalResponsePayloadv201[v201.SetVariablesRes],
		"SignCertificate":               	unmarshalResponsePayloadv201[v201.SignCertificateRes],
		"StatusNotification":            	unmarshalResponsePayloadv201[v201.StatusNotificationRes],
		"TransactionEvent":			  		unmarshalResponsePayloadv201[v201.TransactionEventRes],
		"TriggerMessage":                	unmarshalResponsePayloadv201[v201.TriggerMessageRes],
		"UnlockConnector":               	unmarshalResponsePayloadv201[v201.UnlockConnectorRes],
		"UnpublishFirmware":             	unmarshalResponsePayloadv201[v201.UnpublishFirmwareRes],
		"UpdateFirmware":                	unmarshalResponsePayloadv201[v201.UpdateFirmwareRes],
	}
}