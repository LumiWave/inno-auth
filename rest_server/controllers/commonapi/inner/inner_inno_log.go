package inner

import (
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/inno_log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/log_server"
)

func PostAccountAuthLog(params *inno_log.AccountAuthLog, IsJoined bool) {
	if IsJoined {
		// 3-1. [DB] 신규 사용자 로그 등록
		params.EventID = context.AccountAuthLog_NewAccount
	} else {
		// 3-1. [DB] 기존 사용자 로그 등록
		params.EventID = context.AccountAuthLog_Account
	}
	go log_server.GetInstance().PostAccountAuth(params)
}
