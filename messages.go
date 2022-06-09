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


type cpAction interface {
	v16.BootNotificationReq | v16.AuthorizeReq | v16.DataTransferReq | v16.DiagnosticsStatusNotificationReq 
} 

type csAction interface {
	v16.ChangeAvailabilityConf | v16.ChangeConfigurationConf
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
		5, 
		ce.UniqueId,
		ce.ErrorCode,
		ce.ErrorDescription,
		`{"errorDetails":` + d + `}`,
	}
	raw, _ := json.Marshal(out)
	return &raw
}



// This function converts raw byte to one of the ocpp messages
// There is only one exception where CallResult's payload is returned as json.RawMessage
// because in this function we don't know the type of the payload for CallResult
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
		payload, err = unmarshall_call_payload_from_cp(&mId, &mAction, &rm[3])
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


// Unmarshalls Payload of a Call coming from CP 
func unmarshall_call_payload_from_cp(mId *string, mAction *string, rawPayload *json.RawMessage) (*Payload, error) {
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
		payload, err = unmarshall_cp_action[v16.BootNotificationReq](mId,rawPayload)
		if err != nil {
			return nil, err
		}	
	case "Authorize":
		payload, err = unmarshall_cp_action[v16.AuthorizeReq](mId,rawPayload)
		if err != nil {
			return nil, err
		}
	case "DataTransfer":
		payload, err = unmarshall_cp_action[v16.DataTransferReq](mId, rawPayload)
		if err != nil {
			return nil, err
		}
	case "DiagnosticsStatusNotification":
		payload, err = unmarshall_cp_action[v16.DiagnosticsStatusNotificationReq]( mId,rawPayload)
		if err != nil {
			return nil, err
		}	 	
	}
	return &payload, nil
}

// Unmarshals Payload to a struct of type T, eg. BootNotificationReq
func unmarshall_cp_action[T cpAction](mId *string, rawPayload *json.RawMessage) (*T, error){
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



// Unmarshalls Payload of a CallResult coming from CP 
func unmarshall_call_result_payload_from_cp(mAction *string, rawPayload *json.RawMessage) (*Payload, error) {
	var payload Payload
	var err error
	switch *mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "ChangeAvailability":
		payload, err = unmarshall_cs_action[v16.ChangeAvailabilityConf](rawPayload)
		if err != nil {
			return nil, err
		}
	case "ChangeConfiguration":
		payload, err = unmarshall_cs_action[v16.ChangeConfigurationConf](rawPayload)
		if err != nil {
			return nil, err
		}						
	}
	return &payload, nil
}


// Unmarshals Payload to a struct of type T, eg. ChangeAvailabilityConf
func unmarshall_cs_action[T csAction](rawPayload *json.RawMessage) (*T, error){
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
