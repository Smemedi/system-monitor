def Uptime():

    uptimeFile = open("/proc/uptime", "r")

    uptime = uptimeFile.read();
    print(uptime)

def Memory():
    print("Nothing yet")

def CPU():
    print("Nothing yet")

def main():
    Uptime()
    Memory()
    CPU()

if __name__ == "__main__":
    main()
