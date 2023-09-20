package modules

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Merge(inFolder string, outFolder string, filename string) error {
	files, err := ioutil.ReadDir(inFolder)
	if err != nil {
		log.Println(err)
		return err
	}

	defer os.RemoveAll(inFolder)

	if len(files) < 200 {
		return fmt.Errorf("not enough files to merge")
	}

	outputPath := outFolder + filename
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer outFile.Close()

	for _, file := range files {
		if file.Name() == ".DS_Store" || file.Name() == "playlist.m3u8" {
			continue
		}

		filePath := inFolder + file.Name()

		f, err := os.ReadFile(filePath)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = outFile.Write(f)
		if err != nil {
			log.Println(err)
			continue
		}

		os.Remove(filePath)
	}

	return nil
}
