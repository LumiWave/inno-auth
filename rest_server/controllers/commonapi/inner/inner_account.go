package inner

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

func DecryptInnoUID(innoUID string) string {
	return inno.AESDecrypt(innoUID, []byte(config.GetInstance().Secret.Key), []byte(config.GetInstance().Secret.Iv))
}
