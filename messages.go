package ocpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	v16 "github.com/aliml92/ocpp/v16"
	"github.com/google/uuid"
)



type OCPPError struct {
	id 		   string
	code 	   string
	cause 	   string
} 


func (e *OCPPError) Error() string {
	return e.code + ": " + e.cause
}


type Call struct {
	MessageTypeId 	uint8
	UniqueId 		string
	Action 			string
	Payload 		*Payload
}


func (call *Call) CreateCallResult(r *Payload) ( *[]byte) {
	out := [3]interface{}{
		3, 
		call.UniqueId,
		r,
	}
	raw, _ := json.Marshal(out)
	return &raw
}




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



type CallResult struct {
	MessageTypeId 	uint8
	UniqueId 		string
	Payload 		*json.RawMessage

}

type CallError struct {
	MessageTypeId 		uint8
	UniqueId 			string
	ErrorCode 			string
	ErrorDescription 	string
	ErrorDetails 		interface{}
}


func (ce *CallError) Marshal() *[]byte {
	d := ce.ErrorDetails.(string)
	out := [5]interface{}{
		4, 
		ce.UniqueId,
		ce.ErrorCode,
		ce.ErrorDescription,
		`{"cause":` + d + `}`,
	}
	raw, _ := json.Marshal(out)
	return &raw
}


/*
Converts raw byte to one of the ocpp messages or an error if the message is not valid
*/
func unpack(b *[]byte) (*Call, *CallResult, *CallError, error) {
	var rm []json.RawMessage
	var call *Call
	var callResult *CallResult
	var callError *CallError
	var payload *Payload
	var mType uint8
	var mId string
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
	
	err = json.Unmarshal(rm[1], &mId)
	if err != nil {
		e := &OCPPError{
			id:    		"",
			code:		"ProtocolError",
			cause:		"Message does not contain UniqueId",
		}	
		return nil, nil, nil, e
	}
	if 3 > len(rm) || len(rm) > 5 {
		e := &OCPPError{
			id:    mId,
			code: "ProtocolError",
			cause: "JSON must be an array of range [3,5]",
		}
		log.Println(err)  
		return nil, nil, nil, e
	}
	err = json.Unmarshal(rm[0], &mType)
	if err != nil {
		e := &OCPPError{
			id:    mId,
			code: "PropertyConstraintViolation",
			cause: fmt.Sprintf("MessageTypeId: %v is not valid", rm[0]),
		}
		return nil, nil, nil,  e
	}
	if 2 > mType || mType > 4 {
		e := &OCPPError{
			id:    mId,
			code: "ProtocolError",
			cause: "Message does not contain MessageTypeId",
		}
		return nil, nil, nil, e
	}
	err = json.Unmarshal(rm[1], &mId)
	if err != nil {
		e := &OCPPError{
			id:    mId,
			code: "ProtocolError",
			cause: "Message does not contain UniqueId",
		}	
		return nil, nil, nil, e
	}
	if mType == 2 {
		var mAction string
		err = json.Unmarshal(rm[2], &mAction)
		if err != nil {
			e := &OCPPError{
				id:    mId,
				code: "ProtocolError",
				cause: "Message does not contain Action",
			}
			return nil, nil, nil, e
		}
		payload, err = unmarshal_cp_call_payload(&mId, &mAction, &rm[3])
		if err != nil {
			return nil, nil, nil, err
		}
		call = &Call{
			MessageTypeId: 	mType,
			UniqueId: 		mId,
			Action: 		mAction,
			Payload: 		payload,
		}

	}
	if mType == 3 {
		p := &rm[2]
		callResult = &CallResult{
			MessageTypeId: 	mType,
			UniqueId: 		mId,
			Payload: 		p,
		}
	}
	if mType == 4 {
		var me [5]interface{} 
		_ = json.Unmarshal(*b, &me)
		callError = &CallError{
			MessageTypeId: 		mType,
			UniqueId: 			mId,
			ErrorCode: 			me[2].(string),
			ErrorDescription: 	me[3].(string),
			ErrorDetails: 		me[4],
		}
	}
	return call, callResult, callError, nil

}


/*
Converts raw CP initiated Call Payload to a corresponding Payload struct
*/ 
func unmarshal_cp_call_payload(mId *string, mAction *string, rawPayload *json.RawMessage) (*Payload, error) {
	var payload Payload
	var err error
	switch *mAction {
	default:
		e := &OCPPError{
			id:    *mId,
			code: "NotImplemented",
			cause: fmt.Sprintf("Action %v is not implemented", *mAction),
		}
		return nil, e
	case "BootNotification":
		payload, err = unmarshal_cp_action[v16.BootNotificationReq](mId,rawPayload)
		if err != nil {
			return nil, err
		}	
	case "Authorize":
		payload, err = unmarshal_cp_action[v16.AuthorizeReq](mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshal_cp_action[v16.DataTransferReq](mId, rawPayload)
		if err != nil {
			return nil, err
		}
	case "DiagnosticsStatusNotification":
		payload, err = unmarshal_cp_action[v16.DiagnosticsStatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "FirmwareStatusNotification":
		payload, err = unmarshal_cp_action[v16.FirmwareStatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "Heartbeat":
		payload, err = unmarshal_cp_action[v16.HeartbeatReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "MeterValues":
		payload, err = unmarshal_cp_action[v16.MeterValuesReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "StartTransaction":
		payload, err = unmarshal_cp_action[v16.StartTransactionReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "StatusNotification":
		payload, err = unmarshal_cp_action[v16.StatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "StopTransaction":
		payload, err = unmarshal_cp_action[v16.StopTransactionReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}							 	
	}
	return &payload, nil
}

/*
Unmarshals Payload to a struct of type T, eg. BootNotificationReq
*/
func unmarshal_cp_action[T any](mId *string, rawPayload *json.RawMessage) (*T, error){
	var p *T
	err := json.Unmarshal(*rawPayload, &p)
	if err != nil {
		e := &OCPPError{
			id:    *mId,
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
			id:    *mId,
			code: "PropertyConstraintViolationError",
			cause: "Call Payload is not valid",
		}
		log.Println(err)
		return nil, e
	}
	return p, nil
}



/*
Converts raw CallResult Payload (response to CSMS initiated action) to a corresponding Payload struct
Flow: CP     <--(Call)--    CSMS 
      CP  --(CallResult)--> CSMS
*/ 
func unmarshal_csms_call_result_payload(mAction *string, rawPayload *json.RawMessage) (*Payload, error) {
	var payload Payload
	var err error
	switch *mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "CancelReservation":
		payload, err = unmarshal_cs_action[v16.CancelReservationConf](rawPayload)
		if err != nil {
			return nil, err
		}	
	case "ChangeAvailability":
		payload, err = unmarshal_cs_action[v16.ChangeAvailabilityConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ChangeConfiguration":
		payload, err = unmarshal_cs_action[v16.ChangeConfigurationConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearCache":
		payload, err = unmarshal_cs_action[v16.ClearCacheConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearChargingProfile":
		payload, err = unmarshal_cs_action[v16.ClearChargingProfileConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshal_cs_action[v16.DataTransferConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetCompositeSchedule":
		payload, err = unmarshal_cs_action[v16.GetCompositeScheduleConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetConfiguration":
		payload, err = unmarshal_cs_action[v16.GetConfigurationConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetDiagnostics":
		payload, err = unmarshal_cs_action[v16.GetDiagnosticsConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetLocalListVersion":
		payload, err = unmarshal_cs_action[v16.GetLocalListVersionConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStartTransaction":
		payload, err = unmarshal_cs_action[v16.RemoteStartTransactionConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStopTransaction":
		payload, err = unmarshal_cs_action[v16.RemoteStopTransactionConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ReserveNow":
		payload, err = unmarshal_cs_action[v16.ReserveNowConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "Reset":
		payload, err = unmarshal_cs_action[v16.ResetConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SendLocalList":
		payload, err = unmarshal_cs_action[v16.SendLocalListConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SetChargingProfile":
		payload, err = unmarshal_cs_action[v16.SetChargingProfileConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "TriggerMessage":
		payload, err = unmarshal_cs_action[v16.TriggerMessageConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UnlockConnector":
		payload, err = unmarshal_cs_action[v16.UnlockConnectorConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UpdateFirmware":
		payload, err = unmarshal_cs_action[v16.UpdateFirmwareConf](rawPayload)
		if err != nil {
			return nil, err
		}																				
	}
	return &payload, nil
}


/*
Converts raw Call Payload (CSMS initiated action) to a corresponding Payload struct
Flow:                       CSMS     <--(Call*)--   ThirdParty      // Call* represents CSMS initiated action, can be delivered via various means 
      CP     <--(Call)--    CSMS                                    // Eg. via HTTP, MQTT, Websocket, etc.
*/ 
func UnmarshalCallPayloadFromThirdParty(mAction *string, rawPayload *json.RawMessage) (*Payload, error) {
	var payload Payload
	var err error
	switch *mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "CancelReservation":
		payload, err = unmarshal_cs_action[v16.CancelReservationReq](rawPayload)
		if err != nil {
			return nil, err
		}	
	case "ChangeAvailability":
		payload, err = unmarshal_cs_action[v16.ChangeAvailabilityReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ChangeConfiguration":
		payload, err = unmarshal_cs_action[v16.ChangeConfigurationReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearCache":
		payload, err = unmarshal_cs_action[v16.ClearCacheReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ClearChargingProfile":
		payload, err = unmarshal_cs_action[v16.ClearChargingProfileReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshal_cs_action[v16.DataTransferReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetCompositeSchedule":
		payload, err = unmarshal_cs_action[v16.GetCompositeScheduleReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetConfiguration":
		payload, err = unmarshal_cs_action[v16.GetConfigurationReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetDiagnostics":
		payload, err = unmarshal_cs_action[v16.GetDiagnosticsReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "GetLocalListVersion":
		payload, err = unmarshal_cs_action[v16.GetLocalListVersionReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStartTransaction":
		payload, err = unmarshal_cs_action[v16.RemoteStartTransactionReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "RemoteStopTransaction":
		payload, err = unmarshal_cs_action[v16.RemoteStopTransactionReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ReserveNow":
		payload, err = unmarshal_cs_action[v16.ReserveNowReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "Reset":
		payload, err = unmarshal_cs_action[v16.ResetReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SendLocalList":
		payload, err = unmarshal_cs_action[v16.SendLocalListReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "SetChargingProfile":
		payload, err = unmarshal_cs_action[v16.SetChargingProfileReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "TriggerMessage":
		payload, err = unmarshal_cs_action[v16.TriggerMessageReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UnlockConnector":
		payload, err = unmarshal_cs_action[v16.UnlockConnectorReq](rawPayload)
		if err != nil {
			return nil, err
		}
	case "UpdateFirmware":
		payload, err = unmarshal_cs_action[v16.UpdateFirmwareReq](rawPayload)
		if err != nil {
			return nil, err
		}																				
	}
	return &payload, nil
}

/*
Unmarshals Payload to a struct of type T, eg. ChangeAvailabilityConf
*/
func unmarshal_cs_action[T any](rawPayload *json.RawMessage) (*T, error){
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
