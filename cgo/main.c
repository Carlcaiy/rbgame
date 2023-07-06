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

int shmid = 0;
int semid = 0;

union semun
{
    int val;
    struct semid_ds *buf;
    unsigned short *array;
};


int Add(int a, int b) {
    return a + b;
}

char *ShmGet() {
    shmid = shmget((key_t)0x5005, 1024, 0640|IPC_CREAT);
    if (shmid == -1) {
        return NULL;
    }
    semid = semget(0x5000, 1, 0640|IPC_CREAT);
    if (semid == -1) {
        return NULL;
    }

    union semun sem_union;
    sem_union.val = 1;
    semctl(semid, 0, SETVAL, sem_union);

    char *shmp = shmat(semid, 0, 0);
    if (shmp == (void *)-1) {
        return NULL;
    }
    return shmp;
}

bool shmwait() {
    struct sembuf sem_b;
    sem_b.sem_num = 0;
    sem_b.sem_op = -1;
    sem_b.sem_flg = SEM_UNDO;
    if (semop(shmid, &sem_b, 1) == -1) return false;
    return true;
}

bool shmpost() {
    struct sembuf sem_a;
    sem_a.sem_num = 0;
    sem_a.sem_op = 1;
    sem_a.sem_flg = SEM_UNDO;
    if (semop(shmid, &sem_a, 1) == -1) return false;
    return true;
}

bool shmuget(char *shmp) {
    if (shmdt(shmp) == -1) {
        return false;
    }
    return true;
}

bool shmremove() {
    if (semctl(shmid, 0, IPC_RMID) == -1) return false;
    return true;
}