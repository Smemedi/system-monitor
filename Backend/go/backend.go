package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"regexp"
	"time"
);

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
	memoryFile, err := os.ReadFile("/proc/meminfo");
	if(err != nil){
		panic(err)
	}

	re := regexp.MustCompile(`(?m)^MemTotal:.*$`);
	match := re.FindString(string(memoryFile))
	stringTemp := strings.Fields(match)[1];
	floatTemp, err := strconv.ParseFloat(stringTemp, 64);
	if(err != nil){
		panic(err)
	}

	memoryTotal := floatTemp / 1048576;

	
	re = regexp.MustCompile(`(?m)^MemAvailable:.*$`);
	match = re.FindString(string(memoryFile))
	stringTemp = strings.Fields(match)[1];
	floatTemp, err = strconv.ParseFloat(stringTemp, 64);
	if(err != nil){
		panic(err)
	}
	 memoryUsed := memoryTotal - (floatTemp / 1048576);

	fmt.Printf("Memory: \n\t%.2f GiB / %.2f GiB\n", memoryUsed, memoryTotal);
}

func CPUSnapshot() []string {
	cpuFile, err := os.ReadFile("/proc/stat")
	if(err != nil){
		panic(err)
	}

	re := regexp.MustCompile(`(?m)^cpu\s+.*$`);
	match := re.FindString(string(cpuFile));
	jiffies := strings.Fields(match);

	return jiffies;
}

func CPU(cpuPrevSnapshot[]string, cpuCurrSnapshot []string){

	var delta [8]int;
	totalJiffies := 0;

	for index := 1 ; index < 9 ; index++{
		prev, err := strconv.Atoi(cpuPrevSnapshot[index]);
		if(err != nil){
			panic(err)
		}
		curr, err := strconv.Atoi(cpuCurrSnapshot[index]);
		if(err != nil){
			panic(err)
		}

		delta[index - 1] =  prev - curr;
		totalJiffies += prev - curr;
	}	

	idleJiffies := delta[3] + delta[4];
	busyJiffies := totalJiffies - idleJiffies;

	cpu := (float64(busyJiffies) / float64(totalJiffies)) * 100;

	fmt.Printf("CPU:\n\t%.2f%%", cpu);
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
