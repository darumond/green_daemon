package procinfo

import (
	"github.com/shirou/gopsutil/v3/process"
)

type SingleProcInfo struct {
	proc *process.Process
}

func (p *SingleProcInfo) GetTotalMem() (float64, error) {
	info, err := p.proc.MemoryInfo()
	if err != nil {
		return 0.0, err
	}
	return BytesToMiB(info.RSS), nil
}

func (p *SingleProcInfo) GetCpuTime() (float64, error) {
	ts, err := p.proc.Times()
	if err != nil {
		return 0, err
	}
	return ts.User + ts.System, nil
}

func (p *SingleProcInfo) Name() (string, error) {
	return p.proc.Name()
}

// ProcInfo represents a snapshot of all processes on the system.
type ProcInfo struct {
	procs []SingleProcInfo
}

// GetCurrProcInfo returns the info on all current processes.
func GetCurrProcInfo() (ProcInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return ProcInfo{}, err
	}
	procInfos := []SingleProcInfo{}
	for _, proc := range procs {
		procInfos = append(procInfos, SingleProcInfo{proc})
	}
	return ProcInfo{
		procs: procInfos,
	}, nil
}
