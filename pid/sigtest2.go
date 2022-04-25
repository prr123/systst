package main

/*
#include <stdio.h>
#include <signal.h>
#include <string.h>

struct sigaction old_action;
void handler(int signum, siginfo_t *info, void *context) {
    printf("Sent by %d\n", info->si_pid);
}

void test() {
    struct sigaction action;
    sigaction(SIGUSR1, NULL, &action);
    memset(&action, 0, sizeof action);
    sigfillset(&action.sa_mask);
    action.sa_sigaction = handler;
    action.sa_flags = SA_RESTART | SA_NOCLDSTOP | SA_SIGINFO | SA_ONSTACK;
    sigaction(SIGUSR1, &action, &old_action);
}
*/
import "C"

import (
    "os"
    "syscall"
    "time"
)

func main() {
    C.test()
    pid := os.Getpid()
    for {
        syscall.Kill(pid, syscall.SIGUSR1)
        time.Sleep(time.Second)
    }
}
