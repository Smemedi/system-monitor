import re
import time

def Uptime():

    uptimeFile = open("/proc/uptime", "r")

    uptimeInfo = uptimeFile.read();
    
    match = re.search(r"^\d+", uptimeInfo)
    uptimeMatch = match.group()  # pyright: ignore[reportOptionalMemberAccess]

    uptime = time.strftime('%H:%M:%S', time.gmtime(int(uptimeMatch)))
    
    print(uptime)



def Memory():
    memoryFile = open("/proc/meminfo", "r")

    meminfo = memoryFile.read()

    match = re.search(r"^MemTotal.*?(\d+)", meminfo, re.MULTILINE)
    memTotal = int(match.group(1))  # pyright: ignore[reportOptionalMemberAccess]


    match = re.search(r"^MemAvailable.*?(\d+)", meminfo, re.MULTILINE)
    memAvailable = int(match.group(1))  # pyright: ignore[reportOptionalMemberAccess]

    print(
        str(round((memTotal - memAvailable) / 1048576, 2))
            + 'GB / ' +
            str(round(memTotal / 1048576, 2))
            + 'GB')




def CPUSnapshot():
    cpuFile = open("/proc/stat", "r")
    cpuInfo = cpuFile.read()

    match = re.search(r"^cpu\s+(\d+(?:\s+\d+)*)", cpuInfo, re.MULTILINE)
    return match.group(1).split()  # pyright: ignore[reportOptionalMemberAccess]

def CPU(prevSnapshot: list[str], currSnapshot: list[str]):
    delta: list[int] = []

    for curr, prev in zip(currSnapshot, prevSnapshot):
        delta.append(abs(int(curr) - int(prev)))

    total_jiffies = sum(delta)
    idle_jiffies = delta[3] + delta[4]

    busy_jiffies = total_jiffies - idle_jiffies

    cpu = round((busy_jiffies / total_jiffies) * 100, 2)

    print(str(cpu) + '%')




def StorageSnapshot() -> list[list[str]]:
    drives: list[list[str]] = []

    storageFile = open("/proc/diskstats", "r")
    storageInfo = storageFile.read()

    matches = re.finditer(r"nvme\d+n\d\s.*$", storageInfo, re.MULTILINE)
    
    for match in matches:
        drives.append(match.group().split())

    return drives


def Storage(storagePrevSnapshot: list[list[str]], storageCurrSnapshot: list[list[str]], start: float, end: float):
    timeElapsed = end - start

    for curr, prev in zip(storageCurrSnapshot, storagePrevSnapshot):
        driveName = curr[0]
        readSpeed = (((int(curr[3]) - int(prev[3])) * 512) / timeElapsed) / 1000000
        writeSpeed = (((int(curr[7]) - int(prev[7])) * 512) / timeElapsed) / 1000000

        print(driveName[:-2] +":\n\t" +
            "Read: " + str(round(readSpeed, 2)) + "MB/s\n\t" +
            "Write: " + str(round(writeSpeed, 2)) + "MB/s")




def NetworkSnapshot():
    networkFile = open("/proc/net/dev", "r")

    networkInfo = networkFile.read()

    matches = re.search(r"wlp62s0:\s.*$", networkInfo, re.MULTILINE)

    return matches.group().split()  # pyright: ignore[reportOptionalMemberAccess]


def Network(networkPrevSnapshot: list[str], networkCurrSnapshot: list[str], start: float, end: float):
    downloadSpeed = ((int(networkCurrSnapshot[1]) - int(networkPrevSnapshot[1])) / (end - start)) / 1000000
    uploadSpeed = ((int(networkCurrSnapshot[9]) - int(networkPrevSnapshot[9])) / (end - start)) / 1000000

    print("Network speed:\n\t" +
        "Download Speed: " + str(round(downloadSpeed, 2)) + "MB/s\n\t" +
        "Upload Speed: " + str(round(uploadSpeed, 2)) + "MB/s")




def main(cpuPrevSnapshot: list[str],
         storagePrevSnapshot: list[list[str]],
         networkPrevSnapshot: list[str]):

    start = time.monotonic()

    while 1:
        Uptime()
        Memory()

        currSnapshot = CPUSnapshot()
        CPU(cpuPrevSnapshot, currSnapshot)
        cpuPrevSnapshot = currSnapshot

        storageCurrSnapshot = StorageSnapshot()
        networkCurrSnapshot = NetworkSnapshot()
        end = time.monotonic()
        Storage(storagePrevSnapshot, storageCurrSnapshot, start, end)
        Network(networkPrevSnapshot, networkCurrSnapshot, start, end)
        start = end
        storagePrevSnapshot = storageCurrSnapshot

        time.sleep(5)
        print('\33[12A', end='')

if __name__ == "__main__":
    cpuPrevSnapshot = CPUSnapshot()
    storagePrevSnapshot = StorageSnapshot()
    networkPrevSnapshot = NetworkSnapshot()
    
    time.sleep(3)
    main(cpuPrevSnapshot, storagePrevSnapshot, networkPrevSnapshot)
