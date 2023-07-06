package main

/*
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <error.h>
#include <sys/ipc.h>
#include <sys/shm.h>
#include <sys/sem.h>

#define true 1
#define false 0
typedef unsigned char bool;

char *ShmGet();
bool shmwait();
bool shmpost();
bool shmuget(char *shmp);
bool shmremove();
int Add(int a, int b);
*/
import "C"
import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type Stock struct {
	Value int
	Incr  int
	Dump  int
}

type ShareMem struct {
	shmid   int
	semid   int
	shmaddr uintptr
	stock   *Stock
}

func main() {

	fmt.Println(C.Add(2, 8), C.CString("dadada"), *ShareMemGet())

	p := (*Stock)(unsafe.Pointer(C.ShmGet()))
	for {
		if a := C.shmwait(); a == 0 {
			break
		}
		p.Value -= 1
		p.Incr -= 1
		p.Dump -= 1
		fmt.Println(*p)
		if b := C.shmpost(); b == 0 {
			break
		}
		time.Sleep(time.Second)
	}
}

func ShareMemGet() *ShareMem {
	shmid, _, err := syscall.Syscall(syscall.SYS_SHMGET, 0x5005, 1024, 01000|0640)
	if int(shmid) == -1 {
		fmt.Println(err)
		os.Exit(-1)
	}

	semid, _, err := syscall.Syscall(syscall.SYS_SEMGET, 0x5000, 1, 01000|0640)
	if int(semid) == -1 {
		fmt.Println(err)
		os.Exit(-1)
	}

	shmaddr, _, err := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if int(shmaddr) == -1 {
		fmt.Println(err)
		os.Exit(-1)
	}

	semun := &C.union_semun{}
	fmt.Println(semun)
	ret, _, err := syscall.Syscall6(syscall.SYS_SEMCTL, semid, 0, 16, uintptr(unsafe.Pointer(semun)), 0, 0)
	if int(ret) == -1 {
		fmt.Println(err)
		os.Exit(-1)
	}

	return &ShareMem{
		shmid:   int(shmid),
		semid:   int(semid),
		shmaddr: shmaddr,
		stock:   (*Stock)(unsafe.Pointer(shmaddr)),
	}
}
