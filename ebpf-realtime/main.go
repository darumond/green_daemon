package main

import (
	"bytes"
	"context"
	"ebpf-realtime/proc"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" bpf bpf/tracer.bpf.c -- -I./bpf

// Event struct mirroring our C struct
type Event struct {
	Type       uint8
	_          [3]byte // Padding
	Pid        uint32
	Cpu        uint32
	_          [4]byte // Padding
	Ts         uint64
	DurationNs uint64
	Packetsize uint64
	NextPid    uint32
	_          [4]byte
}

const (
	EventTypeTcpSend     = 1
	EventTypeSchedSwitch = 2
)

// Data structures for storage
type SchedRecord struct {
	Ts       uint64
	Cpu      uint32
	PrevPid  uint32
	NextPid  uint32
	PrevComm string
	NextComm string
}

type TcpRecord struct {
	Ts         uint64
	DurationNs uint64
	Size       uint64
}

func main() {
	pidFlag := flag.Int("pid", 0, "Target PID to monitor for TCP latency (Mandatory)")
	flag.Parse()

	if *pidFlag <= 0 {
		fmt.Println("Error: --pid is required and must be greater than 0")
		flag.Usage()
		os.Exit(1)
	}
	targetPid := uint32(*pidFlag)

	// 2. Load eBPF objects
	spec, err := loadBpf()
	if err != nil {
		log.Fatalf("Failed to load eBPF spec: %v", err)
	}

	if err := spec.RewriteConstants(map[string]interface{}{
		"target_pid": targetPid,
	}); err != nil {
		log.Fatalf("Failed to rewrite constants: %v", err)
	}

	var objs bpfObjects
	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("Failed to load objects: %v", err)
	}
	defer objs.Close()

	// 3. Attach Probes
	kp, err := link.Kprobe("tcp_sendmsg", objs.KprobeTcpSendmsg, nil)
	if err != nil {
		log.Fatalf("Failed to attach kprobe: %v", err)
	}
	defer kp.Close()

	krp, err := link.Kretprobe("tcp_sendmsg", objs.KretprobeTcpSendmsg, nil)
	if err != nil {
		log.Fatalf("Failed to attach kretprobe: %v", err)
	}
	defer krp.Close()

	tp, err := link.AttachRawTracepoint(link.RawTracepointOptions{
		Name:    "sched_switch",
		Program: objs.RawTpSchedSwitch,
	})
	if err != nil {
		log.Fatalf("Failed to attach raw tracepoint: %v", err)
	}
	defer tp.Close()

	rd, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("Failed to open ringbuf reader: %v", err)
	}
	defer rd.Close()

	var schedEvents []SchedRecord
	var tcpEvents []TcpRecord

	log.Printf("Profiling PID %d. Press Ctrl-C to stop and save data...", targetPid)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		rd.Close()
	}()

	pidToName := proc.NewProcNameMap()

	var event Event
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				break
			}
			log.Printf("Read error: %v", err)
			continue
		}

		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("Decode error: %v", err)
			continue
		}

		switch event.Type {
		case EventTypeTcpSend:
			tcpEvents = append(tcpEvents, TcpRecord{
				Ts:         event.Ts,
				DurationNs: event.DurationNs,
				Size:       event.Packetsize,
			})
		case EventTypeSchedSwitch:
			prevName, _ := pidToName.GetName(proc.Pid(event.Pid))
			nextName, _ := pidToName.GetName(proc.Pid(event.NextPid))
			schedEvents = append(schedEvents, SchedRecord{
				Ts:       event.Ts,
				Cpu:      event.Cpu,
				PrevPid:  event.Pid,
				NextPid:  event.NextPid,
				PrevComm: prevName,
				NextComm: nextName,
			})
		}
	}

	printSummary(targetPid, len(tcpEvents), len(schedEvents))

	schedEvents = filterSchedEvents(schedEvents, tcpEvents)

	if err := saveTcpToCsv("tcp_latency.csv", tcpEvents); err != nil {
		log.Printf("Failed to save TCP CSV: %v", err)
	}
	if err := saveSchedToCsv("scheduling_events.csv", schedEvents); err != nil {
		log.Printf("Failed to save Sched CSV: %v", err)
	}
}

func printSummary(pid uint32, tcpCount, schedCount int) {
	fmt.Println("\n--- Profiling Summary ---")
	fmt.Printf("Target PID:             %d\n", pid)
	fmt.Printf("Total TCP Sends:        %d\n", tcpCount)
	fmt.Printf("Total Sched Switches:   %d\n", schedCount)
	fmt.Println("------------------------")
}

func saveTcpToCsv(filename string, data []TcpRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"timestamp_ns", "duration_ns", "size_bytes"})
	for _, r := range data {
		writer.Write([]string{
			strconv.FormatUint(r.Ts, 10),
			strconv.FormatUint(r.DurationNs, 10),
			strconv.FormatUint(r.Size, 10),
		})
	}
	log.Printf("Saved %d TCP records to %s", len(data), filename)
	return nil
}

func saveSchedToCsv(filename string, data []SchedRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"timestamp_ns", "cpu", "prev_pid", "next_pid", "prev_comm", "next_comm"})
	for _, r := range data {
		writer.Write([]string{
			strconv.FormatUint(r.Ts, 10),
			strconv.FormatUint(uint64(r.Cpu), 10),
			strconv.FormatUint(uint64(r.PrevPid), 10),
			strconv.FormatUint(uint64(r.NextPid), 10),
			r.PrevComm,
			r.NextComm,
		})
	}
	log.Printf("Saved %d Scheduling records to %s", len(data), filename)
	return nil
}

// filterSchedEvents prunes the scheduling events to only those that occurred
// during the lifetime of any recorded TCP send operations.
func filterSchedEvents(schedEvents []SchedRecord, tcpEvents []TcpRecord) []SchedRecord {
	// 1. Sort both slices by timestamp to account for multi-CPU ring buffer jitter
	sort.Slice(schedEvents, func(i, j int) bool {
		return schedEvents[i].Ts < schedEvents[j].Ts
	})
	sort.Slice(tcpEvents, func(i, j int) bool {
		return tcpEvents[i].Ts < tcpEvents[j].Ts
	})

	var filtered []SchedRecord
	lastAddedIdx := -1 // Track this to prevent duplicates if TCP windows overlap

	for _, tcp := range tcpEvents {
		startTime := tcp.Ts - tcp.DurationNs
		endTime := tcp.Ts

		// 2. Binary search to find the first scheduler event in this window
		idx := sort.Search(len(schedEvents), func(i int) bool {
			return schedEvents[i].Ts >= startTime
		})

		// Fast-forward if we've already added this index from a previous overlapping TCP event
		if idx <= lastAddedIdx {
			idx = lastAddedIdx + 1
		}

		// 3. Collect all events that fall within the current TCP window
		for i := idx; i < len(schedEvents); i++ {
			if schedEvents[i].Ts > endTime {
				break // We are past the window, stop looking
			}
			filtered = append(filtered, schedEvents[i])
			lastAddedIdx = i
		}
	}

	return filtered
}
