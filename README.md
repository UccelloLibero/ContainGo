# ContainGo: A lightweight containerization tool built with Go

ContainGo is a **learning project** for exploring **containerization concepts** using Go.  
It allows you to run and manage isolated containers using `chroot`.
It is designed for Linux and macOS but lacks full namespace support on macOS.

---

## **Features**

- Run a process inside an isolated root filesystem.  
- List running container processes.  
- Stop a running container.  
- Works on ***Linux*** and ***macOS***.  

---

## **About this project**
This project is a part of my learning journey to understand containerization, Linux namespaces, cgroups, and process isolation. It is not a full-fledged container runtime like Docker, but rather a hands-on experiment to learn how containers work internally. 

---

## **Installation**

### **Clone the repository**
```sh
git clone https://github.com/UccelloLibero/ContainGo.git
cd ContainGo
```

### **Install dependencies**
Run:
```sh
go mod tidy
```

### **Build the binary**
Run:
```sh
go build -o containgo
```

### **Verify the CLI**
Run:
```sh
./containgo --help
```

Expected output:
```
ContainGo allows you to run and manage isolated containers.
It provides filesystem isolation using chroot and executes commands within a containerized environment.

Usage:
  contain-go [flags]
  contain-go [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List all running containers
  run         Run a process inside an isolated root filesystem
  stop        Stop a running container

Flags:
      --config string   Config file (default is $HOME/.contain-go.yaml)
  -h, --help            help for contain-go

Use "contain-go [command] --help" for more information about a command.
```

---

## **Usage**
Since ContainGo uses `chroot`, it needs a minimal root filesystem:
```sh
mkdir -p test_rootfs/bin
cp /bin/sh test_rootfs/bin/  # Copy a shell into the test filesystem
```
Now `test_rootfs`is ready.

---

## **Run a  container**
Run:
```sh
./containgo run test_rootfs
```
Expected: You enter a shell inside a new root filesystem.
Run `ls /` inside the container to verify.
Exit using:
```sh
exit
```

---

## **Run a container in the background**
Run:
```sh
./containgo run test_rootfs &
```
The `&` runs it in the background.

---

## **List running containers**
Run:
```sh
./containgo list
```
Expected output:
```
Container ID | PID
12345        | 56789
```

---

## **Stopping a running container**
Find the PID:
```sh
ps aux | grep containgo
```
Stop the process:
```sh
./containgo stop <PID>
```

Example:
```sh
./containgo stop 56789
```
Expected: The container is killed.

---

## **Error handling**
1. Stopping a non-existent process:
```sh
./containgo stop 999999
```
Expected: Error message.

2. Running `run` without valid `rootfs`
```sh
./containgo run non_existent_rootfs
```
Expected: Error message saying `rootfs` is missing.

---

## **Troubleshooting**
1. Binary not found after build
Try rebuilding: 
```sh
go build -o containgo
```
Make sure you are in correct directory.

2. Process not stopping
Run:
```sh
ps aux | grep containgo
kill <PID>
```

3. `chroot` errors
Make sure `/bin/sh` exists inside your rootfs:
```sh
ls test_rootfs/bin/
```

---

### **Author**
Maya
GitHub: [UccelloLibero](https://github.com/UccelloLibero)

## **License**
This project is licensed under the MIT License.