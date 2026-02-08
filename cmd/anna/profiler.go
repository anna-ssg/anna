package anna

import (
	"runtime"
	"time"
)

func (cmd *Cmd) PrintStats(elapsedTime time.Duration) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	cmd.InfoLogger.Printf("Memory Usage: %d bytes", memStats.Alloc)
	cmd.InfoLogger.Printf("Time Elapsed: %s", elapsedTime)
	cpuUsage := runtime.NumCPU()
	threads := runtime.GOMAXPROCS(0)
	runtime.ReadMemStats(memStats)

	cmd.InfoLogger.Printf("Threads: %d", threads)
	cmd.InfoLogger.Printf("Cores: %d", cpuUsage)
	cmd.InfoLogger.Printf("Time Taken: %s", elapsedTime)
	cmd.InfoLogger.Printf("Allocated Memory: %d bytes", memStats.Alloc)
	cmd.InfoLogger.Printf("Total Memory Allocated: %d bytes", memStats.TotalAlloc)
	cmd.InfoLogger.Printf("Heap Memory In Use: %d bytes", memStats.HeapInuse)
	cmd.InfoLogger.Printf("Heap Memory Idle: %d bytes", memStats.HeapIdle)
	cmd.InfoLogger.Printf("Heap Memory Released: %d bytes", memStats.HeapReleased)
	cmd.InfoLogger.Printf("Number of Goroutines: %d", runtime.NumGoroutine())

	// Get the function with the highest CPU usage
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	cmd.InfoLogger.Printf("Function with Highest CPU Usage: %s", function.Name())
}
