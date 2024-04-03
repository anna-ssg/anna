package anna

import (
	"log"
	"net/http"
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

	// Get the function with the highest CPU usage
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	log.Printf("Function with Highest CPU Usage: %s", function.Name())

	log.Println("-----------------------------------")
}

func StopProfiling() {
	pprof.StopCPUProfile()
}

func StartProfiling(annaCmd *Cmd) {
	file, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := pprof.StartCPUProfile(file); err != nil {
		log.Fatalf("%v", err)
	}
	defer pprof.StopCPUProfile()

	for i := 0; i < 500; i++ {
		annaCmd.VanillaRender()
	}

	memProfileFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer memProfileFile.Close()

	if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
		log.Fatalf("%v", err)
	}

	log.Println("-----------------------------------")
	log.Println("Memory and CPU profiles written to file")
}

func RunProfilingServer() {
	log.Println("Profiling server started at http://localhost:8080/debug/pprof")
	log.Println("-----------------------------------")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
