package ocpp

import (
	"reflect"
	"testing"
)

type RawMessage struct {
	data []byte
	proto string
}


// Add more test cases
func TestUnpack(t *testing.T) {
	cases := []struct{
		name 	string
		rawMsg  RawMessage
		want1   OcppMessage
		want2   error    
	}{
		{	"non array json", 
		  	RawMessage{ []byte(`{"some": "data"}`),"ocppv16",}, 
			nil,
			&ocppError{
				id:    "-1",
				code:  "ProtocolError",
				cause: "Invalid JSON format",
			},
		},
		
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			got1, got2 := unpack(v.rawMsg.data, v.rawMsg.proto)
			
			if got1 != v.want1{
				t.Errorf("got %v want %v", got1, v.want1)
			}

			if !reflect.DeepEqual(got2, v.want2) {
				t.Errorf("got %v want %v", got2, v.want2)
			}
		})
	}
}