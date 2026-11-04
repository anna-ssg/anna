package anna

import (
	"runtime"
	"time"
)

func (cmd *Cmd) PrintStats(elapsedTime time.Duration) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	cmd.InfoLogger.Info().Msgf("Memory Usage: %d bytes", memStats.Alloc)
	cmd.InfoLogger.Info().Msgf("Time Elapsed: %s", elapsedTime)
	cpuUsage := runtime.NumCPU()
	threads := runtime.GOMAXPROCS(0)
	runtime.ReadMemStats(memStats)

	cmd.InfoLogger.Info().Msgf("Threads: %d", threads)
	cmd.InfoLogger.Info().Msgf("Cores: %d", cpuUsage)
	cmd.InfoLogger.Info().Msgf("Time Taken: %s", elapsedTime)
	cmd.InfoLogger.Info().Msgf("Allocated Memory: %d bytes", memStats.Alloc)
	cmd.InfoLogger.Info().Msgf("Total Memory Allocated: %d bytes", memStats.TotalAlloc)
	cmd.InfoLogger.Info().Msgf("Heap Memory In Use: %d bytes", memStats.HeapInuse)
	cmd.InfoLogger.Info().Msgf("Heap Memory Idle: %d bytes", memStats.HeapIdle)
	cmd.InfoLogger.Info().Msgf("Heap Memory Released: %d bytes", memStats.HeapReleased)
	cmd.InfoLogger.Info().Msgf("Number of Goroutines: %d", runtime.NumGoroutine())

	// Get the function with the highest CPU usage
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	cmd.InfoLogger.Info().Msgf("Function with Highest CPU Usage: %s", function.Name())
}
