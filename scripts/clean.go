package scripts

import (
	"io/ioutil"
	"jazzine/m3u8recorder/libs/logging"
	"jazzine/m3u8recorder/modules"
	"log"
)

func CleanSegments(segmentsPath string) {
	files, err := ioutil.ReadDir(segmentsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			err := modules.Merge(segmentsPath+f.Name()+"/", "./videos/", f.Name()+".ts")
			logging.Log().
				WithError(err).
				Error("cleaning err: failed merging segments of: " + f.Name())
		}
	}
}
