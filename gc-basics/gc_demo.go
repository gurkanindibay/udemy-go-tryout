package main

import (
    "fmt"
    "runtime"
    "time"
)

func printMem(prefix string) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("%s: Alloc=%d MB	HeapSys=%d MB	NumGC=%d\n",
        prefix, m.Alloc/1024/1024, m.HeapSys/1024/1024, m.NumGC)
}

func main() {
    fmt.Println("Go GC demo: allocations will trigger garbage collections (see GODEBUG=gctrace=1)")
    printMem("initial")

    // Allocate several 10MB slices to grow the heap
    holders := make([][]byte, 0, 200)
    for i := 0; i < 60; i++ {
        holders = append(holders, make([]byte, 10*1024*1024)) // 10 MB
        printMem(fmt.Sprintf("after alloc %02d", i+1))
        time.Sleep(150 * time.Millisecond)
        if (i+1)%15 == 0 {
            fmt.Println("-> forcing runtime.GC()")
            runtime.GC()
            time.Sleep(100 * time.Millisecond)
            printMem(fmt.Sprintf("post runtime.GC() %02d", i+1))
        }
    }

    // Drop references and force a GC to reclaim memory
    holders = nil
    fmt.Println("Dropped references; forcing final GC")
    runtime.GC()
    time.Sleep(200 * time.Millisecond)
    printMem("after drop and GC")

    // small pause so GODEBUG trace lines (if any) are printed before exit
    time.Sleep(300 * time.Millisecond)
}
