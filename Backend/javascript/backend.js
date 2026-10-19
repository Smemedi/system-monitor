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
		((memTotal[1] - memAvailable[1]) / 1048576).toFixed(2) + "GB/s / " +
		(Number(memTotal[1] / 1048576)).toFixed(2) + "GB/s")
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
	console.log(delta);

	let totalJiffies = 0;
	delta.forEach(x => {
			totalJiffies += x;
	});

	let idleJiffies = delta[3] + delta[4];
	let busyJiffies = totalJiffies - idleJiffies;

	let cpu = ((busyJiffies / totalJiffies) * 100).toFixed(2)

	console.log(cpu +"%");
}


async function main(cpuPrevSnapshot){
	let start = performance.now();

	while(true){
		console.clear();
		Uptime();
		Memory();

		let cpuCurrSnapshot = CPUSnapshot();
		CPU(cpuPrevSnapshot, cpuCurrSnapshot);
		cpuPrevSnapshot = cpuCurrSnapshot;

		await new Promise(r => setTimeout(r, 5000));
	}
}

let cpuPrevSnapshot = CPUSnapshot();
main(cpuPrevSnapshot);
