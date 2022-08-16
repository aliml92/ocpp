package v201

import (
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)


var Validate = validator.New()

func contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
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
	Validate.RegisterValidation("APNAuthenticationEnumType", isAPNAuthenticationEnumType)
	Validate.RegisterValidation("AttributeEnumType", isAttributeEnumType)
	Validate.RegisterValidation("AuthorizationStatusEnumType", isAuthorizationStatusEnumType)
	Validate.RegisterValidation("AuthorizeCertificateStatusEnumType", isAuthorizeCertificateStatusEnumType)
	Validate.RegisterValidation("BootReasonEnumType", isBootReasonEnumType)
	Validate.RegisterValidation("CancelReservationStatusEnumType", isCancelReservationStatusEnumType)
	Validate.RegisterValidation("CertificateActionEnumType", isCertificateActionEnumType)
	Validate.RegisterValidation("CertificateSignedStatusEnumType",  isCertificateSignedStatusEnumType)
	Validate.RegisterValidation("CertificateSigningUseEnumType", isCertificateSigningUseEnumType)
	Validate.RegisterValidation("ChangeAvailabilityStatusEnumType",  isChangeAvailabilityStatusEnumType)
	Validate.RegisterValidation("ChargingLimitSourceEnumType",  isChargingLimitSourceEnumType)
	Validate.RegisterValidation("ChargingProfileKindEnumType", isChargingProfileKindEnumType)
	Validate.RegisterValidation("ChargingProfilePurposeEnumType",  isChargingProfilePurposeEnumType)
	Validate.RegisterValidation("ChargingProfileStatusEnumType",  isChargingProfileStatusEnumType)
	Validate.RegisterValidation("ChargingRateUnitEnumType",  isChargingRateUnitEnumType)
	Validate.RegisterValidation("ChargingStateEnumType",  isChargingStateEnumType)
	Validate.RegisterValidation("ClearCacheStatusEnumType",  isClearCacheStatusEnumType)
	Validate.RegisterValidation("ClearChargingProfileStatusEnumType",  isClearChargingProfileStatusEnumType)
	Validate.RegisterValidation("ClearMessageStatusEnumType",  isClearMessageStatusEnumType)
	Validate.RegisterValidation("ClearMonitoringStatusEnumType",  isClearMonitoringStatusEnumType)
	Validate.RegisterValidation("ComponentCriterionEnumType",  isComponentCriterionEnumType)
	Validate.RegisterValidation("ConnectorEnumType",  isConnectorEnumType)
	Validate.RegisterValidation("ConnectorStatusEnumType",  isConnectorStatusEnumType)
	Validate.RegisterValidation("CostKindEnumType",  isCostKindEnumType)
	Validate.RegisterValidation("CustomerInformationStatusEnumType",  isCustomerInformationStatusEnumType)
	Validate.RegisterValidation("DataEnumType",  isDataEnumType)
	Validate.RegisterValidation("DataTransferStatusEnumType",  isDataTransferStatusEnumType)
	Validate.RegisterValidation("DeleteCertificateStatusEnumType",  isDeleteCertificateStatusEnumType)
	Validate.RegisterValidation("DisplayMessageStatusEnumType",  isDisplayMessageStatusEnumType)
	Validate.RegisterValidation("EnergyTransferModeEnumType",  isEnergyTransferModeEnumType)
	Validate.RegisterValidation("EventNotificationEnumType",  isEventNotificationEnumType)
	Validate.RegisterValidation("EventTriggerEnumType",  isEventTriggerEnumType)
	Validate.RegisterValidation("FirmwareStatusEnumType",  isFirmwareStatusEnumType)
	Validate.RegisterValidation("GenericDeviceModelStatusEnumType",  isGenericDeviceModelStatusEnumType)
	Validate.RegisterValidation("GenericStatusEnumType",  isGenericStatusEnumType)
	Validate.RegisterValidation("GetCertificateIdUseEnumType",  isGetCertificateIdUseEnumType)
	Validate.RegisterValidation("GetCertificateStatusEnumType",  isGetCertificateStatusEnumType)
	Validate.RegisterValidation("GetChargingProfileStatusEnumType",  isGetChargingProfileStatusEnumType)
	Validate.RegisterValidation("GetDisplayMessagesStatusEnumType",  isGetDisplayMessagesStatusEnumType)
	Validate.RegisterValidation("GetInstalledCertificateStatusEnumType",  isGetInstalledCertificateStatusEnumType)
	Validate.RegisterValidation("GetVariableStatusEnumType",  isGetVariableStatusEnumType)
	Validate.RegisterValidation("HashAlgorithmEnumType",  isHashAlgorithmEnumType)
	Validate.RegisterValidation("IdTokenEnumType",  isIdTokenEnumType)
	Validate.RegisterValidation("InstallCertificateStatusEnumType",  isInstallCertificateStatusEnumType)
	Validate.RegisterValidation("InstallCertificateUseEnumType",  isInstallCertificateUseEnumType)
	Validate.RegisterValidation("Iso15118EVCertificateStatusEnumType",  isIso15118EVCertificateStatusEnumType)
	Validate.RegisterValidation("LocationEnumType",  isLocationEnumType)
	Validate.RegisterValidation("LogEnumType",  isLogEnumType)
	Validate.RegisterValidation("LogStatusEnumType",  isLogStatusEnumType)
	Validate.RegisterValidation("MeasurandEnumType",  isMeasurandEnumType)
	Validate.RegisterValidation("MessageFormatEnumType",  isMessageFormatEnumType)
	Validate.RegisterValidation("MessagePriorityEnumType",  isMessagePriorityEnumType)
	Validate.RegisterValidation("MessageStateEnumType",  isMessageStateEnumType)
	Validate.RegisterValidation("MessageTriggerEnumType",  isMessageTriggerEnumType)
	Validate.RegisterValidation("MonitorEnumType",  isMonitorEnumType)
	Validate.RegisterValidation("MonitoringBaseEnumType",  	isMonitoringBaseEnumType)
	Validate.RegisterValidation("MonitoringCriterionEnumType",  isMonitoringCriterionEnumType)
	Validate.RegisterValidation("MutabilityEnumType", isMutabilityEnumType)
	Validate.RegisterValidation("NotifyEVChargingNeedsStatusEnumType",  isNotifyEVChargingNeedsStatusEnumType)
	Validate.RegisterValidation("OCPPInterfaceEnumType",  isOCPPInterfaceEnumType)
	Validate.RegisterValidation("OCPPTransportEnumType",  isOCPPTransportEnumType)
	Validate.RegisterValidation("OCPPVersionEnumType",  isOCPPVersionEnumType)
	Validate.RegisterValidation("OperationalStatusEnumType",  isOperationalStatusEnumType)
	Validate.RegisterValidation("PhaseEnumType",  isPhaseEnumType)
	Validate.RegisterValidation("PublishFirmwareStatusEnumType",  isPublishFirmwareStatusEnumType)
	Validate.RegisterValidation("ReadingContextEnumType",  isReadingContextEnumType)
	Validate.RegisterValidation("ReasonEnumType",  isReasonEnumType)
	Validate.RegisterValidation("RecurrencyKindEnumType",  isRecurrencyKindEnumType)
	Validate.RegisterValidation("RegistrationStatusEnumType",  isRegistrationStatusEnumType)
	Validate.RegisterValidation("ReportBaseEnumType", 	isReportBaseEnumType)
	Validate.RegisterValidation("RequestStartStopStatusEnumType",  isRequestStartStopStatusEnumType)
	Validate.RegisterValidation("ReservationUpdateStatusEnumType",  isReservationUpdateStatusEnumType)
	Validate.RegisterValidation("ReserveNowStatusEnumType",  isReserveNowStatusEnumType)
	Validate.RegisterValidation("ResetEnumType",  isResetEnumType)
	Validate.RegisterValidation("ResetStatusEnumType",  isResetStatusEnumType)
	Validate.RegisterValidation("SendLocalListStatusEnumType",  isSendLocalListStatusEnumType)
	Validate.RegisterValidation("SetMonitoringStatusEnumType",  isSetMonitoringStatusEnumType)
	Validate.RegisterValidation("SetNetworkProfileStatusEnumType",  isSetNetworkProfileStatusEnumType)
	Validate.RegisterValidation("SetVariableStatusEnumType",  isSetVariableStatusEnumType)
	Validate.RegisterValidation("TransactionEventEnumType",  isTransactionEventEnumType)
	Validate.RegisterValidation("TriggerMessageStatusEnumType",  isTriggerMessageStatusEnumType)
	Validate.RegisterValidation("TriggerReasonEnumType",  isTriggerReasonEnumType)
	Validate.RegisterValidation("UnlockStatusEnumType",  isUnlockStatusEnumType)
	Validate.RegisterValidation("UnpublishFirmwareStatusEnumType",  isUnpublishFirmwareStatusEnumType)
	Validate.RegisterValidation("UpdateEnumType",  UpdateEnumType)
	Validate.RegisterValidation("UpdateFirmwareStatusEnumType",  isUpdateFirmwareStatusEnumType)
	Validate.RegisterValidation("UpdateFirmwareStatusEnumType",  isUpdateFirmwareStatusEnumType)
	Validate.RegisterValidation("VPNEnumType",  isVPNEnumType)				

}


func IsISO8601Date(fl validator.FieldLevel) bool {
    ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
    ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
  	return ISO8601DateRegex.MatchString(fl.Field().String())
}


func isAPNAuthenticationEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "CHAP", "NONE", "PAP", "AUTO":
		return true
	default:
		return false
	}
}	

func isAttributeEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Actual", "Target", "MinSet", "MaxSet":
		return true
	default:
		return false
	}
}


func isAuthorizationStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted", 
		"Blocked", 
		"ConcurrentTx",
		"Expired",
		"Invalid",
		"NoCredit",
		"NotAllowedTypeEVSE",
		"NotAtThisLocation",
		"NotAtThisTime",
		"Unknown",
	}
	return contains(cases, enum)
}


func isAuthorizeCertificateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted", 
		"SignatureError", 
		"CertificateExpired",
		"NoCertificateAvailable",
		"CertChainError",
		"ContractCancelled",
	}
	return contains(cases, enum)
}

func isBootReasonEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"ApplicationReset", 
		"FirmwareUpdate", 
		"LocalReset",
		"PowerUp",
		"RemoteReset",
		"ScheduledReset",
		"Triggered",
		"Unknown",
		"Watchdog",
	}
	return contains(cases, enum)
}


func isCancelReservationStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}

func isCertificateActionEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Install", "Update":
		return true
	default:
		return false
	}
}


func isCertificateSignedStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}

func isCertificateSigningUseEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "ChargingStationCertificate", "V2GCertificate":
		return true
	default:
		return false
	}
}


func isChangeAvailabilityStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "Scheduled":
		return true
	default:
		return false
	}
}

func isChargingLimitSourceEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "EMS", "Other", "SO", "CSO":
		return true
	default:
		return false
	}
}


func isChargingProfileKindEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Absolute", "Recurring", "Relative":
		return true
	default:
		return false
	}
}


func isChargingProfilePurposeEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "ChargingStationExternalConstraints", "ChargingStationMaxProfile", "TxDefaultProfile", "TxProfile":
		return true
	default:
		return false
	}
}


func isChargingProfileStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}

func isChargingRateUnitEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "A", "W":
		return true
	default:
		return false
	}
}

func isChargingStateEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Charging",
		"EVConnected",
		"SuspendedEV",
		"SuspendedEVSE",
		"Idle",
	}
	return contains(cases, enum)
}


func isClearCacheStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}


func isClearChargingProfileStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Unknown":
		return true
	default:
		return false
	}
}


func isClearMessageStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Unknown":
		return true
	default:
		return false
	}
}

func isClearMonitoringStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "NotFound":
		return true
	default:
		return false
	}
}

func isComponentCriterionEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Active", "Available", "Enabled", "Problem":
		return true
	default:
		return false
	}
}

func isConnectorEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"cCCS1",
		"cCCS2",
		"cG105",
		"cTesla",
		"cType1",
		"cType2",
		"s309-1P-16A",
		"s309-1P-32A",
		"s309-3P-16A",
		"s309-3P-32A",
		"sBS1361",
		"sCEE-7-7",
		"sType2",
		"sType3",
		"Other1PhMax16A",
		"Other1PhOver16A",
		"Other3Ph",
		"Pan",
		"wInductive",
		"wResonant",
		"Undetermined",
		"Unknown",
	}
	return contains(cases, enum)
}

func isConnectorStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Available",
		"Occupied",
		"Reserved",
		"Unavailable",
		"Faulted",
	}
	return contains(cases, enum)
}

func isCostKindEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "CarbonDioxideEmission", "RelativePricePercentage", "RenewableGenerationPercentage":
		return true
	default:
		return false
	}
}

func isCustomerInformationStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "Invalid":
		return true
	default:
		return false
	}
}


func isDataEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"string",
		"decimal",
		"integer",
		"dateTime",
		"boolean",
		"OptionList",
		"SequenceList",
		"MemberList",
	}
	return contains(cases, enum)
}

func isDataTransferStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "UnknownMessageId", "UnknownVendorId":
		return true
	default:
		return false
	}
}