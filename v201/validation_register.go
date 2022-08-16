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
	Validate.RegisterValidation("UpdateEnumType",  isUpdateEnumType)
	Validate.RegisterValidation("UpdateFirmwareStatusEnumType",  isUpdateFirmwareStatusEnumType)
	Validate.RegisterValidation("UpdateFirmwareStatusEnumType",  isUpdateFirmwareStatusEnumType)
	Validate.RegisterValidation("UploadLogStatusEnumType",  isUploadLogStatusEnumType)
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

func isDeleteCertificateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Failed", "NotFound":
		return true
	default:
		return false
	}
}


func isDisplayMessageStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted",
		"NotSupportedMessageFormat",
		"Rejected",
		"NotSupportedPriority",
		"NotSupportedState",
		"UnknownTransaction",
	}
	return contains(cases, enum)
}

func isEnergyTransferModeEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "DC", "AC_single_phase", "AC_two_phase", "AC_three_phase":
		return true
	default:
		return false
	}
}


func isEventNotificationEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "HardWiredNotification", "HardWiredMonitor", "PreconfiguredMonitor", "CustomMonitor":
		return true
	default:
		return false
	}
}

func isEventTriggerEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Alerting", "Delta", "Periodic":
		return true
	default:
		return false
	}
}

func isFirmwareStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Downloaded",
		"DownloadFailed",
		"Downloading",
		"DownloadScheduled",
		"DownloadPaused",
		"Idle",
		"InstallationFailed",
		"Installing",
		"Installed",
		"Installed",
		"InstallRebooting",
		"InstallScheduled",
		"InstallVerificationFailed",
		"InvalidSignature",
		"SignatureVerified",
	}
	return contains(cases, enum)
}

func isGenericDeviceModelStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "NotSupported", "EmptyResultSet":
		return true
	default:
		return false
	}
}


func isGenericStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}

func isGetCertificateIdUseEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"V2GRootCertificate",
		"MORootCertificate",
		"CSMSRootCertificate",
		"V2GCertificateChain",
		"ManufacturerRootCertificate",
	}
	return contains(cases, enum)
}

func isGetCertificateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Failed":
		return true
	default:
		return false
	}
}

func isGetChargingProfileStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "NoProfiles":
		return true
	default:
		return false
	}
}

func isGetDisplayMessagesStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Unknown":
		return true
	default:
		return false
	}
}

func isGetInstalledCertificateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "NotFound":
		return true
	default:
		return false
	}
}

func isGetVariableStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted",
		"Rejected",
		"UnknownComponent",
		"UnknownVariable",
		"NotSupportedAttributeType",
	}
	return contains(cases, enum)
}


func isHashAlgorithmEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "SHA256", "SHA384", "SHA512":
		return true
	default:
		return false
	}
}


func isIdTokenEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Central",
		"eMAID",
		"ISO14443",
		"ISO15693",
		"KeyCode",
		"Local",
		"MacAddress",
		"NoAuthorization",
	}
	return contains(cases, enum)
}

func isInstallCertificateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "Failed":
		return true
	default:
		return false
	}
}

func isInstallCertificateUseEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "V2GRootCertificate", "MORootCertificate", "CSMSRootCertificate", "ManufacturerRootCertificate":
		return true
	default:
		return false
	}
}


func isIso15118EVCertificateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Failed":
		return true
	default:
		return false
	}
}


func isLocationEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Body",
		"Cable",
		"EV",
		"Inlet",
		"Outlet",
	}
	return contains(cases, enum)
}


func isLogEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "DiagnosticsLog", "SecurityLog":
		return true
	default:
		return false
	}
}


func isLogStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "AcceptedCancelled":
		return true
	default:
		return false
	}
}

func isMeasurandEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Current.Export",
		"Current.Import",
		"Current.Offered",
		"Energy.Active.Export.Register",
		"Energy.Active.Import.Register",
		"Enerfy.Reactive.Export.Register",
		"Energy.Reactive.Import.Register",
		"Energy.Active.Export.Interval",
		"Energy.Active.Import.Interval",
		"Energy.Active.Net",
		"Energy.Reactive.Export.Interval",
		"Energy.Reactive.Import.Interval",
		"Energy.Reactive.Net",
		"Energy.Apparent.Net",
		"Energy.Apparent.Import",
		"Energy.Apparent.Export",
		"Frequncy",
		"Power.Active.Export",
		"Power.Active.Import",
		"Power.Factor",
		"Power.Offered",
		"Power.Reactive.Export",
		"Power.Reactive.Import",
		"SoC",
		"Voltage",
	}
	return contains(cases, enum)
}

func isMessageFormatEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "ASCII", "HTML", "URI", "UTF8":
		return true
	default:
		return false
	}
}


func isMessagePriorityEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "AlwaysFront", "InFront", "NormalCycle":
		return true
	default:
		return false
	}
}


func isMessageStateEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Charging", "Faulted", "Idle", "Unavailable":
		return true
	default:
		return false
	}
}


func isMessageTriggerEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"BootNotification",
		"LogStatusNotification",
		"FirmawareStatusNotification",
		"Heartbeat",
		"MeterValues",
		"SignChargingStationCertificate",
		"SignV2GCertificate",
		"StatusNotification",
		"TransactionEvent",
		"SignCombinedCertificate",
		"PublishFirmwareStatusNotification",
	}
	return contains(cases, enum)
}

func isMonitorEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"UpperThreshold",
		"LowerThreshold",
		"Delta",
		"Periodic",
		"PeriodicClockAligned",
	}
	return contains(cases, enum)
}

func isMonitoringBaseEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "All", "FactoryDefault", "HardWiredOnly":
		return true
	default:
		return false
	}
}


func isMonitoringCriterionEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "ThresholdMonitoring", "DeltaMonitoring", "PeriodicMonitoring":
		return true
	default:
		return false
	}
}


func isMutabilityEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "ReadOnly", "WriteOnly", "ReadWrite":
		return true
	default:
		return false
	}
}

func isNotifyEVChargingNeedsStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "Processing":
		return true
	default:
		return false
	}
}

func isOCPPInterfaceEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Wired0",
		"Wired1",
		"Wired2",
		"Wired3",
		"Wireless0",
		"Wireless1",
		"Wireless2",
		"Wireless3",
	}
	return contains(cases, enum)
}



func isOCPPTransportEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "JSON", "SOAP":
		return true
	default:
		return false
	}
}



func isOCPPVersionEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "OCPP12", "OCPP15", "OCPP16", "OCPP20":
		return true
	default:
		return false
	}
}


func isOperationalStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "InOperative", "Operative":
		return true
	default:
		return false
	}
}

func isPhaseEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
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
		"L3-L1",
	}
	return contains(cases, enum)
}

func isPublishFirmwareStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Idle",
		"DownloadScheduled",
		"Downloading",
		"Downloaded",
		"Published",
		"DownloadFailed",
		"DownloadPaused",
		"InvalidChecksum",
		"ChecksumVerified",
		"PublishFailed",
	}
	return contains(cases, enum)
}


func isReadingContextEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
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
	return contains(cases, enum)
}

func isReasonEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"DeAuthorized",
		"EmergencyStop",
		"EnergyLimitReached",
		"EVDisconnected",
		"GroundFault",
		"ImmediateReset",
		"Local",
		"LocalOutOfCredit",
		"MeterPass",
		"Other",
		"OvercurrentFault",
		"PowerLoss",
		"PowerQuality",
		"Reboot",
		"Remote",
		"SOCLimitReached",
		"StoppedByEV",
		"TimeLimitReached",
		"Timeout",
	}
	return contains(cases, enum)
}


func isRecurrencyKindEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Daily", "Weekly":
		return true
	default:
		return false
	}
}


func isRegistrationStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Pending", "Rejected":
		return true
	default:
		return false
	}
}


func isReportBaseEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "ConfigurationInventory", "FullInventory", "SummaryInventory":
		return true
	default:
		return false
	}
}


func isRequestStartStopStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected":
		return true
	default:
		return false
	}
}


func isReservationUpdateStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Expired", "Removed":
		return true
	default:
		return false
	}
}


func isReserveNowStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted",
		"Faulted",
		"Occupied",
		"Rejected",
		"Unavailable",
	}
	return contains(cases, enum)
}


func isResetEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Immediate", "OnIdle":
		return true
	default:
		return false
	}
}


func isResetStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "Scheduled":
		return true
	default:
		return false
	}
}


func isSendLocalListStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "VersionMismatch":
		return true
	default:
		return false
	}
}


func isSetMonitoringStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted",
		"UnkownComponent",
		"UnkownVariable",
		"UnsupportedMonitorType",
		"Rejected",
		"Duplicate",
	}
	return contains(cases, enum)
}


func isSetNetworkProfileStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "Failed":
		return true
	default:
		return false
	}
}


func isSetVariableStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted",
		"Rejected",
		"UnkownComponent",
		"UnkownVariable",
		"NotSupportedAttributeType",
		"RebootRequired",
	}
	return contains(cases, enum)
}


func isTransactionEventEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Ended", "Started", "Updated":
		return true
	default:
		return false
	}
}

func isTriggerMessageStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Accepted", "Rejected", "NotImplemented":
		return true
	default:
		return false
	}
}


func isTriggerReasonEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Authorized",
		"CablePluggedIn",
		"ChargingRateChanged",
		"ChargingStateChanged",
		"Deauthorized",
		"EnergyLimitReached",
		"EVCommunicationLost",
		"EVConnectTimeout",
		"MeterValueClock",
		"MeterValuePeriodic",
		"TimeLimitReached",
		"Trigger",
		"UnlockCommand",
		"StopAuthorized",
		"EVDeparted",
		"EVDetected",
		"RemoteStop",
		"RemoteStart",
		"AbnormalCondition",
		"SignedDataReceived",
		"ResetCommand",
	}
	return contains(cases, enum)
}



func isUnlockStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Unlocked", "UnlockFailed", "OngoingAuthorizedTransaction", "UnkwownConnector":
		return true
	default:
		return false
	}
}	

func isUnpublishFirmwareStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "DownloadOngoing", "NoFirmware", "Unpublished":
		return true
	default:
		return false
	}
}


func isUpdateEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "Differential", "Full":
		return true
	default:
		return false
	}
}

func isUpdateFirmwareStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"Accepted",
		"Rejected",
		"AcceptedCancelled",
		"InvalidCertificate",
		"RevokedCertificate",
	}
	return contains(cases, enum)
}


func isUploadLogStatusEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	cases := []string{
		"BadMessage",
		"Idle",
		"NotSupportedOperation",
		"PermissionDenied",
		"Uploaded",
		"UploadFailure",
		"Uploading",
		"AcceptedCancelled",
	}
	return contains(cases, enum)
}

func isVPNEnumType(fl validator.FieldLevel) bool {
	enum := fl.Field().String()
	switch enum {
	case "IKEv2", "IPSec", "L2TP", "PPTP":
		return true
	default:
		return false
	}
}