package ocpp

import (
	"encoding/json"
	"errors"

	v16 "github.com/aliml92/ocpp/v16"
)


type Call struct {
	MessageTypeId 	uint8
	UniqueId 		string
	Action 			string
	Payload 		interface{}
}

type CallResult struct {
	MessageTypeId 	uint8
	UniqueId 		string
	Payload 		json.RawMessage

}

type CallError struct {
	MessageTypeId 		uint8
	UniqueId 			string
	ErrorCode 			string
	ErrorDescription 	string
	ErrorDetails 		interface{}
}




// converts ocpp bytes to ocpp message struct or error
// TODO CallResult is implemented later
func UnmarshalOCPPMessage(raw []byte) (*Call, *CallResult, *CallError, error) {
	var mm []json.RawMessage
	var call *Call
	var callResult *CallResult
	var callError *CallError
	var payload interface{}

	err := json.Unmarshal(raw, &mm)
	if err != nil {
		// TODO: send a proper CallError 
		return nil, nil, nil,  err
	}
	l := len(mm)
	if l < 3 || l > 5 {
		// TODO: send a proper CallError
		return nil, nil, nil,  errors.New("invalid message")
	}
	// unmarshal the first two elements of message, sice they are always the same
	var mType uint8
	var mId string
	err = json.Unmarshal(mm[0], &mType)
	if err != nil {
		return nil, nil, nil,  err
	}
	err = json.Unmarshal(mm[1], &mId)
	if err != nil {
		return nil, nil, nil,  err
	}
	if mType == 2 {
		var mAction string
		err = json.Unmarshal(mm[2], &mAction)
		if err != nil {
			return nil, nil, nil,  err
		}
		payload, err = UnmarshalReqPayload(mAction, mm[3])
		if err != nil {
			return nil, nil, nil,  err
		}
		call = &Call{
			MessageTypeId: 	mType,
			UniqueId: 		mId,
			Action: 		mAction,
			Payload: 		payload,
		}

	}
	if mType == 3 {
		callResult = &CallResult{
			MessageTypeId: 	mType,
			UniqueId: 		mId,
			Payload: 		mm[2],
		}
	}
	if mType == 4 {
		var me [5]interface{} 
		err = json.Unmarshal(raw, &me)
	
		if err != nil {
			return nil, nil, nil,  err
		}
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


// This function is used to unmarshal the payload of the message by action name
// for now this is used only for incoming Call messages
func UnmarshalReqPayload(mAction string, rawPayload json.RawMessage) (ReqPayload, error) {
	var payload interface{}
	var err error
	switch mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "BootNotification":
		var p * v16.BootNotificationReq
		err = json.Unmarshal(rawPayload, &p)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(p)
		if err != nil {
			return nil, err
		}
		payload = p	
	case "Authorize":
		var p * v16.AuthorizeReq
		err = json.Unmarshal(rawPayload, &p)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(p)
		if err != nil {
			return nil, err
		}
		payload = p
	}
	return payload, nil
}

func UnmarshalResPayload(mAction string, rawPayload json.RawMessage) (ResPayload, error) {
	var payload interface{}
	var err error
	switch mAction {
	default:
		err = errors.New("invalid action")
		return nil, err
	case "ChangeAvailability":
		var p * v16.ChangeAvailabilityConf
		err = json.Unmarshal(rawPayload, &p)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(p)
		if err != nil {
			return nil, err
		}
		payload = p	
	}
	return payload, nil
}
