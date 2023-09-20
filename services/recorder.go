package services

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/chromedp/chromedp"

	"jazzine/m3u8recorder/datasource"
	"jazzine/m3u8recorder/libs/logging"
	"jazzine/m3u8recorder/libs/recording"
	"jazzine/m3u8recorder/model"
	"jazzine/m3u8recorder/modules"
	"jazzine/m3u8recorder/scripts"
)

const segmentsRootDir = "./segments/"
const baseUrl = "https://your-site/"

func Record(ctx context.Context, dao *datasource.Dao) {
	createDirIfNotExist(segmentsRootDir)

	staticSlugs := []string{}

	scripts.CleanSegments(segmentsRootDir)

	slugs := []string{}
	currentlyRecording := []string{}
	for {
		slugs = append(slugs, staticSlugs...)

		if len(slugs) <= len(currentlyRecording) {
			time.Sleep(time.Duration(5) * time.Minute)
		}

		for _, slug := range slugs {
			if Contains(currentlyRecording, slug) != -1 {
				continue
			}

			room, err := getRoomDossier(baseUrl + slug)
			if err != nil {
				logging.Log().
					WithError(err).
					Warn("failed getting room dossier: " + slug)
			} else if room != nil && room.RoomStatus == "public" {
				currentlyRecording = append(currentlyRecording, slug)
				id := getNewRecordingId(slug)

				go func(roomSlug string, recordingId string) {
					defer func() {
						currentlyRecording, err = RemoveFromStringSlice(currentlyRecording, roomSlug)
						if err != nil {
							logging.Log().
								WithError(err).
								Error("failed removing slug from currentlyRecording slice: " + roomSlug)
						}
					}()

					logging.Log().
						Info("start recording: " + recordingId)

					segmentsFolderPath := segmentsRootDir + recordingId + "/"

					_ = recording.Record(room.HlsSource, 2000, 720, segmentsFolderPath)

					err = modules.Merge(segmentsFolderPath, "./videos/", recordingId+".ts")
					if err != nil {
						logging.Log().
							WithError(err).
							Error("failed merging segments of: " + recordingId)
						return
					}

					// DB
					account, err := dao.FindAccountOrCreate(ctx, roomSlug)
					if err != nil || account == nil {
						logging.Log().
							WithError(err).
							Error("error finding or creating account: " + roomSlug)
						return
					}

					record := model.Record{
						AccountID: account.ID,
						URL:       recordingId,
					}
					err = dao.SaveRecord(ctx, record)
					if err != nil {
						logging.Log().
							WithError(err).
							Error("error saving record: " + recordingId)
						return
					}
				}(slug, id)
			}

			time.Sleep(time.Duration(30) * time.Second)
		}
	}
}

func getRoomDossier(url string) (*modules.RoomDossier, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var initialRoomDossierJson string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.EvaluateAsDevTools("window.initialRoomDossier", &initialRoomDossierJson, chromedp.EvalAsValue),
	)
	if err != nil {
		return nil, err
	}

	var initialRoomDossier modules.RoomDossier
	err = json.Unmarshal([]byte(initialRoomDossierJson), &initialRoomDossier)
	if err != nil {
		return nil, err
	}

	return &initialRoomDossier, nil
}

func RemoveFromStringSlice(slice []string, needle string) ([]string, error) {
	index := Contains(slice, needle)
	if index == -1 {
		return slice, errors.New("element to remove not in the slice")
	}

	return append(slice[:index], slice[index+1:]...), nil
}

func getNewRecordingId(slug string) string {
	now := time.Now()
	nowToString := now.Format("2006_01_02_15_59")
	return slug + "_" + nowToString
}

func Contains(slice []string, needle string) int {
	for i, element := range slice {
		if needle == element {
			return i
		}
	}
	return -1
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
