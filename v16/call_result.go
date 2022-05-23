package v16

import (
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)


var Validate = validator.New()



type AuthorizeConf struct {
	IdTagInfo  *IdTagInfo		    `json:"idTagInfo" validate:"required"`
}


type BootNotificationConf struct {
	CurrentTime string   			`json:"currentTime" validate:"required,ISO8601date"`
	Interval    int                 `json:"interval" validate:"gte=0"`
	Status      RegistrationStatus  `json:"status" validate:"required,registrationStatus"`
}


type StartTransactionConf struct {
	IdTagInfo  	  *IdTagInfo 		`json:"idTagInfo" validate:"required"`
	TransactionId int 				`json:"transactionId" validate:"required"`
}









type ChangeAvailabilityConf struct {
	Status 					string 			`json:"status" validate:"required"`
}


















func init(){


	// register function to get tag name from json tags.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	Validate.RegisterValidation("ISO8601date", IsISO8601Date)
	Validate.RegisterValidation("registrationStatus", isValidRegistrationStatus)
	Validate.RegisterValidation("authorizationStatus", isValidAuthorizationStatus)
}


func IsISO8601Date(fl validator.FieldLevel) bool {
    ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
    ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
  return ISO8601DateRegex.MatchString(fl.Field().String())
}



func isValidRegistrationStatus(fl validator.FieldLevel) bool {
	status := RegistrationStatus(fl.Field().String())
	switch status {
	case "Accepted", "Pending", "Rejected":
		return true
	default:
		return false
	}
}

func isValidAuthorizationStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Blocked", "Expired", "Invalid", "ConcurrentTx":
		return true
	default:
		return false
	}
}