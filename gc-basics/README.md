# Go Garbage Collector (GC) — guide and demo

This folder contains a small demonstration program `gc_demo.go` that allocates memory and shows how Go's garbage collector behaves. This README documents the key GC concepts, how the demo maps to them, how to reproduce the traces, tuning knobs, common pitfalls, and next steps for investigation.

## What this explains
- High-level model of Go's GC (concurrent, non-generational, tri-color mark-and-sweep)
- Key mechanisms: write barrier, pacer, mutator assist, STW pauses, sweeper, finalizers
- How to observe GC with `GODEBUG` and `pprof`
- How `gc_demo.go` demonstrates the behavior
- Practical knobs (`GOGC`, `debug.SetGCPercent`, `runtime.GC()`), tips and caveats

## Quick reproduction
Open a bash shell and run from this directory:

```bash
cd /c/source/Tryouts/golang/gc-basics

# Run demo normally
go run gc_demo.go

# Run demo with GC tracing and pacer tracing (shows GC events and pacer decisions)
GODEBUG=gctrace=1,gcpacertrace=1 GOGC=100 go run gc_demo.go
```

Notes:
- `GODEBUG=gctrace=1` prints summary GC traces (`gc N @...` lines).
- `gcpacertrace=1` prints pacer decisions and mutator-assist ratios.
- `GOGC` controls aggressiveness (default `100` → allow heap to grow ~2× live). Lower means more frequent GC and smaller heap.

## High-level model
- Go's GC is part of the runtime and is linked into each Go binary.
- It is a concurrent, stop-the-world (brief) aware, tri-color mark-and-sweep collector.
- Not generational: there is no young/old heap split in the standard GC.
- The runtime attempts to work concurrently with the program (mutator) and keeps GC pause times short.

## GC phases (simplified)
1. Trigger: when heap growth reaches the pacer goal (based on live heap and `GOGC`).
2. Concurrent mark: background workers + mutator assistance traverse reachable objects using a tri-color invariant.
3. Mark termination: short stop-the-world (STW) pause to finish marking and prepare to sweep.
4. Sweep: reclaimed memory is returned to runtime free lists; sweeping happens in background.
5. Finalizers: if set, finalizers run in separate goroutines and delay actual reclamation until they complete.

## Important mechanisms explained
- Write barrier: executed on pointer writes; ensures mutator changes won't break concurrent marking (keeps tri-color invariant).
- Mutator assist: when allocation outpaces background marking, mutators help finish mark work to keep pacing goals.
- Pacer: calculates how much marking work must be done and schedules background workers, aiming to keep heap near the goal.
- STW pauses: unavoidable short pauses (stack scan, mark termination); typically small but dependent on workload and heap size.

## Interpreting `GODEBUG` trace lines
A typical `gctrace` line looks like:

```
gc 1 @0.014s 6%: 0.84+0.50+0.52 ms clock, 18+0/2.0/0+11 ms cpu, 3->4->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 22 P
```

- `gc 1` — collection number.
- `@0.014s` — time since program start when GC started.
- `6%` — approximate CPU percentage GC used.
- `0.84+0.50+0.52 ms clock` — STW (initial) + concurrent + STW (final) in wall-clock ms.
- `18+0/2.0/0+11 ms cpu` — CPU-time breakdown for subphases.
- `3->4->1 MB` — heap sizes at key points (before mark -> at mark termination -> after sweep or similar).
- `4 MB goal` — heap goal used by the pacer.

`gcpacertrace` adds details about the pacer’s calculations and the assist ratio (how much mutators must help).

## How `gc_demo.go` maps to concepts
- The demo performs repeated 10 MB allocations to grow the heap, prints `runtime.ReadMemStats()` snapshots, and calls `runtime.GC()` periodically and at the end.
- This shows:
  - When the runtime triggers concurrent GC as the heap grows.
  - How `NumGC` increases and `Alloc`/`HeapSys` change.
  - The difference between `Alloc` (live bytes) and `HeapSys` (heap memory obtained from OS/runtime arenas).
  - Effects of forced GC (`runtime.GC()`) which show `(forced)` in trace lines.

## Practical knobs and commands
- Observe GC traces: `GODEBUG=gctrace=1`.
- See pacer details: `GODEBUG=gctrace=1,gcpacertrace=1`.
- Change aggressiveness: `GOGC=50 go run ...` or at runtime with `runtime/debug.SetGCPercent(n)`.
- Force a GC (for tests/demos): `runtime.GC()` — not recommended in production generally.
- Disable automatic GC (only for controlled tests): `debug.SetGCPercent(-1)` (be careful — can cause OOM).
- Capture heap profiles: start an HTTP pprof server (`import _ "net/http/pprof"`) or use `runtime/pprof` to write profiles and inspect with `go tool pprof`.

Example: capture heap profile programmatically (simple sketch)

```go
// import "runtime/pprof" and create a file to write the profile
f, _ := os.Create("heap.prof")
pprof.WriteHeapProfile(f)
f.Close()
```

Then analyze:

```bash
go tool pprof -http=:8080 path/to/binary heap.prof
```

## Tuning and best practices
- Use default GC for most apps — it's conservative and low-maintenance.
- Lower `GOGC` to reduce resident memory at the cost of more CPU spent in GC.
- Increase `GOGC` to reduce GC CPU but accept larger memory usage.
- Avoid excessive allocations in hot code paths (object pooling, reuse slices, sync.Pool where appropriate).
- Use profiling (`pprof`) to find actual allocation hot spots before micro-tuning.

## Edge cases and caveats
- cgo / C allocations: memory allocated in C is not tracked by Go GC. If you allocate large amounts in C, Go's GC won't see it; you must manage it manually.
- Pinned memory: passing Go pointers to C or runtime pinning can prevent the GC from moving/reclaiming objects and increase retention.
- Finalizers: finalizers delay reclamation and can cause unexpected retention if they allocate or are slow.
- `plugin` / shared libs and GC: using plugins and shared objects has version/platform constraints and may complicate memory behavior.
- Very large heaps: GC CPU overhead grows with heap size (concurrent marking work increases), so tuning and architecture choices matter for very large memory apps.

## Debugging tips
- Start with `GODEBUG=gctrace=1` to see whether GC frequency or pause times are concerns.
- Use `gcpacertrace=1` if you need to see why the GC pacer is asking for more/less work.
- Capture heap profiles and flamegraphs with `go tool pprof` and visual tools (pprof web UI, speedscope, etc.).
- Reproduce with controlled allocations (like `gc_demo.go`) to reason about pacing and `GOGC` behavior.

## Contract (what this demo/document shows)
- Inputs: a Go program that allocates memory; runtime environment variables `GOGC` and `GODEBUG`.
- Outputs: GC traces (via `GODEBUG`), heap stats (via `runtime.ReadMemStats()`), and profile output (via `pprof`).
- Error modes: high allocation rates can cause high CPU usage for GC or OOM if the OS cannot satisfy allocations. Disabling GC can lead to OOM.

## Next steps / suggested experiments
1. Run `GODEBUG=gctrace=1,gcpacertrace=1 GOGC=50 go run gc_demo.go` and compare traces.
2. Add `net/http/pprof` to a small server and capture heap/alloc profiles under real workload.
3. Replace `make([]byte, 10*1024*1024)` with smaller allocations to see how allocation granularity affects pacing and sweep behavior.
4. Experiment with `debug.SetGCPercent()` in code to adjust GC at runtime and measure memory/CPU impact.

---
If you'd like, I can also add a tiny `pprof` example in this folder, collect a heap profile from `gc_demo.go`, and add instructions for interpreting the pprof results. What would you like me to add next?
