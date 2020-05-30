package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shirou/gopsutil/process"
)

func main() {

	pidArg := flag.Int("pid", 0, "PID of the process to throttle. If provided with -pids, will be merged with list of pids.")
	pidsArg := flag.String("pids", "[]", "Comma separated list of PIDs of the processes to throttle. Example: -pids='[1,2,3,4]'. If provided with -pid, -pid value will be added to the list.")
	cpu := flag.Int("cpu", 100, "Percentage of CPU limit, from 0 to 100. Setting to >=100 will have no effect. Setting to 0 will freeze processes. Setting <0 will kill processes.")
	flag.Parse()

	if *pidArg == 0 && *pidsArg == "" {
		log.Fatal("Incorrect params")
	}

	var procPids []int
	err := json.Unmarshal([]byte(*pidsArg), &procPids)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	if *pidArg != 0 {
		procPids = append(procPids, *pidArg)
	}

	var processes []*process.Process

	for _, pid := range procPids {
		p, pidErr := process.NewProcess(int32(pid))

		if pidErr != nil {
			log.Fatal(pidErr)
			return
		}
		processes = append(processes, p)
	}

	if *cpu >= 100 {
		resumeSuspended(processes)
		return
	}

	if *cpu == 0 {
		suspendProcesses(processes)
		return
	}

	if *cpu < 0 {
		for _, proc := range processes {
			killResult := proc.Kill()
			if killResult != nil {
				log.Fatal(killResult)
			}
		}
		return
	}

	throttle(processes, *cpu)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Send SIGINT to remove limit")
	<-done

	resumeSuspended(processes)

	fmt.Println("Exiting and removing CPU limit")
}
