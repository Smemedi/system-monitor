import re
import time
from flask import Flask, jsonify

app = Flask(__name__)

class SystemMonitor:
    def __init__(self) -> None:
        self.cpuPrevSnapshot: list[str] = CPUSnapshot()
        self.storagePrevSnapshot: list[list[str]]= StorageSnapshot()
        self.storageTime: float = time.monotonic()
        self.networkPrevSnapshot: list[str] = NetworkSnapshot()
        self.networkTime: float = time.monotonic()

@app.route('/api/uptime', methods=['GET'])
def Uptime():
    uptimeFile = open("/proc/uptime", "r")

    uptimeInfo = uptimeFile.read();
    
    match = re.search(r"^\d+", uptimeInfo)
    uptimeMatch = match.group()  # pyright: ignore[reportOptionalMemberAccess]

    uptime = time.strftime('%H:%M:%S', time.gmtime(int(uptimeMatch)))
    
    json_data = {
        'uptime': uptime
    }
    return jsonify(json_data)


@app.route('/api/memory', methods=['GET'])
def Memory():
        memoryFile = open("/proc/meminfo", "r")

        meminfo = memoryFile.read()

        match = re.search(r"^MemTotal.*?(\d+)", meminfo, re.MULTILINE)
        memTotal = int(match.group(1))  # pyright: ignore[reportOptionalMemberAccess]


        match = re.search(r"^MemAvailable.*?(\d+)", meminfo, re.MULTILINE)
        memAvailable = int(match.group(1))  # pyright: ignore[reportOptionalMemberAccess]
        json_data = {
            'usedMem': round((memTotal - memAvailable) / 1048576, 2),
            'memAvailable': round(memAvailable / 1048576, 2)
        }
        return jsonify(json_data)



def CPUSnapshot() -> list[str]:
    cpuFile = open("/proc/stat", "r")
    cpuInfo = cpuFile.read()

    match = re.search(r"^cpu\s+(\d+(?:\s+\d+)*)", cpuInfo, re.MULTILINE)
    return match.group(1).split()  # pyright: ignore[reportOptionalMemberAccess]

@app.route('/api/cpu', methods=['GET'])
def CPU():
    delta: list[int] = []
    cpuPrevSnapshot = monitor.cpuPrevSnapshot
    cpuCurrSnapshot = CPUSnapshot()


    for curr, prev in zip(cpuCurrSnapshot, cpuPrevSnapshot):
        delta.append(abs(int(curr) - int(prev)))

    total_jiffies = sum(delta)
    idle_jiffies = delta[3] + delta[4]

    busy_jiffies = total_jiffies - idle_jiffies

    if total_jiffies == 0:
        cpu = 0
    else:
        cpu = round((busy_jiffies / total_jiffies) * 100, 2)

    json_data = {
        "cpu%": cpu
    }

    monitor.cpuPrevSnapshot = cpuCurrSnapshot

    return jsonify(json_data)




def StorageSnapshot() -> list[list[str]]:
    drives: list[list[str]] = []

    storageFile = open("/proc/diskstats", "r")
    storageInfo = storageFile.read()

    matches = re.finditer(r"nvme\d+n\d\s.*$", storageInfo, re.MULTILINE)
    
    for match in matches:
        drives.append(match.group().split())

    return drives

@app.route('/api/storage', methods=['GET'])
def Storage():
    storagePrevSnapshot = monitor.storagePrevSnapshot
    start = monitor.storageTime
    storageCurrSnapshot = StorageSnapshot()
    end = time.monotonic()

    timeElapsed = end - start

    drives = {}

    for curr, prev in zip(storageCurrSnapshot, storagePrevSnapshot):
        driveName = curr[0]
        readSpeed = (((int(curr[3]) - int(prev[3])) * 512) / timeElapsed) / 1000000
        writeSpeed = (((int(curr[7]) - int(prev[7])) * 512) / timeElapsed) / 1000000

        drives[driveName] = {
            'readSpeed': round(readSpeed, 2),
            'writeSpeed': round(writeSpeed, 2)
        }

    monitor.storagePrevSnapshot = storageCurrSnapshot
    monitor.storageTime = end
    return jsonify(drives)



def NetworkSnapshot() -> list[str]:
    networkFile = open("/proc/net/dev", "r")

    networkInfo = networkFile.read()

    matches = re.search(r"wlp62s0:\s.*$", networkInfo, re.MULTILINE)

    return matches.group().split()  # pyright: ignore[reportOptionalMemberAccess]


@app.route('/api/network', methods=['GET'])
def Network():
    networkPrevSnapshot = monitor.networkPrevSnapshot
    start = monitor.networkTime
    networkCurrSnapshot = NetworkSnapshot()
    end = time.monotonic()

    downloadSpeed = ((int(networkCurrSnapshot[1]) - int(networkPrevSnapshot[1])) / (end - start)) / 1000000
    uploadSpeed = ((int(networkCurrSnapshot[9]) - int(networkPrevSnapshot[9])) / (end - start)) / 1000000

    json_data = {
        'downloadSpeed': round(downloadSpeed, 2),
        'uploadSpeed': round(uploadSpeed, 2)
    }

    monitor.networkPrevSnapshot = networkCurrSnapshot
    monitor.networkTime = end

    return jsonify(json_data)


if __name__ == "__main__":
    monitor = SystemMonitor()
    app.run()
    
