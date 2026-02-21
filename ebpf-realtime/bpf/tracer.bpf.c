// go:build ignore

#define __TARGET_ARCH_x86

#include "vmlinux.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

char __license[] SEC("license") = "Dual MIT/GPL";

// The PID we want to scope the TCP tracing to.
// Go will override this before loading the program.
const volatile u32 target_pid = 0;

// Event Types
#define EVENT_TYPE_TCP_SEND 1
#define EVENT_TYPE_SCHED_SWITCH 2

// The structure we will send to Go via the Ring Buffer
struct event_t {
  u8 type;
  u32 pid;
  u32 cpu;
  u64 ts; // Timestamp of the event

  // Fields used for TCP send
  u64 duration_ns;

  // Fields used for Sched Switch
  u32 next_pid;
  char comm[16];      // Prev task name
  char next_comm[16]; // Next task name
};

// Ring buffer to push events to userspace
struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 256 * 1024);
} events SEC(".maps");

// Hash map to store the start time of tcp_sendmsg
struct {
  __uint(type, BPF_MAP_TYPE_HASH);
  __uint(max_entries, 10240);
  __type(key, u64);   // thread id (tgid/pid)
  __type(value, u64); // start timestamp
} tcp_start_times SEC(".maps");

// --- 1. TCP Sendmsg Probes ---

SEC("kprobe/tcp_sendmsg")
int BPF_KPROBE(kprobe_tcp_sendmsg) {
  u64 id = bpf_get_current_pid_tgid();
  u32 pid = id >> 32;

  // Filter by our target process
  if (target_pid != 0 && pid != target_pid) {
    return 0;
  }

  u64 ts = bpf_ktime_get_ns();
  bpf_map_update_elem(&tcp_start_times, &id, &ts, BPF_ANY);
  return 0;
}

SEC("kretprobe/tcp_sendmsg")
int BPF_KRETPROBE(kretprobe_tcp_sendmsg, int ret) {
  u64 id = bpf_get_current_pid_tgid();
  u32 pid = id >> 32;

  if (target_pid != 0 && pid != target_pid) {
    return 0;
  }

  u64 *start_ts = bpf_map_lookup_elem(&tcp_start_times, &id);
  if (!start_ts) {
    return 0; // Missed the start event
  }

  struct event_t *event =
      bpf_ringbuf_reserve(&events, sizeof(struct event_t), 0);
  if (!event) {
    return 0; // Ring buffer full
  }

  event->type = EVENT_TYPE_TCP_SEND;
  event->pid = pid;
  event->cpu = bpf_get_smp_processor_id();
  event->ts = bpf_ktime_get_ns();
  event->duration_ns = event->ts - *start_ts;

  bpf_get_current_comm(&event->comm, sizeof(event->comm));

  bpf_ringbuf_submit(event, 0);
  bpf_map_delete_elem(&tcp_start_times, &id);

  return 0;
}

// --- 2. Sched Switch Tracepoint ---

SEC("raw_tracepoint/sched_switch")
int raw_tp_sched_switch(struct bpf_raw_tracepoint_args *ctx) {
  // ctx->args[0] is prev task_struct, ctx->args[1] is next task_struct
  struct task_struct *prev = (struct task_struct *)ctx->args[0];
  struct task_struct *next = (struct task_struct *)ctx->args[1];

  struct event_t *event =
      bpf_ringbuf_reserve(&events, sizeof(struct event_t), 0);
  if (!event) {
    return 0;
  }

  event->type = EVENT_TYPE_SCHED_SWITCH;
  event->ts = bpf_ktime_get_ns();
  event->cpu = bpf_get_smp_processor_id();

  // Read PIDs
  event->pid = BPF_CORE_READ(prev, tgid);
  event->next_pid = BPF_CORE_READ(next, tgid);

  // Read Command names
  bpf_core_read_str(&event->comm, sizeof(event->comm), &prev->comm);
  bpf_core_read_str(&event->next_comm, sizeof(event->next_comm), &next->comm);

  bpf_ringbuf_submit(event, 0);
  return 0;
}
