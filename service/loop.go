/*
 *   This class defines how the service reacts to 
 *   the Windows service manager.
 * 
 *   It also defines the work provided by the service, which 
 *   is a loop that wakes up every often to take a screenshot
 *   and go back to sleep.
 *
*/

package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    ticker = time.NewTicker(15 * time.Second)
    servicePaused = false
)

type server struct {
    data chan int
    exit chan struct{}
    wg   sync.WaitGroup
}

func (s *server) start() {
    s.data = make(chan int)
    s.exit = make(chan struct{})
    servicePaused = false
}

func (s *server) stop() error {
    close(s.exit)
    servicePaused = true
    return nil
}

func (s *server) prepareStart() {
    ticker = time.NewTicker(time.Duration(interval_mins * 60) * time.Second)
    s.takeShot()
    servicePaused = false
    go s.startLoop()
}

func (s *server) startLoop() {
    for {
        select {
            case <-ticker.C:
                if !servicePaused {
                    s.takeShot()
                } else {
                    // do nothing
                }
                
        }
    }
}

func (s *server) takeShot() {
    evtlog.Info(1, fmt.Sprintf("Program to be invoked %s.", pgm_fullpath) )

    err := StartProcessAsCurrentUser( pgm_fullpath, pgm_fullpath, "" )
    if err != nil {
        evtlog.Info(1, "An error happened trying to take a screenshot")
        return
    }
    evtlog.Info(1, "A screenshot has been taken.")
}

