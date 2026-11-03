package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"net/http"
	"log"
	"encoding/json"
	"math"
);

var(
	bibToGiB = 1048576;
	bytetoMB = 1000000
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

type drive struct{
	driveName string
	read int64
	write int64
}

type networkData struct{
	upload float64
	download float64
}

type SystemMonitor struct{
	cpuPrevSnapshot *jiffies
	storagePrevSnapshot *[]drive
	storageTime time.Time
	networkPrevSnapshot *networkData
	networkTime time.Time
}



func Uptime( w http.ResponseWriter, r *http.Request ){
	w.Header().Set("Content-Type", "application/json")

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

	uptimeString := fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)

	json.NewEncoder(w).Encode(map[string]string{
			"Uptime": uptimeString,
	})

}



func Memory( w http.ResponseWriter, r *http.Request ){
	w.Header().Set("Content-Type", "application/json")

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
			break
		}
	}

	if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
	}

	memoryUsed := math.Round(((memTotalFloat - memAvailableFloat) / float64(bibToGiB)) * 100) /100;

	defer memoryFile.Close()
	
	json.NewEncoder(w).Encode(map[string]float64{
		"Memory Used": memoryUsed,
		"Memory Available": math.Round((memAvailableFloat / float64(bibToGiB)) * 100) / 100,
	})
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

func CPU( w http.ResponseWriter, r *http.Request, monitor *SystemMonitor ){
	w.Header().Set("Content-Type", "application/json")

	cpuPrevSnapshot := monitor.cpuPrevSnapshot
	cpuCurrSnapshot := CPUSnapshot()

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

	cpu := math.Round((float64(busyJiffies) / float64(total)) * 10000) / 100

	monitor.cpuPrevSnapshot = cpuCurrSnapshot

	json.NewEncoder(w).Encode(map[string]float64{
		"CPU": cpu,
	})

}



func StorageSnapshot() *[]drive{
	drives := []drive{}
	storageFile, err := os.Open("/proc/diskstats")
	if (err != nil){
		panic(err)
	}

	scanner := bufio.NewScanner(storageFile)

	for scanner.Scan(){
		fileString := strings.Fields(scanner.Text())
		driveName := fileString[2]


		if(driveName[len(driveName) - 2] != 'p' && driveName[:2] != "sd"){
			read, err := strconv.ParseInt(fileString[5], 10, 64)
			if(err != nil){
				panic(err)
			}

			write, err := strconv.ParseInt(fileString[9], 10, 64)
			if(err != nil){
				panic(err)
			}

			newDrive := drive{
			driveName: fileString[2],
			read: read,
			write: write,
			}

			drives = append(drives, newDrive)
		}

	}
	if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
	}

	defer storageFile.Close()
	return &drives
}

func Storage( w http.ResponseWriter, r *http.Request, monitor *SystemMonitor ){
	w.Header().Set("Content-Type", "application/json")

	storagePrevSnapshot := monitor.storagePrevSnapshot
	storageCurrSnapshot := StorageSnapshot()
	start := monitor.storageTime
	end := time.Now()

	storage := make(map[string]map[string]float64)

	timeElapsed := end.Sub(start).Seconds()
	for i := 0 ; i < len(*storageCurrSnapshot) ; i++ {
		
		driveName := (*storageCurrSnapshot)[i].driveName
		prev := (*storagePrevSnapshot)[i]
		curr := (*storageCurrSnapshot)[i]

		readSpeed := (float64(curr.read - prev.read) * float64(512)) / timeElapsed / float64(bytetoMB)
		writeSpeed := (float64(curr.write - prev.write) * float64(512)) / timeElapsed / float64(bytetoMB)

		storage[driveName] = map[string]float64{
			"read":  readSpeed,
			"write": writeSpeed,
		}
	}
	monitor.storagePrevSnapshot = storageCurrSnapshot
	monitor.storageTime = end

	json.NewEncoder(w).Encode(storage)
}



func NetworkSnapshot() *networkData{
	data := networkData{}
	networkFile, err := os.Open("/proc/net/dev")
	if(err != nil){
		panic(err)
	}
	scanner := bufio.NewScanner(networkFile)

	for scanner.Scan(){
		if( scanner.Text()[:7] == "wlp62s0"){
			networkString := strings.Fields(scanner.Text())

			data.download, err = strconv.ParseFloat(networkString[1], 64)
			if (err != nil){
				panic(err)
			}

			data.upload, err = strconv.ParseFloat(networkString[9], 64)
			if (err != nil){
				panic(err)
			}

			break

		}
	}
	if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
	}

	defer networkFile.Close()
	return &data
}

func Network( w http.ResponseWriter, r *http.Request, monitor *SystemMonitor ){
	w.Header().Set("Content-Type", "application/json")

	networkPrevSnapshot := monitor.networkPrevSnapshot
	networkCurrSnapshot := NetworkSnapshot()
	start := monitor.networkTime
	end := time.Now()

	timeElapsed := end.Sub(start).Seconds()

	deltaUp := ((networkCurrSnapshot.upload - networkPrevSnapshot.upload) / timeElapsed) / float64(bytetoMB)
	deltaDown := ((networkCurrSnapshot.download - networkPrevSnapshot.download) / timeElapsed) / float64(bytetoMB)

	monitor.networkPrevSnapshot = networkCurrSnapshot
	monitor.networkTime = end

	json.NewEncoder(w).Encode(map[string]float64{
		"Upload Speed": deltaUp,
		"Download Speed": deltaDown,
	})

}



func main(){

	monitor := SystemMonitor{
		cpuPrevSnapshot: CPUSnapshot(),
		storagePrevSnapshot: StorageSnapshot(),
		storageTime: time.Now(),
		networkPrevSnapshot: NetworkSnapshot(),
		networkTime: time.Now(),
	}

	http.HandleFunc("/api/uptime", Uptime)
	http.HandleFunc("/api/memory", Memory)
	http.HandleFunc("/api/cpu", func(w http.ResponseWriter, r *http.Request) {
			CPU(w, r, &monitor)
	})
	http.HandleFunc("/api/storage", func(w http.ResponseWriter, r *http.Request) {
			Storage(w, r, &monitor)
	})
	http.HandleFunc("/api/network", func(w http.ResponseWriter, r *http.Request) {
			Network(w, r, &monitor)
	})
  log.Fatal(http.ListenAndServe(":8080", nil))
	
}
