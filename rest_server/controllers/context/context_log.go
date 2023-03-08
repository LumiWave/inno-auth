package context

const (
	AccountAuthLog_Auth       = 4 // INNO Web 플랫폼 인증 LogID
	AccountAuthLog_NewAccount = 5 // INNO 플랫폼 신규 계정 Web 로그인 EventID
	AccountAuthLog_Account    = 6 // INNO 플랫폼 기존 계정 Web 로그인 EventID

	MemberAuthLog_Auth      = 4 // INNO App 플랫폼 인증 logID
	MemberAuthLog_NewMember = 7 // 신규 게임 계정 생성 EventID
	MemberAuthLog_Member    = 8 // 기존 게임 계정 로그인 EventID
)
