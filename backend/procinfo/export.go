package procinfo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	memUtlization = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_utlization",
			Help: "The memory utilization of a process, labeled by name and pid",
		},
		[]string{"name", "pid"},
	)
	cpuTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_time_total",
			Help: "The CPU time of a process, labeled by name and pid",
		},
		[]string{"name", "pid"},
	)
)

func PollLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for t := range ticker.C {
		fmt.Println("Tick at", t.Format("15:04:05"))
		pollAndUpdate()
	}
}

func pollAndUpdate() error {
	infos, err := GetCurrProcInfo()
	if err != nil {
		return err
	}

	for _, info := range infos.procs {
		// Get the labels for the process.
		name, err := info.Name()
		if err != nil {
			continue
		}
		pid := strconv.Itoa(int(info.proc.Pid))
		labels := []string{name, pid}

		cpu, err := info.GetCpuTime()
		if err != nil {
			continue
		}
		mem, err := info.GetTotalMem()
		if err != nil {
			continue
		}

		memUtlization.WithLabelValues(labels...).Set(mem)
		cpuTime.WithLabelValues(labels...).Set(cpu)
	}

	return nil
}
