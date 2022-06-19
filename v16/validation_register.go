package v16

import (
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)


var Validate = validator.New()



// TODO: replace some validators with GenericStatusEnumType  
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
	Validate.RegisterValidation("RegistrationStatus", isValidRegistrationStatus)
	Validate.RegisterValidation("AuthorizationStatus", isValidAuthorizationStatus)
	Validate.RegisterValidation("DiagnosticsStatus", isValidDiagnosticsStatus)
	Validate.RegisterValidation("FirmwareStatus", isValidFirmwareStatus)
	Validate.RegisterValidation("ReadingContext", isValidReadingContext)
	Validate.RegisterValidation("ValueFormat", isValidValueFormat)
	Validate.RegisterValidation("Measurand", isValidMeasurand)
	Validate.RegisterValidation("Phase", isValidPhase)
	Validate.RegisterValidation("Location", isValidLocation)
	Validate.RegisterValidation("UnitOfMeasure", isValidUnitOfMeasure)
	Validate.RegisterValidation("ChargePointErrorCode", isValidChargePointErrorCode)
	Validate.RegisterValidation("ChargePointStatus", isValidChargePointStatus)
	Validate.RegisterValidation("Reason", isValidReason)
	Validate.RegisterValidation("DataTransferStatus", isValidDataTransferStatus)
	Validate.RegisterValidation("AvailabilityType", isValidAvailabilityType)
	Validate.RegisterValidation("AvailabilityStatus", isValidAvailabilityStatus)
	Validate.RegisterValidation("ConfigurationStatus", isValidConfigurationStatus)
	Validate.RegisterValidation("ClearCacheStatus", isValidGenericStatusEnumType)            // generic status enum type
	Validate.RegisterValidation("ChargingProfilePurposeType", isValidChargingProfilePurposeType)
	Validate.RegisterValidation("ChargingRateUnitType", isValidChargingRateUnitType)
	Validate.RegisterValidation("ChargingProfileKindType", isValidChargingProfileKindType)
	Validate.RegisterValidation("RecurrencyKindType", isValidRecurrencyKindType)
	Validate.RegisterValidation("ResetType", isValidResetType)
	Validate.RegisterValidation("MessageTrigger", isValidMessageTrigger)
	Validate.RegisterValidation("ClearChargingProfileStatus", isValidClearChargingProfileStatus)
	Validate.RegisterValidation("RemoteStartStopStatus", isValidGenericStatusEnumType)        // generic status enum type
	Validate.RegisterValidation("ReservationStatus", isValidReservationStatus)
	Validate.RegisterValidation("ResetStatus", isValidGenericStatusEnumType)                 // generic status enum type
	Validate.RegisterValidation("UpdateStatus", isValidUpdateStatus)
	Validate.RegisterValidation("ChargingProfileStatus", isValidChargingProfileStatus)
	Validate.RegisterValidation("TriggerMessageStatus", isValidTriggerMessageStatus)
	Validate.RegisterValidation("UnlockStatus", isValidUnlockStatus)
	Validate.RegisterValidation("CancelReservationStatus", isValidGenericStatusEnumType)         // generic status enum type
	Validate.RegisterValidation("GetCompositeScheduleStatus", isValidGenericStatusEnumType)      // generic status enum type
	Validate.RegisterValidation("FirmwareStatus", isValidFirmwareStatusEnumType)
	Validate.RegisterValidation("CertificateSignedStatusEnumType", isValidGenericStatusEnumType) // generic status enum type
	Validate.RegisterValidation("CertificateStatusEnumType", isValidCertificateStatusEnumType)
	Validate.RegisterValidation("CertificateUseTypeEnumType", isValidCertificateUseTypeEnumType)
	Validate.RegisterValidation("DeleteCertificateStatusEnumType", isValidDeleteCertificateStatusEnumType)
	Validate.RegisterValidation("FirmwareStatusEnumType", isValidFirmwareStatusEnumType)
	Validate.RegisterValidation("GenericStatusEnumType", isValidGenericStatusEnumType)
	Validate.RegisterValidation("GetInstalledCertificateStatusEnumType", isValidGetInstalledCertificateStatusEnumType)
	Validate.RegisterValidation("HashAlgorithmEnumType", isValidHashAlgorithmEnumType)
	Validate.RegisterValidation("LogEnumType", isValidLogEnumType)
	Validate.RegisterValidation("LogStatusEnumType", isValidLogStatusEnumType)
	Validate.RegisterValidation("MessageTriggerEnumType", isValidMessageTriggerEnumType)
	Validate.RegisterValidation("TriggerMessageStatusEnumType", isValidTriggerMessageStatusEnumType)
	Validate.RegisterValidation("UpdateFirmwareStatusEnumType", isValidUpdateFirmwareStatusEnumType)
	Validate.RegisterValidation("UploadLogStatusEnumType", isValidUploadLogStatusEnumType)
}


func IsISO8601Date(fl validator.FieldLevel) bool {
    ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
    ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
  	return ISO8601DateRegex.MatchString(fl.Field().String())
}


func isValidRegistrationStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
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

func isValidDiagnosticsStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Idle", "Uploaded", "UploadFailed", "Uploading":
		return true
	default:
		return false
	}
}

func isValidFirmwareStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	cases := []string{
		"Downloaded", 
		"DownloadFailed", 
		"Downloading",
		"Idle", 
		"InstallationFailed", 
		"Installing", 
		"Installed",
	}
	return contains(cases, status)
}

func isValidReadingContext(fl validator.FieldLevel) bool {
	context := fl.Field().String()
	cases := []string{
		"Interruption.Begin", 
		"Interruption.End", 
		"Other",
		"Sample.Clock",
		"Sample.Periodic",
		"Transaction.Begin",
		"Transaction.End",
		"Trigger",
	}
	return contains(cases, context)
}

func isValidValueFormat(fl validator.FieldLevel) bool {
	format := fl.Field().String()
	switch format {
	case "Raw", "SignedData":
		return true
	default:
		return false
	}
}

func isValidMeasurand(fl validator.FieldLevel) bool {
	measurand := fl.Field().String()
	cases := []string{
		"Energy.Active.Import.Register",
		"Energy.Active.Export.Register",
		"Energy.Reactive.Import.Register",
		"Energy.Reactive.Export.Register",
		"Energy.Active.Import.Interval",
		"Energy.Active.Export.Interval",
		"Energy.Reactive.Import.Interval",
		"Energy.Reactive.Export.Interval",
		"Frequency",
		"Power.Active.Export",
		"Power.Active.Import",
		"Power.Reactive.Export",
		"Power.Reactive.Import",
		"Power.Factor",
		"Power.Offered",
		"RPM",
		"SoC",
		"Temperature",
		"Voltage",}
	return contains(cases, measurand)
}


func isValidPhase(fl validator.FieldLevel) bool {
	phase := fl.Field().String()
	cases := []string{
		"L1",
		"L2",
		"L3",
		"N",
		"L1-N",
		"L2-N",
		"L3-N",
		"L1-L2",
		"L2-L3",
		"L3-L1",}
	return contains(cases, phase)
}


func isValidLocation(fl validator.FieldLevel) bool {
	location := fl.Field().String()
	cases := []string{
		"Body",
		"Cable",
		"EV",
		"Inlet",
		"Outlet",
	}
	return contains(cases, location)
}



func isValidUnitOfMeasure(fl validator.FieldLevel) bool {
	unit := fl.Field().String()
	cases := []string{
		"Wh",
		"kWh",
		"varh",
		"kvarh",
		"W",
		"VA",
		"kVA",
		"var",
		"kvar",
		"A",
		"V",
		"Celsius",
		"Fahrenheit",
		"K",
		"Percent",
	}
	return contains(cases, unit)
}



func isValidChargePointErrorCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	cases := []string{
		"ConnectorLockFailure",
		"EVCommunicationError",
		"GroundFailure",
		"HighTemperature",
		"InternalError",
		"LocalListConflict",
		"NoError",
		"OtherError",
		"OverCurrentFailure",
		"OverVoltage",
		"PowerMeterFailure",
		"PowerSwitchFailure",
		"ReaderFailure",
		"ResetFailure",
		"UnderVoltage",
		"WeakSignal",
	}
	return contains(cases, code)
}



func isValidChargePointStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	cases := []string{
		"Available",
		"Preparing",
		"Charging",
		"SuspendedEVSE",
		"SuspendedEV",
		"Finishing",
		"Reserved",
		"Unavailable",
		"Faulted",
	}
	return contains(cases, status)
}


func isValidReason(fl validator.FieldLevel) bool {
	reason := fl.Field().String()
	cases := []string{
		"DeAuthorized",
		"EmergencyStop",
		"EVDisconnected",
		"HardReset",
		"Local",
		"Other",
		"PowerLoss",
		"Reboot",
		"Remote",
		"SoftReset",
		"UnlockCommand",
	}
	return contains(cases, reason)
}




func isValidDataTransferStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "UnknownMessageId", "UnknownVendorId":
		return true
	default:
		return false
	}	
}



func isValidAvailabilityType(fl validator.FieldLevel) bool {
	type_ := fl.Field().String()
	switch type_ {
	case "Inoperative", "Operative":
		return true
	default:	
		return false
	}

}

func isValidAvailabilityStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "Scheduled":
		return true
	default:
		return false
	}
}


func isValidConfigurationStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "RebootRequired", "NotSupported":
		return true
	default:
		return false
	}

}



func isValidChargingProfilePurposeType(fl validator.FieldLevel) bool {
	purpose := fl.Field().String()
	switch purpose {
	case "ChargePointMaxProfile", "TxDefaultProfile", "TxProfile":
		return true
	default:
		return false
	}
}


func isValidChargingRateUnitType(fl validator.FieldLevel) bool {
	kind := fl.Field().String()
	switch kind {
	case "W", "A":
		return true
	default:
		return false
	}
}


func isValidChargingProfileKindType(fl validator.FieldLevel) bool {
	kind := fl.Field().String()
	switch kind {
	case "Absolute", "Recurring", "Relative":
		return true
	default:
		return false
	}
}


func isValidRecurrencyKindType(fl validator.FieldLevel) bool {
	kind := fl.Field().String()
	switch kind {
	case "Daily", "Weekly":
		return true
	default:
		return false
	}
}


func isValidResetType(fl validator.FieldLevel) bool {
	kind := fl.Field().String()
	switch kind {
	case "Hard", "Soft":
		return true
	default:
		return false
	}
}

func isValidMessageTrigger(fl validator.FieldLevel) bool {
	trigger := fl.Field().String()
	cases := []string{
		"BootNotification",
		"DiagnosticsStatusNotification",
		"FirmwareStatusNotification",
		"Heartbeat",
		"MeterValues",
		"StatusNotification",
	}
	return contains(cases, trigger)
}

 func isValidClearChargingProfileStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Unknown":
		return true
	default:
		return false
	}
}




func isValidReservationStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Faulted", "Occupied", "Rejected", "Unavailable":
		return true
	default:
		return false
	}
}




func isValidUpdateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Failed", "NotSupported", "VersionMismatch":
		return true
	default:
		return false
	}
}


func isValidChargingProfileStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "NotSupported":
		return true
	default:
		return false
	}
}

func isValidTriggerMessageStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "NotImplemented":
		return true
	default:
		return false
	}
}


func isValidUnlockStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Unlocked", "UnlockFailed", "NotSupported":
		return true
	default:
		return false
	}
}




func contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}



func isValidCertificateStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Failed", "Rejected":
		return true
	default:
		return false
	}
}

func isValidCertificateUseTypeEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "CentralSystemRootCertificate", "ManufacturerRootCertificate":
		return true
	default:
		return false
	}
}

func isValidDeleteCertificateStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Failed", "NotFound":
		return true
	default:
		return false
	}
}

func isValidFirmwareStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	cases := []string{
		"Downloaded",
		"DownloadFailed",
		"Downloading",
		"DownloadScheduled",
		"DownloadPaused",
		"Idle",
		"InstallationFailed",
	}
	return contains(cases, status)
}

func isValidGenericStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}

func isValidGetInstalledCertificateStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "NotFound":
		return true
	default:
		return false
	}
}

func isValidHashAlgorithmEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "SHA256", "SHA384", "SHA512":
		return true
	default:
		return false
	}
}

func isValidLogEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "DiagnosticsLog", "SecurityLog":
		return true
	default:
		return false
	}
}

func isValidLogStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "AcceptedCanceled":
		return true
	default:
		return false
	}
}

func isValidMessageTriggerEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	cases := []string{
		"BootNotification",
		"LogStatusNotification",
		"FirmwareStatusNotification",
		"Heartbeat",
		"MeterValues",
		"SignChargePointCertificate",
		"StatusNotification",
	}
	return contains(cases, status)
}

func isValidTriggerMessageStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "Accepted", "Rejected", "NotImplemented":
		return true
	default:
		return false
	}
}

func isValidUpdateFirmwareStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	cases := []string{
		"Accepted",
		"Rejected",
		"AcceptedCancled",
		"InvalidCertificate",
		"RevokedCertificate",
	}
	return contains(cases, status)
}


func isValidUploadLogStatusEnumType(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	cases := []string{
		"BadMessage",
		"Idle",
		"NotSupportedOperation",
		"PermissionDenied",
		"Uploaded",
		"UploadFailure",
		"Uploading",
	}
	return contains(cases, status)
}