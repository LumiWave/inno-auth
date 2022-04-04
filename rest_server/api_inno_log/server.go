package api_inno_log

var gServerInfo *ServerInfo

type HostInfo struct {
	IntHostUri string
	ExtHostUri string
	IntVer     string // m1.0
	ExtVer     string // v1.0
}

type AuthInfo struct {
	ApiKey string
}

type ServerInfo struct {
	HostInfo

	AuthInfo
}

func GetInstance() *ServerInfo {
	return gServerInfo
}

func NewServerInfo(apiKey string, hostInfo HostInfo) *ServerInfo {
	if gServerInfo == nil {
		gServerInfo = &ServerInfo{
			HostInfo: hostInfo,
			AuthInfo: AuthInfo{
				ApiKey: apiKey,
			},
		}
	}

	return gServerInfo
}

func (o *ServerInfo) SetApiKey(key string) {
	gServerInfo.ApiKey = key
}
