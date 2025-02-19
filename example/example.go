package main

import (
	"fmt"
	"time"

	"github.com/dongs0104/activity-tracker/internal/pkg/logging"
	"github.com/dongs0104/activity-tracker/pkg/tracker"
)

func main() {

	logger := logging.New()

	heartbeatInterval := 60 //value always in seconds
	workerInterval := 4     //seconds

	activityTracker := &tracker.Instance{
		HeartbeatInterval: heartbeatInterval,
		WorkerInterval:    workerInterval,
		LogLevel:          logging.Info,
	}

	//This starts the tracker for all handlers. It gives you a channel
	//which you can listen to for heartbeat objects
	heartbeatCh := activityTracker.Start()

	// //if you only want to track certain handlers, you can use StartWithhandlers
	// heartbeatCh := activityTracker.StartWithHandlers(handler.MouseClickHandler(),
	// 	handler.MouseCursorHandler())

	loopTime := 2
	timeToKill := time.NewTicker(time.Second * time.Duration((heartbeatInterval*loopTime)+1))

	logger.Infof("starting activity tracker with %vs heartbeat and %vs worker interval ...", heartbeatInterval, workerInterval)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.WasAnyActivity {
				logger.Infof("no activity detected in the last %v seconds\n\n\n", int(heartbeatInterval))
			} else {
				logger.Infof("activity detected in the last %v seconds.", int(heartbeatInterval))
				logger.Infof("Activity type:\n")
				for activityType, times := range heartbeat.ActivityMap {
					logger.Infof("activityType : %v times: %v\n", activityType, len(times))
				}
				fmt.Printf("\n\n\n")
			}
		case <-timeToKill.C:
			logger.Infof("time to kill app")
			activityTracker.Quit()
			return
		}
	}
}
