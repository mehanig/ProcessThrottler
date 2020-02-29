# ProcessThrottler

It's a simple CPU throttler.
It accepts 2 params:

- `-pid`: the pid of process to throttle
- `-pids`: Comma separated list of PIDs passed as JSON string array: `-pids='[1,2,3,4]'`
- `-cpu`: amount of throttling value in % from 1 to 99. Should be of type `Int`

It works by suspending and resuming the pid.
If both `pids` and `pid` args provided, values will be merged and throttling will be applied for every process.

Cross-platform work is made possible by using this: "github.com/shirou/gopsutil/process"

This is not recommended for any production use.
