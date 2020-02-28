package main

import (
	"flag"
	"log"
	"time"

	"github.com/shirou/gopsutil/process"
)

func main() {

	pid := flag.Int("pid", 0, "PID of the process to throttle")
	cpu := flag.Int("cpu", 100, "Percentage of CPU limit, from 0 to 100. Setting to >=100 will have no effect. Setting to 0 will freeze process. Setting <0 will kill process.")
	flag.Parse()

	if *pid == 0 {
		log.Fatal("Incorrect params")
	}

	p, pidErr := process.NewProcess(int32(*pid))

	if pidErr != nil {
		log.Fatal(pidErr)
	}

	if *cpu >= 100 {
		return
	}
	if *cpu <= 0 {
		killResult := p.Kill()
		if killResult != nil {
			log.Fatal(killResult)
		}
		return
	}
	for true {
		p.Suspend()
		time.Sleep(time.Duration(10*(100-(*cpu))) * time.Microsecond)
		p.Resume()
		time.Sleep(time.Duration(10*(*cpu)) * time.Microsecond)
	}
}
