// main
package recording

import (
	"fmt"
	"strconv"

	"io"
	"log"

	"github.com/grafov/m3u8"

	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

var client = &http.Client{}

func Record(inUrl string, maxWidth int, maxHeight int, segmentsPath string) error {
	if !strings.HasPrefix(inUrl, "http") {
		return fmt.Errorf("cms17> Playlist URL must begin with http/https")
	}

	theURL, err := url.Parse(inUrl)
	if err != nil {
		return fmt.Errorf("cms18> %w", err)
	}

	if _, err := os.Stat(segmentsPath); os.IsNotExist(err) {
		err = os.Mkdir(segmentsPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cms19> %w", err)
		}
	}

	var highestId = -1
	for {
		err := getPlaylist(theURL, maxWidth, maxHeight, segmentsPath, &highestId)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(2) * time.Second)
	}
}

func getPlaylist(u *url.URL, maxWidth int, maxHeight int, segmentsPath string, highestId *int) error {

	content, err := getContent(u)
	if err != nil {
		return fmt.Errorf("cms9> %w", err)
	}

	playlist, listType, err := m3u8.DecodeFrom(*content, true)
	if err != nil {
		return fmt.Errorf("cms10> %w", err)
	}
	(*content).Close()

	if listType != m3u8.MEDIA && listType != m3u8.MASTER {
		return fmt.Errorf("cms11> Not a valid playlist")
	}

	if listType == m3u8.MASTER {
		masterpl := playlist.(*m3u8.MasterPlaylist)
		var variant *m3u8.Variant
		for _, v := range masterpl.Variants {
			resolution := strings.Split(v.VariantParams.Resolution, "x")
			if len(resolution) == 2 {
				width, err := strconv.Atoi(resolution[0])
				if err == nil {
					if width > maxWidth {
						break
					}

					variant = v
				}
			}
		}

		if variant == nil {
			variant = masterpl.Variants[len(masterpl.Variants)-1]
		}

		if variant != nil {
			msURL, err := absolutize(variant.URI, u)
			if err != nil {
				return fmt.Errorf("cms12> %w", err)
			}

			err = getPlaylist(msURL, maxWidth, maxHeight, segmentsPath, highestId)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if listType == m3u8.MEDIA {
		mediapl := playlist.(*m3u8.MediaPlaylist)
		for _, segment := range mediapl.Segments {
			if segment != nil {

				msURL, err := absolutize(segment.URI, u)
				if err != nil {
					return fmt.Errorf("cms15> %w", err)
				}

				id := extractIdFromUrl(*msURL)
				if id > *highestId {
					err := download(msURL, segmentsPath, strconv.Itoa(id))
					if err != nil {
						return fmt.Errorf("cms16> %w", err)
					}
					*highestId = id
				}

			}
		}

		time.Sleep(time.Duration(int64(mediapl.TargetDuration)) * time.Second)
	}

	return nil
}

func extractIdFromUrl(u url.URL) int {
	base := path.Base(u.String())
	parts := strings.Split(base, "_")
	idAndExtension := parts[len(parts)-1]
	id := strings.Split(string(idAndExtension), ".")[0]

	v, err := strconv.Atoi(id)
	if err != nil {
		return -1
	}
	return v
}

func getContent(u *url.URL) (*io.ReadCloser, error) {
	var USER_AGENT string

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("cms1> %w", err)
	}

	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cms2> %w", err)
	}

	if resp.StatusCode != 200 {
		return &resp.Body, fmt.Errorf("cms400> Received HTTP %v for %v", resp.StatusCode, u.String())
	}

	return &resp.Body, nil
}

func absolutize(rawurl string, u *url.URL) (uri *url.URL, err error) {

	suburl := rawurl
	uri, err = u.Parse(suburl)
	if err != nil {
		return
	}

	if rawurl == u.String() {
		return
	}

	if !uri.IsAbs() { // relative URI
		if rawurl[0] == '/' { // from the root
			suburl = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, rawurl)
		} else { // last element
			splitted := strings.Split(u.String(), "/")
			splitted[len(splitted)-1] = rawurl

			suburl = strings.Join(splitted, "/")
		}
	}

	suburl, err = url.QueryUnescape(suburl)
	if err != nil {
		return
	}

	uri, err = u.Parse(suburl)
	if err != nil {
		return
	}

	return
}

func writePlaylist(u *url.URL, mpl m3u8.Playlist, outPath string) error {
	fileName := path.Base(u.Path)
	out, err := os.Create(outPath + fileName)
	if err != nil {
		return fmt.Errorf("cms3> %w", err)
	}
	defer out.Close()

	_, err = mpl.Encode().WriteTo(out)
	if err != nil {
		return fmt.Errorf("cms4> %w", err)
	}

	return nil
}

func download(u *url.URL, outPath string, id string) error {
	fileName := path.Base(id + ".ts")

	out, err := os.Create(outPath + fileName)
	if err != nil {
		log.Println("cms5> " + err.Error())
	}
	defer out.Close()

	content, err := getContent(u)
	if content != nil {
		defer (*content).Close()
	}
	if err != nil {
		return fmt.Errorf("cms6> %w", err)
	}

	_, err = io.Copy(out, *content)
	if err != nil {
		return fmt.Errorf("cms7> Failed to download "+fileName+" %w", err)
	}

	return nil
}
