import { readFileSync } from 'fs';
import { performance } from 'node:perf_hooks';

function Uptime(){
	let uptimeInfo = readFileSync("/proc/uptime", 'utf8');
	
	let uptime = uptimeInfo.match(/^\d+/gm);

	let hours = Math.floor(uptime[0] / 3600)
	uptime[0] = uptime[0] % 3600;
	let minutes = Math.floor(uptime[0] / 60);
	uptime[0] = uptime[0] % 60
	let seconds = uptime[0]

	console.log(hours + ":" + String(minutes).padStart(2, '0') + ":" + String(seconds).padStart(2, '0'))
}



function Memory(){
	let memInfo = readFileSync("/proc/meminfo", 'utf8');
	
	let memTotal = memInfo.match(/^MemTotal.*?(\d+)/m);
	let memAvailable = memInfo.match(/^MemAvailable.*?(\d+)/m);

	console.log("Memory: " +
		((memTotal[1] - memAvailable[1]) / 1048576).toFixed(2) + "GiB / " +
		(Number(memTotal[1] / 1048576)).toFixed(2) + "GiB")
}



function CPUSnapshot(){
	let cpuInfo = readFileSync("/proc/stat", 'utf8');
	let cpuStats = cpuInfo.match(/^cpu\s+(\d+(?:\s+\d+)*)/m)[1].split(/\s+/);

	return cpuStats;
}

function CPU(cpuPrevSnapshot, cpuCurrSnapshot){
	let delta = []
	for(let i = 0 ; i < cpuPrevSnapshot.length ; i++){
		delta.push( cpuCurrSnapshot[i] - cpuPrevSnapshot[i] );
	}

	let totalJiffies = 0;
	delta.forEach(x => {
			totalJiffies += x;
	});

	let idleJiffies = delta[3] + delta[4];
	let busyJiffies = totalJiffies - idleJiffies;

	let cpu = ((busyJiffies / totalJiffies) * 100).toFixed(2)

	console.log(cpu +"%");
}



function StorageSnapshot(){
	let storageFile = readFileSync("/proc/diskstats", 'utf8');
	let storageInfo = storageFile.match(/nvme\d+n\d\s.*$/mg);

	let drives = [];

	for(let i = 0 ; i < storageInfo.length ; i++){
		drives.push(storageInfo[i].split(/\s+/));
	}

	return drives;
}

function Storage(storagePrevSnapshot, storageCurrSnapshot, start, end){
	let timeElapsed = end - start;

	for(let i = 0 ; i < storageCurrSnapshot.length ; i++){
		let driveName = storageCurrSnapshot[i][0];
		let readSpeed = (((storageCurrSnapshot[i][3] - storagePrevSnapshot[i][3]) * 512) / timeElapsed) / 1000000;
		let writeSpeed = (((storageCurrSnapshot[i][7] - storagePrevSnapshot[i][7]) * 512) / timeElapsed) / 1000000;

		console.log(driveName + ":\n\t" +
				"Read speed: " + readSpeed.toFixed(2) + " MB/s\n\t" +
				"Write speed: " + writeSpeed.toFixed(2) + " MB/s");
	}
}

function NetworkSnapshot(){
	let networkFile = readFileSync("/proc/net/dev", 'utf8');
	let networkInfo = networkFile.match(/wlp62s0:\s.*$/gm)[0].split(/\s+/);

	return networkInfo;
}

function Network(networkPrevSnapshot, networkCurrSnapshot, start, end){
	let downloadSpeed = ((networkCurrSnapshot[1] - networkPrevSnapshot[1]) / (end - start)) / 1000000;
	let uploadSpeed = ((networkCurrSnapshot[9] - networkPrevSnapshot[9]) / (end - start)) / 1000000;

	console.log("Network Speed: \n\t" +
			"Download Speed: " + downloadSpeed.toFixed(2) + " MB/s\n\t" +
			"Upload Speed: " + uploadSpeed.toFixed(2) + " MB/s");
}



async function main(cpuPrevSnapshot, storagePrevSnapshot, networkPrevSnapshot){
	let start = performance.now() / 1000;

	while(true){
		console.clear();
		Uptime();
		Memory();

		let cpuCurrSnapshot = CPUSnapshot();
		CPU(cpuPrevSnapshot, cpuCurrSnapshot);
		cpuPrevSnapshot = cpuCurrSnapshot;

		let storageCurrSnapshot = StorageSnapshot();
		let networkCurrSnapshot = NetworkSnapshot();
		let end = performance.now() / 1000;

		Storage(storagePrevSnapshot, storageCurrSnapshot, start, end);
		Network(networkPrevSnapshot, networkCurrSnapshot, start, end);

		storagePrevSnapshot = storageCurrSnapshot;
		networkPrevSnapshot = networkCurrSnapshot;
		start = end;


		await new Promise(r => setTimeout(r, 5000));
	}
}

let cpuPrevSnapshot = CPUSnapshot();
let storagePrevSnapshot = StorageSnapshot();
let networkPrevSnapshot = NetworkSnapshot();
main(cpuPrevSnapshot, storagePrevSnapshot, networkPrevSnapshot);
