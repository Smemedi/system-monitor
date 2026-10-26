package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"bufio"
);

var(
	bibToGiB = 1048576;
)

type jiffies struct{
	user int64
	nice int64
	system int64
	idle int64
	iowait int64
	irq int64
	softirq int64
	steal int64
}



func Uptime(){
	uptimeFile, err := os.ReadFile("/proc/uptime")
	if(err != nil){
		panic(err);
	}

	uptime := strings.Fields(string(uptimeFile))[0]
	uptimeFloat, err := strconv.ParseFloat(uptime, 64);
	if(err != nil){
		panic(err);
	}
	totalSeconds := int(uptimeFloat);

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	fmt.Printf("Uptime: \n\t%d:%02d:%02d\n", hours, minutes, seconds);

}



func Memory(){
	memoryFile, err := os.Open("/proc/meminfo");
	if (err != nil){
		panic(err);
	}

	scanner := bufio.NewScanner(memoryFile);

	memTotalFloat := 0.0;
	memAvailableFloat := 0.0;

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "MemTotal:"){
			memTotal := strings.Fields(scanner.Text())[1];
			memTotalFloat, err = strconv.ParseFloat(memTotal, 64);

			if(err != nil){
				panic(err);
			}
		}

		if strings.Contains(scanner.Text(), "MemAvailable:"){
			memAvailable := strings.Fields(scanner.Text())[1];
			memAvailableFloat, err = strconv.ParseFloat(memAvailable, 64);

			if(err != nil){
				panic(err);
			}
		}
	}

	if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
	}

	memoryUsed := (memTotalFloat - memAvailableFloat) / float64(bibToGiB);
	fmt.Printf("Memory: \n\t%.2f Gib / %.2f GiB\n", memoryUsed, (memTotalFloat) / float64(bibToGiB));
	memoryFile.Close();
}



func parseJiffie(jiffie string) int64 {
	parsed, err := strconv.ParseInt(jiffie, 10, 64);
	if (err != nil){
		panic(err)
	}

	return parsed
}

func CPUSnapshot() *jiffies {
	j := jiffies{}
	cpuFile, err := os.Open("/proc/stat")
	if (err != nil){
		panic(err)
	}

	scanner := bufio.NewScanner(cpuFile)
	scanner.Scan()
	jiffiesString := strings.Fields(scanner.Text());

	j.user = parseJiffie(jiffiesString[1])
	j.nice = parseJiffie(jiffiesString[2])
	j.system = parseJiffie(jiffiesString[3])
	j.idle = parseJiffie(jiffiesString[4])
	j.iowait = parseJiffie(jiffiesString[5])
	j.irq = parseJiffie(jiffiesString[6])
	j.softirq = parseJiffie(jiffiesString[7])
	j.steal = parseJiffie(jiffiesString[8])

	defer cpuFile.Close()
	return &j
}

func CPU(cpuPrevSnapshot *jiffies, cpuCurrSnapshot *jiffies){
	deltaUser := cpuCurrSnapshot.user - cpuPrevSnapshot.user
	deltaNice := cpuCurrSnapshot.nice - cpuPrevSnapshot.nice
	deltaSystem := cpuCurrSnapshot.system - cpuPrevSnapshot.system
	deltaIdle := cpuCurrSnapshot.idle - cpuPrevSnapshot.idle
	deltaIowait := cpuCurrSnapshot.iowait - cpuPrevSnapshot.iowait
	deltaIrq := cpuCurrSnapshot.irq - cpuPrevSnapshot.irq
	deltaSoftirq := cpuCurrSnapshot.softirq - cpuPrevSnapshot.softirq
	deltaSteal := cpuCurrSnapshot.steal - cpuPrevSnapshot.steal

	total := deltaUser + deltaNice + deltaSystem + deltaIdle + deltaIowait + deltaIrq + deltaSoftirq + deltaSteal
	idleJiffies := deltaIdle + deltaIowait
	busyJiffies := total - idleJiffies

	cpu := (float64(busyJiffies) / float64(total)) * 100

	fmt.Printf("CPU: \n\t%.2f%%", cpu)

}



func main(){
	cpuPrevSnapshot := CPUSnapshot();
	time.Sleep(3 * time.Second);

	for {
		fmt.Print("\033[H\033[2J");
		Uptime();
		Memory();

		cpuCurrSnapshot := CPUSnapshot();
		CPU(cpuPrevSnapshot, cpuCurrSnapshot);
		cpuPrevSnapshot = cpuCurrSnapshot;

		time.Sleep(5 * time.Second);
	}
}
