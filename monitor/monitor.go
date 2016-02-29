package monitor

import (
	"os"
	"path/filepath"
	"time"

	//"comiclaser/lzmadec"
	mdl "comiclaser/model"

	log "github.com/Sirupsen/logrus"
	"github.com/rjeczalik/notify"
)

// Watch a specific folder for changes
func Watch(folder string) {

	log.Debug("Monitor started on folder " + folder)
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)
	timer := time.NewTimer(time.Second)

	if err := notify.Watch(folder+"/...", c, notify.All); err != nil {
		log.WithField("error", err).Panic("Failed to watch comic directory")
	}

	defer notify.Stop(c)

	for e := range c {
		log.WithField("event", e).Debug("folder event")
		timer.Stop()
		timer = time.NewTimer(time.Second * 5)
		go func() {
			<-timer.C
			filepath.Walk(folder, visit) // timer expired
		}()

		/*if e.Event() == notify.Create {

			/*dir, file := filepath.Split(e.Path())
			log.WithFields(log.Fields{
				"path": e.Path(),
				"dir":  dir,
				"file": file,
				"loc":  "./" + folder + "/" + file,
			}).Debug("Opening comic file")
			a, err := lzmadec.NewArchive("comics\\A-Force 003 (2015) (Digital) (Zone-Empire).cbr")
			if err != nil {
				log.WithField("error", err).Error("Failed to open archive")
			} else {
				for _, f := range a.Entries {
					log.Debugf("name: %s, size: %d", f.Path, f.Size)
				}
			}

		}*/
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		log.WithFields(log.Fields{
			"path": path,
			"f":    f,
		}).Info("Visited")

		mdl.CreateComic(mdl.Comic{
			Path:   path,
			Series: "test",
		})
	}

	return nil
}
