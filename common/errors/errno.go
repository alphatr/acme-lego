package errors

/**
 * 错误代码列表
 * 错误码: ABBCCDDD
 * A 错误分组 main: 1; common: 2; controller: 3; model: 4; unknow: 9;
 * BB 错误模块 common: 0; config: 1; bootstrap 2;
 * CC 错误子模块
 * DDD 错误功能
 */
const (
	MainInitErrno                  ErrorNum = 10000001
	CommonFileNotExistErrno        ErrorNum = 20001001
	CommonFileCreateErrno          ErrorNum = 20001002
	CommonFileReadErrno            ErrorNum = 20001003
	CommonFileWriteErrno           ErrorNum = 20001004
	CommonFileCloseErrno           ErrorNum = 20001005
	CommonPathAbsErrno             ErrorNum = 20001006
	CommonMakeDirErrno             ErrorNum = 20001007
	CommonFileIsExistErrno         ErrorNum = 20001008
	CommonTOMLUnmarshalErrno       ErrorNum = 20002001
	CommonJSONUnmarshalErrno       ErrorNum = 20003001
	CommonJSONMarshalErrno         ErrorNum = 20003002
	CommonCommandRunErrno          ErrorNum = 20004001
	CommonParsePrivateErrno        ErrorNum = 20005001
	CommonParseCertificateErrno    ErrorNum = 20005002
	CommonUnknowBlockErrno         ErrorNum = 20005003
	CommonMarshalPrivateErrno      ErrorNum = 20005004
	CommonParseHostPortErrno       ErrorNum = 20006001
	ConfigInitErrno                ErrorNum = 20101001
	ConfigParseTOMLErrno           ErrorNum = 20101002
	ConfigBaseInitErrno            ErrorNum = 20102001
	ConfigDomainInitErrno          ErrorNum = 20103001
	BootstrapInitErrno             ErrorNum = 20201001
	BootstrapInitLoggerErrno       ErrorNum = 20202001
	BootstrapInitHandlerErrno      ErrorNum = 20203001
	ConGetAccountErrno             ErrorNum = 30101001
	ConInitClientErrno             ErrorNum = 30101002
	ConRequireParamErrno           ErrorNum = 30101003
	ConErrorParamErrno             ErrorNum = 30101004
	ConAccCreateErrno              ErrorNum = 30201001
	ConAccRegisterErrno            ErrorNum = 30201002
	ConAccSaveErrno                ErrorNum = 30201003
	ConCertSetupChallengeErrno     ErrorNum = 30301001
	ConCertGenerateKeyErrno        ErrorNum = 30301002
	ConCertObtainErrno             ErrorNum = 30301003
	ConCertObtainDomainErrno       ErrorNum = 30301004
	ConCertRenewDomainErrno        ErrorNum = 30301005
	ConCertCheckFolderErrno        ErrorNum = 30301006
	ConCertSaveCertErrno           ErrorNum = 30301007
	ConCertRunAfterRenewErrno      ErrorNum = 30301008
	ConCertLoadPrivateErrno        ErrorNum = 30301009
	ConCertRenewIgnoreErrno        ErrorNum = 30301010
	ModelClientInitErrno           ErrorNum = 40101001
	ModelClientRegisterErrno       ErrorNum = 40101002
	ModelClientObtainErrno         ErrorNum = 40101003
	ModelClientUnknowProviderErrno ErrorNum = 40101101
	ModelClientProviderErrno       ErrorNum = 40101102
	ModelClientSetProviderErrno    ErrorNum = 40101103
	ModelClientTypeProviderErrno   ErrorNum = 40101104
	ModelAccSaveConfigErrno        ErrorNum = 40201001
	ModelAccSavePrivateErrno       ErrorNum = 40201002
	ModelAccLoadPrivateErrno       ErrorNum = 40201003
	ModelAccNotECDSAErrno          ErrorNum = 40201004
	ModelAccGenerateKeyErrno       ErrorNum = 40201005
	ModelChalHTTPInitErrno         ErrorNum = 40301001
	ModelChalServerStartErrno      ErrorNum = 40301002
	ModelChalDNSConfigErrno        ErrorNum = 40301003
	UnknowErrno                    ErrorNum = 90000000
)

// ErrorMap 错误 Map 列表
var ErrorMap = map[ErrorNum]ErrorContent{
	MainInitErrno:                  {"init-error", 0},
	CommonFileNotExistErrno:        {"file-not-exist(%s)", 0},
	CommonFileCreateErrno:          {"create-file(%s)", 0},
	CommonFileReadErrno:            {"read-file(%s)", 0},
	CommonFileWriteErrno:           {"write-file(%s)", 0},
	CommonFileCloseErrno:           {"close-file(%s)", 0},
	CommonPathAbsErrno:             {"get-absolute-path(%s)", 0},
	CommonMakeDirErrno:             {"mkdir(%s)", 0},
	CommonFileIsExistErrno:         {"file-is-exist(%s)", 0},
	CommonTOMLUnmarshalErrno:       {"toml-decode", 0},
	CommonJSONUnmarshalErrno:       {"json-unmarshal", 0},
	CommonJSONMarshalErrno:         {"json-marshal", 0},
	CommonCommandRunErrno:          {"run-command(%s)", 0},
	CommonParsePrivateErrno:        {"parse-private-key(%s)", 0},
	CommonParseCertificateErrno:    {"parse-certificate(%s)", 0},
	CommonUnknowBlockErrno:         {"unknow-pem-block(%s)", 0},
	CommonMarshalPrivateErrno:      {"marshal-private-key(%s)", 0},
	CommonParseHostPortErrno:       {"parse-host-port", 0},
	ConfigInitErrno:                {"init-config", 0},
	ConfigParseTOMLErrno:           {"parse-toml-config", 0},
	ConfigBaseInitErrno:            {"init-base-config", 0},
	ConfigDomainInitErrno:          {"init-domain-config", 0},
	BootstrapInitErrno:             {"init-bootstrap", 0},
	BootstrapInitLoggerErrno:       {"init-logger", 0},
	BootstrapInitHandlerErrno:      {"bootstrap-init-handle(%s)", 0},
	ConGetAccountErrno:             {"get-account", 0},
	ConInitClientErrno:             {"init-client", 0},
	ConRequireParamErrno:           {"require-param(%s)", 0},
	ConErrorParamErrno:             {"error-param(%s)", 0},
	ConAccCreateErrno:              {"create-account", 0},
	ConAccRegisterErrno:            {"account-register", 0},
	ConAccSaveErrno:                {"save-account", 0},
	ConCertSetupChallengeErrno:     {"setup-challenge", 0},
	ConCertGenerateKeyErrno:        {"generate-private-key", 0},
	ConCertObtainErrno:             {"obtain-certificate(%s, %s)", 0},
	ConCertObtainDomainErrno:       {"obtain-domain-certificate(%s)", 0},
	ConCertRenewDomainErrno:        {"renew-domain-certificate(%s)", 0},
	ConCertCheckFolderErrno:        {"check-folder(%s)", 0},
	ConCertSaveCertErrno:           {"save-certificate(%s, %s)", 0},
	ConCertRunAfterRenewErrno:      {"run-after-renew", 0},
	ConCertLoadPrivateErrno:        {"load-private-key", 0},
	ConCertRenewIgnoreErrno:        {"renew-ignore", 0},
	ModelClientInitErrno:           {"init-client", 0},
	ModelClientRegisterErrno:       {"register-account", 0},
	ModelClientObtainErrno:         {"obtain-certificate", 0},
	ModelClientUnknowProviderErrno: {"unknow-provider(%s)", 0},
	ModelClientProviderErrno:       {"provider-server", 0},
	ModelClientSetProviderErrno:    {"client-set-provider", 0},
	ModelClientTypeProviderErrno:   {"client-set-type-provider(%s)", 0},
	ModelAccSaveConfigErrno:        {"save-account-config", 0},
	ModelAccSavePrivateErrno:       {"save-account-private-key", 0},
	ModelAccLoadPrivateErrno:       {"load-account-private-key", 0},
	ModelAccNotECDSAErrno:          {"not-ecdsa-private-key", 0},
	ModelAccGenerateKeyErrno:       {"generate-private-key", 0},
	ModelChalHTTPInitErrno:         {"init-http-provider", 0},
	ModelChalServerStartErrno:      {"server-start", 0},
	ModelChalDNSConfigErrno:        {"init-dns-config(%s)", 0},
	UnknowErrno:                    {"unknow-error %s", 0},
}
