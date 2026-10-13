import re
import time

def Uptime():

    uptimeFile = open("/proc/uptime", "r")

    uptimeInfo = uptimeFile.read();
    
    match = re.search(r"^\d+", uptimeInfo)
    uptimeMatch = match.group()

    uptime = time.strftime('%H:%M:%S', time.gmtime(int(uptimeMatch)))
    
    print(uptime)

def Memory():
    memoryFile = open("/proc/meminfo", "r")

    meminfo = memoryFile.read()

    match = re.search(r"^MemTotal.*?(\d+)", meminfo, re.MULTILINE)
    memTotal = int(match.group(1))


    match = re.search(r"^MemAvailable.*?(\d+)", meminfo, re.MULTILINE)
    memAvailable = int(match.group(1))

    print(
        str(round((memTotal - memAvailable) / 1048576, 2))
            + 'GB / ' +
            str(round(memTotal / 1048576, 2))
            + 'GB')

def CPU():
    print("Nothing yet")

def main():
    Uptime()
    Memory()
    CPU()

if __name__ == "__main__":
    main()
