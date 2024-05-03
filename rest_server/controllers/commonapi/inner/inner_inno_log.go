package inner

import (
	"github.com/LumiWave/baseInnoClient/inno_log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/controllers/log_server"
)

func PostAccountAuthLog(params *inno_log.AccountAuthLog, IsJoined bool) {
	if IsJoined {
		// INNO 플랫폼 신규 계정 Web 로그인 EventID
		params.EventID = context.AccountAuthLog_NewAccount
	} else {
		// INNO 플랫폼 기존 계정 Web 로그인 EventID
		params.EventID = context.AccountAuthLog_Account
	}
	go log_server.GetInstance().PostAccountAuth(params)
}

func PostMemberAuthLog(params *inno_log.MemberAuthLog, IsJoined bool) {
	if IsJoined {
		// 신규 게임 계정 생성 EventID
		params.EventID = context.MemberAuthLog_NewMember
	} else {
		// 기존 게임 계정 로그인 EventID
		params.EventID = context.MemberAuthLog_Member
	}
	go log_server.GetInstance().PostMemberAuth(params)
}
