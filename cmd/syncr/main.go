package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chlunde/syncr"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %s\n", err)
	}

	var syncrs []*syncr.Syncr
	for _, conf := range config {
		for _, host := range conf.Hosts {
			s, err := syncr.NewSyncr(conf.Source, fmt.Sprintf("%s:%s", host, conf.Destination))

			if s != nil {
				// prefer to show error in UI if s exists (s.Error will be shown)
				syncrs = append(syncrs, s)
			} else if err != nil {
				log.Println("Failed to create syncr for %s->%s (%s): %s", conf.Source, conf.Destination, host, err)
			}
		}
	}

	buf := &bytes.Buffer{}
	for {
		time.Sleep(100 * time.Millisecond)

		// Write to a buffer instead of stdout to prevent any flickering
		buf.Reset()
		anyAlive := Display(buf, syncrs)
		os.Stdout.Write(buf.Bytes())

		if !anyAlive {
			break
		}
	}
}
