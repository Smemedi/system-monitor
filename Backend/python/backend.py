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


def StorageSnapshot() -> list[str]:
    print("wip")


def Storage(storagePrevSnapshot: list[str], storageCurrSnapshot: list[str]):
    drives: list[str] = []

    storageFile = open("/proc/diskstats", "r")
    storageInfo = storageFile.read()

    matches = re.finditer(r"nvme\d+n\d\s.*$", storageInfo, re.MULTILINE)
    
    for match in matches:
        drives.append(match.group().split())  # pyright: ignore[reportArgumentType]

    for drive in drives:
        readSpeed = 5
        writeSpeed = 5

        print(drive[0] + ":\n\t" + ""
              )



def main(cpuPrevSnapshot: list[str], storagePrevSnapshot: list[str]):
    while 1:
        Uptime()
        Memory()

        currSnapshot = CPUSnapshot()
        CPU(cpuPrevSnapshot, currSnapshot)
        cpuPrevSnapshot = currSnapshot

        storageCurrSnapshot = StorageSnapshot()
        Storage(storagePrevSnapshot, storageCurrSnapshot)
        storagePrevSnapshot = storageCurrSnapshot

        time.sleep(5)
        print('\33[4A', end='')

if __name__ == "__main__":
    cpuPrevSnapshot = CPUSnapshot()
    storagePrevSnapshot = StorageSnapshot()
    
    time.sleep(3)
    main(cpuPrevSnapshot, storagePrevSnapshot)
