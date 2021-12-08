package resultcode

const (
	Result_Success = 0

	Result_DBError        = 13000
	Result_DBNotExistItem = 13001
	Result_RedisError     = 13002

	Result_Auth_InvalidJwt = 20000

	Result_Auth_NotMatchAppAccount = 21001
	Result_Auth_EmptyAccessID      = 21002
	Result_Auth_EmptyAccessPW      = 21003

	Result_Auth_EmptyAccountSocialID      = 22001
	Result_Auth_EmptyAccountSocialType    = 22002
	Result_Procedure_Auth_Members         = 22003
	Result_Api_Post_Point_Member_Register = 22004
	Result_Api_Get_Token_Address_New      = 22005
	Result_Procedure_Add_Account_Coins    = 22006

	Result_Auth_MakeTokenError = 23001
)

var ResultCodeText = map[int]string{
	Result_Success: "success",

	Result_DBError:        "Internal DB error",
	Result_DBNotExistItem: "Not exist item",
	Result_RedisError:     "Redis Error",

	Result_Auth_InvalidJwt: "Invalid jwt token",

	Result_Auth_NotMatchAppAccount: "Account information does not match",
	Result_Auth_EmptyAccessID:      "Empty Access ID",
	Result_Auth_EmptyAccessPW:      "Empty Access PW",

	Result_Auth_EmptyAccountSocialID:      "Empty Account SocialID",
	Result_Auth_EmptyAccountSocialType:    "Empty Account SocialType",
	Result_Procedure_Auth_Members:         "Error Procedure Auth Members",
	Result_Api_Post_Point_Member_Register: "Error API Post Point Member Register",
	Result_Api_Get_Token_Address_New:      "Error API Get Token Address New",
	Result_Procedure_Add_Account_Coins:    "Error Procedure Add Accounts Coins",

	Result_Auth_MakeTokenError: "Make Token Error.",
}
