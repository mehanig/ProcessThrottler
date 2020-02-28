# ProcessThrottler

It's a simple CPU throttler.
It accepts 2 params:

- `-pid`: the pid of process to throttle
- `-cpu`: amount of throttling value in % from 1 to 99. Should be of type `Int`

It works by suspending and resuming the pid.
Cross-platform work is made possible by using this: "github.com/shirou/gopsutil/process"

This is not recommended for any production use.
