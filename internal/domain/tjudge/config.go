package tjudge

import "time"

type Config struct {
	Title        string
	AllowLateReg bool
	Maintenance  bool
	JudgeTimeout time.Time
	FileLimit    int // TODO: kb?
}
