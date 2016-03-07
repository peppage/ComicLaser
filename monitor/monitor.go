package monitor

import (
	"os"
	"path/filepath"
	"time"

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
		switch e.Event() {
		case notify.Create:
			log.WithField("event", e).Debug("folder create event")
			timer.Stop()
			timer = time.NewTimer(time.Second * 5)
			go func() {
				<-timer.C
				filepath.Walk(folder, visit) // timer expired
				mdl.DbUpdated()
			}()
			break
		case notify.Remove:
			log.WithField("event", e).Debug("folder remove event")

			break
		}

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

		c, err := mdl.CreateComic(path)
		if err != nil {
			log.WithError(err).Error("Cannot create comic")
			return err
		}

		c.FileName = f.Name()
		c.Size = f.Size()

		err = mdl.SaveComic(c)
		if err != nil {
			log.WithError(err).Error("Failed saving comic to DB")
		}
	}

	return nil
}

// Update goes through the folder and adds all found comics
func Update(folder string) {
	filepath.Walk(folder, visit)
}

// Remove looks at all files in DB and removes if gone
func Remove(folder string) {
	comics, err := mdl.GetAllComics()
	if err != nil {
		log.WithError(err).Error("Unable to get comics from DB")
		return
	}

	for _, c := range *comics {
		go func(c mdl.Comic) {
			if _, err := os.Stat(c.Path); os.IsNotExist(err) {
				// path/to/whatever does not exist
				mdl.RemoveComic(c.ID)
			}
		}(c)
	}
}
