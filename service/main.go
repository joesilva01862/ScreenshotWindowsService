/*
     This is the main func(). 
     
     This application defines a Windows service that takes a screenshot
     every n minutes defined at the command line, and sends the screenshot
     to a server (currently hardcoded in the program that invokes the web service).
     
     How it works:
     1. This service is installed, and is recognized by Windows.
     2. Every n minutes (defined at installation time at the command line) this program
        wakes up and invokes a Windows API CreateProcessAsUser,
        which allows us to invoke the screenshot.exe program to access Windows GUI.
     2. That screenshot.exe program takes and send a screenshot to the server.
        
     Based on the project found here: https://github.com/judwhite/go-svc
     
     To build:
     go build -o service.exe main.go loop.go process.go
*/   

package main

import (
    "os"
    "fmt"

    "github.com/judwhite/go-svc/svc"
    "golang.org/x/sys/windows/svc/eventlog"
    "golang.org/x/sys/windows/registry"
)

// implements svc.Service
type program struct {
    LogFile *os.File
    svr     *server
}

// global vars
var (
   evtlog * eventlog.Log
   SERVICE_NAME = "Screenshot Service"
   interval_mins uint64
   pgm_fullpath string
)

func main() {
    // prepare to write to Windows Event Log
    var err error

    k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Microservice\Netpipe`, registry.QUERY_VALUE)
    check(err)
    
    interval_mins, _, err = k.GetIntegerValue("INTERVAL_MINUTES")
    check(err)
    pgm_fullpath, _, err = k.GetStringValue("PGM_TO_INVOKE")
    check(err)
    
    // open Windows event logger
    evtlog, err = eventlog.Open(SERVICE_NAME)
    check(err)
    
    // register our service name to Windows event logger
    eventlog.InstallAsEventCreate(SERVICE_NAME, eventlog.Info | eventlog.Warning | eventlog.Error)
    defer evtlog.Close()
    
    prg := program {
        svr: &server{},
    }
    
    // call svc.Run to start your program/service
    // svc.Run will call Init, Start, and Stop
    if err := svc.Run(&prg); err != nil {
        evtlog.Error(1, "Service error while trying to start running.")
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (p *program) Init(env svc.Environment) error {
    evtlog.Info(1, fmt.Sprintf("is win service? %v\n", env.IsWindowsService()))

/*
    // write to "example.log" when running as a Windows Service
    if env.IsWindowsService() {
        dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            return err
        }

        logPath := filepath.Join(dir, "example.log")

        f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            return err
        }

        p.LogFile = f

        log.SetOutput(f)
    }
*/
    return nil
}

func (p *program) Start() error {
    evtlog.Info(1, "Starting...\n")
    
    // take first shot and start the loop
    p.svr.prepareStart()
    
    go p.svr.start()
    return nil
}

func (p *program) Stop() error {
    evtlog.Info(1, "Stopping...\n")
    if err := p.svr.stop(); err != nil {
        return err
    }
    evtlog.Info(1, "Stopped.\n")
    return nil
}
