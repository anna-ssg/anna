package anna

import (
	"log"
	"runtime"
	"time"
)

func (cmd *Cmd) PrintStats(elapsedTime time.Duration) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	log.Printf("Memory Usage: %d bytes", memStats.Alloc)
	log.Printf("Time Elapsed: %s", elapsedTime)
	cpuUsage := runtime.NumCPU()
	threads := runtime.GOMAXPROCS(0)
	runtime.ReadMemStats(memStats)

	log.Printf("Threads: %d", threads)
	log.Printf("Cores: %d", cpuUsage)
	log.Printf("Time Taken: %s", elapsedTime)
	log.Printf("Allocated Memory: %d bytes", memStats.Alloc)
	log.Printf("Total Memory Allocated: %d bytes", memStats.TotalAlloc)
	log.Printf("Heap Memory In Use: %d bytes", memStats.HeapInuse)
	log.Printf("Heap Memory Idle: %d bytes", memStats.HeapIdle)
	log.Printf("Heap Memory Released: %d bytes", memStats.HeapReleased)
	log.Printf("Number of Goroutines: %d", runtime.NumGoroutine())

	// Get the function with the highest CPU usage
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	log.Printf("Function with Highest CPU Usage: %s", function.Name())
}
