package monitor

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rjeczalik/notify"
)

// Watch a specific folder for changes
func Watch(folder string) {

	log.Debug("Monitor started on folder " + folder)
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	if err := notify.Watch(folder+"/...", c, notify.All); err != nil {
		log.WithField("error", err).Panic("Failed to watch comic directory")
	}

	defer notify.Stop(c)

	for e := range c {
		log.WithField("event", e).Debug("something moved in the folder")
	}
}
