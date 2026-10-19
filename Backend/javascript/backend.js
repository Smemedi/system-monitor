import { readFileSync } from 'fs';

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




async function main(){
	while(true){
		console.clear();
		Uptime();
		Memory();

		await new Promise(r => setTimeout(r, 5000));
	}
}

main()
