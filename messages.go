package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aliml92/ocpp/v16"
	"github.com/google/uuid"
)



// intented to be used as a generic error container
// internal use only
type OCPPError struct {
	id 		   string
	code 	   string
	cause 	   string
} 


func (e *OCPPError) Error() string {
	return e.code + ": " + e.cause
}

// Represents OCPP Call
type Call struct {
	MessageTypeId 	uint8
	UniqueId 		string
	Action 			string
	Payload 		Payload
}

// Create CallResult from a received Call
func (call *Call) CreateCallResult(r Payload) ( *[]byte) {
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
func (call *Call) CreateCallError(err error) ( *[]byte) {
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
			UniqueId: id,
			ErrorCode: code,
			ErrorDescription: "",
			ErrorDetails: cause,
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
	return callError.Marshal()
}


// Represents OCPP CallResult
type CallResult struct {
	MessageTypeId 	uint8
	UniqueId 		string
	Payload 		*json.RawMessage

}

// Represents OCPP CallError
type CallError struct {
	MessageTypeId 		uint8
	UniqueId 			string
	ErrorCode 			string
	ErrorDescription 	string
	ErrorDetails 		interface{}
}


func (ce *CallError) Marshal() *[]byte {
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
	Message 	string 
}


func (e *TimeoutError) Error() string {
	return fmt.Sprintf("TimeoutError: %s", e.Message)
}

type DisconnectedError struct {
	Message string
}


func (e *DisconnectedError) Error() string {
	return fmt.Sprintf("DisconnectedError: %s", e.Message)
}

// Converts raw byte to one of the ocpp messages or an error if the message is not valid
// [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}] 
func unpack(b *[]byte) (*Call, *CallResult, *CallError, error) {
	var rm []json.RawMessage
	var mti uint8		//  MessageTypeId
	var ui string		//  UniqueId
	var a string 		//  Action
	var p Payload		//  Payload
	var c *Call         
	var cr *CallResult
	var ce *CallError
	err := json.Unmarshal(*b, &rm)
	if err != nil {
		e := &OCPPError{
			id:    "",    
			code: "ProtocolError",
			cause: "Invalid JSON format",
		}
		log.Println(err)  
		return nil, nil, nil, e
	}
	err = json.Unmarshal(rm[1], &ui)
	if err != nil {
		e := &OCPPError{
			id:    ui,
			code: "ProtocolError",
			cause: "Message does not contain UniqueId",
		}	
		return nil, nil, nil, e
	}
	if 3 > len(rm) || len(rm) > 5 {
		e := &OCPPError{
			id:    ui,
			code: "ProtocolError",
			cause: "JSON must be an array of range [3,5]",
		}
		log.Println(err)  
		return nil, nil, nil, e
	}
	err = json.Unmarshal(rm[0], &mti)
	if err != nil {
		e := &OCPPError{
			id:    ui,
			code: "PropertyConstraintViolation",
			cause: fmt.Sprintf("MessageTypeId: %v is not valid", rm[0]),
		}
		return nil, nil, nil,  e
	}
	if 2 > mti || mti > 4 {
		e := &OCPPError{
			id:    ui,
			code: "ProtocolError",
			cause: "Message does not contain MessageTypeId",
		}
		return nil, nil, nil, e
	}
	if mti == 2 {
		err = json.Unmarshal(rm[2], &a)
		if err != nil {
			e := &OCPPError{
				id:    ui,
				code: "ProtocolError",
				cause: "Message does not contain Action",
			}
			return nil, nil, nil, e
		}
		p, err = unmarshalCPCallPayload(ui, a, &rm[3])
		println(p)
		if err != nil {
			return nil, nil, nil, err
		}
		c = &Call{
			MessageTypeId: 	mti,
			UniqueId: 		ui,
			Action: 		a,
			Payload: 		p,
		}

	}
	if mti == 3 {
		p := &rm[2]
		cr = &CallResult{
			MessageTypeId: 	mti,
			UniqueId: 		ui,
			Payload: 		p,
		}
	}
	if mti == 4 {
		var me [5]interface{} 
		_ = json.Unmarshal(*b, &me)
		ce = &CallError{
			MessageTypeId: 		mti,
			UniqueId: 			ui,
			ErrorCode: 			me[2].(string),
			ErrorDescription: 	me[3].(string),
			ErrorDetails: 		me[4],
		}
	}
	return c, cr, ce, nil

}






// Converts raw CP initiated Call Payload to a corresponding Payload struct
func unmarshalCPCallPayload(mId string, mAction string, rawPayload *json.RawMessage) (Payload, error) {
	var payload Payload
	var err error
	switch mAction {
	default:
		e := &OCPPError{
			id:    mId,
			code: "NotImplemented",
			cause: fmt.Sprintf("Action %v is not implemented", mAction),
		}
		return nil, e
	case "BootNotification":
		payload, err = unmarshalCPAction[v16.BootNotificationReq](mId,rawPayload)
		if err != nil {
			return nil, err
		}	
	case "Authorize":
		payload, err = unmarshalCPAction[v16.AuthorizeReq](mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshalCPAction[v16.DataTransferReq](mId, rawPayload)
		if err != nil {
			return nil, err
		}
	case "DiagnosticsStatusNotification":
		payload, err = unmarshalCPAction[v16.DiagnosticsStatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "FirmwareStatusNotification":
		payload, err = unmarshalCPAction[v16.FirmwareStatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "Heartbeat":
		payload, err = unmarshalCPAction[v16.HeartbeatReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "MeterValues":
		payload, err = unmarshalCPAction[v16.MeterValuesReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "StartTransaction":
		payload, err = unmarshalCPAction[v16.StartTransactionReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "StatusNotification":
		payload, err = unmarshalCPAction[v16.StatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "StopTransaction":
		payload, err = unmarshalCPAction[v16.StopTransactionReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}							 	
	}
	return payload, nil
}


// Unmarshals Payload to a struct of type T, eg. BootNotificationReq
func unmarshalCPAction[T any](mId string, rawPayload *json.RawMessage) (*T, error){
	var p *T
	err := json.Unmarshal(*rawPayload, &p)
	if err != nil {
		e := &OCPPError{
			id:    mId,
			code: "TypeConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Println(err)
		return nil, e
	}
	err = validate.Struct(*p)
	if err != nil {
		// TODO: construct more detailed error
		e := &OCPPError{
			id:    mId,
			code: "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Println(err)
		return nil, e
	}
	return p, nil
}




// Converts raw CallResult Payload (response to CSMS initiated action) to a corresponding Payload struct
func unmarshalCSMSCallResultPayload(mAction string, rawPayload *json.RawMessage) (Payload, error) {
	var payload Payload
	var err error
	switch mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "CancelReservation":
		payload, err = unmarshalCSMSAction[v16.CancelReservationConf](rawPayload)
		if err != nil {
			return nil, err
		}	
	case "ChangeAvailability":
		payload, err = unmarshalCSMSAction[v16.ChangeAvailabilityConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ChangeConfiguration":
		payload, err = unmarshalCSMSAction[v16.ChangeConfigurationConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearCache":
		payload, err = unmarshalCSMSAction[v16.ClearCacheConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearChargingProfile":
		payload, err = unmarshalCSMSAction[v16.ClearChargingProfileConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshalCSMSAction[v16.DataTransferConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetCompositeSchedule":
		payload, err = unmarshalCSMSAction[v16.GetCompositeScheduleConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetConfiguration":
		payload, err = unmarshalCSMSAction[v16.GetConfigurationConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetDiagnostics":
		payload, err = unmarshalCSMSAction[v16.GetDiagnosticsConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetLocalListVersion":
		payload, err = unmarshalCSMSAction[v16.GetLocalListVersionConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStartTransaction":
		payload, err = unmarshalCSMSAction[v16.RemoteStartTransactionConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStopTransaction":
		payload, err = unmarshalCSMSAction[v16.RemoteStopTransactionConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ReserveNow":
		payload, err = unmarshalCSMSAction[v16.ReserveNowConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "Reset":
		payload, err = unmarshalCSMSAction[v16.ResetConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SendLocalList":
		payload, err = unmarshalCSMSAction[v16.SendLocalListConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SetChargingProfile":
		payload, err = unmarshalCSMSAction[v16.SetChargingProfileConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "TriggerMessage":
		payload, err = unmarshalCSMSAction[v16.TriggerMessageConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UnlockConnector":
		payload, err = unmarshalCSMSAction[v16.UnlockConnectorConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UpdateFirmware":
		payload, err = unmarshalCSMSAction[v16.UpdateFirmwareConf](rawPayload)
		if err != nil {
			return nil, err
		}																				
	}
	return payload, nil
}



// Converts raw Call Payload (CSMS initiated action) to a corresponding Payload struct
func UnmarshalCSMSCallPayload(mAction string, rawPayload *json.RawMessage) (Payload, error) {
	var payload Payload
	var err error
	switch mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "CancelReservation":
		payload, err = unmarshalCSMSAction[v16.CancelReservationReq](rawPayload)
		if err != nil {
			return nil, err
		}	
	case "ChangeAvailability":
		payload, err = unmarshalCSMSAction[v16.ChangeAvailabilityReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ChangeConfiguration":
		payload, err = unmarshalCSMSAction[v16.ChangeConfigurationReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearCache":
		payload, err = unmarshalCSMSAction[v16.ClearCacheReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearChargingProfile":
		payload, err = unmarshalCSMSAction[v16.ClearChargingProfileReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshalCSMSAction[v16.DataTransferReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetCompositeSchedule":
		payload, err = unmarshalCSMSAction[v16.GetCompositeScheduleReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetConfiguration":
		payload, err = unmarshalCSMSAction[v16.GetConfigurationReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetDiagnostics":
		payload, err = unmarshalCSMSAction[v16.GetDiagnosticsReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetLocalListVersion":
		payload, err = unmarshalCSMSAction[v16.GetLocalListVersionReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStartTransaction":
		payload, err = unmarshalCSMSAction[v16.RemoteStartTransactionReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStopTransaction":
		payload, err = unmarshalCSMSAction[v16.RemoteStopTransactionReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ReserveNow":
		payload, err = unmarshalCSMSAction[v16.ReserveNowReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "Reset":
		payload, err = unmarshalCSMSAction[v16.ResetReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SendLocalList":
		payload, err = unmarshalCSMSAction[v16.SendLocalListReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SetChargingProfile":
		payload, err = unmarshalCSMSAction[v16.SetChargingProfileReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "TriggerMessage":
		payload, err = unmarshalCSMSAction[v16.TriggerMessageReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UnlockConnector":
		payload, err = unmarshalCSMSAction[v16.UnlockConnectorReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UpdateFirmware":
		payload, err = unmarshalCSMSAction[v16.UpdateFirmwareReq](rawPayload)
		if err != nil {
			return nil, err
		}																				
	}
	return payload, nil
}


// Unmarshals Payload to a struct of type T, eg. ChangeAvailabilityConf
func unmarshalCSMSAction[T any](rawPayload *json.RawMessage) (*T, error){
	var p *T
	err := json.Unmarshal(*rawPayload, &p)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(*p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
