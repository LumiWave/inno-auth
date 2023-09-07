package context

///////////////////////////////////////////
// member auth log

const (
	state_normal      = 0
	state_maintenance = 1
	state_abnormal    = 2
)

type MaintenanceInfo struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`

	Content string `json:"content,omitempty"`
}

type StatusMain struct {
	IsMaintenance int             `json:"state"`
	Info          MaintenanceInfo `json:"info,omitempty"`
}

func NewStatusMain() *StatusMain {
	return new(StatusMain)
}

///////////////////////////////////////////
