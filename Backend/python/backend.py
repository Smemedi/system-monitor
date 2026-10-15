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




def Snapshot():
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


def Storage():
    print("WIP")




def main(prevSnapshot: list[str]):
    while 1:
        Uptime()
        Memory()

        currSnapshot = Snapshot()
        CPU(prevSnapshot, currSnapshot)
        prevSnapshot = currSnapshot

        Storage()

        time.sleep(5)
        print('\33[4A', end='')

if __name__ == "__main__":
    prevSnapshot = Snapshot()
    time.sleep(3)
    main(prevSnapshot)
