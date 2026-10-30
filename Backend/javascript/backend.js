import { readFileSync } from 'fs';
import { performance } from 'node:perf_hooks';
import express from 'express';

const app = express();
const port = 8080;


class SystemMonitor {
	constructor() {
		this.cpuPrevSnapshot = CPUSnapshot();
		this.storagePrevSnapshot = StorageSnapshot();
		this.storageTime = performance.now();
		this.networkPrevSnapshot = NetworkSnapshot();
		this.networkTime = performance.now();
	}
}

function Uptime(){
	let uptimeInfo = readFileSync("/proc/uptime", 'utf8');
	
	let uptime = uptimeInfo.match(/^\d+/gm);

	let hours = Math.floor(uptime[0] / 3600)
	uptime[0] = uptime[0] % 3600;
	let minutes = Math.floor(uptime[0] / 60);
	uptime[0] = uptime[0] % 60
	let seconds = uptime[0]

	return(hours + ":" + String(minutes).padStart(2, '0') + ":" + String(seconds).padStart(2, '0'))
}



function Memory(){
	let memInfo = readFileSync("/proc/meminfo", 'utf8');
	
	let memTotal = memInfo.match(/^MemTotal.*?(\d+)/m);
	let memAvailable = memInfo.match(/^MemAvailable.*?(\d+)/m);

	return {
		"used": ((memTotal[1] - memAvailable[1]) / 1048576).toFixed(2),
		"total": (memTotal[1] / 1048576).toFixed(2)
	};
}



function CPUSnapshot(){
	let cpuInfo = readFileSync("/proc/stat", 'utf8');
	let cpuStats = cpuInfo.match(/^cpu\s+(\d+(?:\s+\d+)*)/m)[1].split(/\s+/);

	return cpuStats;
}

function CPU(){
	let delta = [];
	let cpuPrevSnapshot = monitor.cpuPrevSnapshot
	let cpuCurrSnapshot = CPUSnapshot();

	for(let i = 0 ; i < cpuPrevSnapshot.length ; i++){
		delta.push( cpuCurrSnapshot[i] - cpuPrevSnapshot[i] );
	}

	let totalJiffies = 0;
	delta.forEach(x => {
			totalJiffies += x;
	});

	let idleJiffies = delta[3] + delta[4];
	let busyJiffies = totalJiffies - idleJiffies;

	let cpu = ((busyJiffies / totalJiffies) * 100).toFixed(2);

	monitor.cpuPrevSnapshot = cpuCurrSnapshot;

	return cpu;
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

function Storage(){
	let storagePrevSnapshot = monitor.storagePrevSnapshot;
	let storageCurrSnapshot = StorageSnapshot();
	let start = monitor.storageTime;
	let end = performance.now();

	let timeElapsed = (end - start) / 1000;

	let drives = {};

	for(let i = 0 ; i < storageCurrSnapshot.length ; i++){
		let driveName = storageCurrSnapshot[i][0];
		let readSpeed = (((storageCurrSnapshot[i][3] - storagePrevSnapshot[i][3]) * 512) / timeElapsed) / 1000000;
		let writeSpeed = (((storageCurrSnapshot[i][7] - storagePrevSnapshot[i][7]) * 512) / timeElapsed) / 1000000;

		drives[driveName] = {
			"read": readSpeed,
			"write": writeSpeed
		}
	}

	monitor.storagePrevSnapshot = storageCurrSnapshot;
	monitor.storageTime = end

	return drives;
}

function NetworkSnapshot(){
	let networkFile = readFileSync("/proc/net/dev", 'utf8');
	let networkInfo = networkFile.match(/wlp62s0:\s.*$/gm)[0].split(/\s+/);

	return networkInfo;
}

function Network(){
	let networkPrevSnapshot = monitor.networkPrevSnapshot;
	let networkCurrSnapshot = NetworkSnapshot();
	let end = performance.now();
	let start = monitor.networkTime;

	let timeElapsed = (end - start) / 1000;

	let downloadSpeed = ((networkCurrSnapshot[1] - networkPrevSnapshot[1]) / timeElapsed) / 1000000;
	let uploadSpeed = ((networkCurrSnapshot[9] - networkPrevSnapshot[9]) / timeElapsed) / 1000000;

	monitor.networkPrevSnapshot = networkCurrSnapshot;
	monitor.networkTime = end;

	return {
		"Download Speed": downloadSpeed.toFixed(2),
		"Upload Speed": uploadSpeed.toFixed(2)
	};
}



const monitor = new SystemMonitor();

app.get('/api/uptime', (_, res) => {

  res.json({
			'uptime': Uptime()
			});
});

app.get('/api/memory', (_, res) => {

  res.json(Memory());
});

app.get('/api/cpu', (_, res) => {

  res.json({
			'CPU': CPU()
			});
});

app.get('/api/storage', (_, res) => {

  res.json(Storage());
});

app.get('/api/network', (_, res) => {

  res.json(Network());
});


app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`);
});

