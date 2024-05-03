package inner

import (
	"github.com/LumiWave/baseapp/auth/inno"
	"github.com/LumiWave/inno-auth/rest_server/config"
)

func DecryptInnoUID(innoUID string) string {
	return inno.AESDecrypt(innoUID, []byte(config.GetInstance().Secret.Key), []byte(config.GetInstance().Secret.Iv))
}
