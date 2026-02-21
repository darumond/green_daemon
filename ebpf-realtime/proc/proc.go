package proc

import (
	"github.com/shirou/gopsutil/v3/process"
)

type Pid uint32

type ProcNameMap struct {
	pidToName map[Pid]string
}

func NewProcNameMap() ProcNameMap {
	return ProcNameMap{
		pidToName: make(map[Pid]string),
	}
}

func (p *ProcNameMap) GetName(pid Pid) (string, error) {
	if name, ok := p.pidToName[pid]; ok {
		return name, nil
	}
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return "unknown", err
	}
	return proc.Name()
}
