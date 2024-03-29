package anna

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func PrintStats(elapsedTime time.Duration) {
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
	print("-----------------------------------\n")
}

func StopProfiling() {
	pprof.StopCPUProfile()
}

func StartProfiling() {
	runtime.MemProfileRate = 1
	go func() {
		for {
			time.Sleep(5 * time.Second)
			memProfile := pprof.Lookup("heap")
			if memProfile != nil {
				file, err := os.Create("mem_profile.pprof")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				memProfile.WriteTo(file, 1)
			}
			cpuProfile := pprof.Lookup("goroutine")
			if cpuProfile != nil {
				file, err := os.Create("cpu_profile.pprof")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				cpuProfile.WriteTo(file, 1)
			}
		}
	}()
}
