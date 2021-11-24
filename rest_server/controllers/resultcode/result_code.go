package resultcode

const (
	Result_Success = 0

	Result_DBError        = 13000
	Result_DBNotExistItem = 13001
	Result_RedisError     = 13002

	Result_Auth_InvalidJwt     = 20000
	Result_Auth_EmptyCpName    = 20001
	Result_Auth_ExistsCpName   = 20002
	Result_Auth_NotFoundCpName = 20003
	Result_Auth_NotFoundCpIdx  = 20004

	Result_Auth_EmptyAppName    = 21001
	Result_Auth_ExistsAppName   = 21002
	Result_Auth_NotFoundAppName = 21003

	Result_Auth_EmptyMemberSocialInfo  = 22001
	Result_Auth_ExistsMemberSocialInfo = 22002
)

var ResultCodeText = map[int]string{
	Result_Success: "success",

	Result_DBError:        "Internal DB error",
	Result_DBNotExistItem: "Not exist item",
	Result_RedisError:     "Redis Error",

	Result_Auth_InvalidJwt:     "Invalid jwt token",
	Result_Auth_EmptyCpName:    "Empty CP Name",
	Result_Auth_ExistsCpName:   "Exists CP Name",
	Result_Auth_NotFoundCpName: "Not Found CP Name",
	Result_Auth_NotFoundCpIdx:  "Not Found CP Idx",

	Result_Auth_EmptyAppName:    "Empty App Name",
	Result_Auth_ExistsAppName:   "Exists App Name",
	Result_Auth_NotFoundAppName: "Not Found App Name",

	Result_Auth_EmptyMemberSocialInfo:  "Empty Member Social Info",
	Result_Auth_ExistsMemberSocialInfo: "Exists Member Social Info",
}
