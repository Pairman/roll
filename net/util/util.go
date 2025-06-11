package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	netUrl "net/url"
	"path"
	"strconv"
	"strings"
)

func ShareURLFromObjectID(id string) string {
	return "https://pan-yz.chaoxing.com/external/m/file/" + id
}

func ObjectIDFromURL(url string) (string, error) {
	u, err := netUrl.Parse(url)
	if err != nil {
		return "", err
	}
	id := strings.SplitN(path.Base(u.Path), ".", 1)[0]

	if _, err := strconv.Atoi(id); err == nil {
		req, _ := http.NewRequest("GET", ShareURLFromObjectID(id), nil)
		req.Header = GlobalHeader.Clone()

		res, err := (&http.Client{}).Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return "", errors.New(res.Status)
		}
		data, _ := io.ReadAll(res.Body)
		if i := bytes.Index(data, []byte("'objectId': '")) + 13; i != 13-1 &&
			len(data) >= i+32 {
			return string(data[i : i+32]), nil
		}
	} else if len(id) == 32 {
		return id, nil
	}
	return "", errors.New("no valid object ID found")
}

type CloudfileStatusJson struct {
	Download string `json:"download"`
	Filename string `json:"filename"`
	CRC      string `json:"crc"`
	Length   int    `json:"length"`
	HTTP     string `json:"http"`
	ObjectID string `json:"objectid"`
	Key      string `json:"key"`
}

func StatusURLFromObjectID(id string) string {
	return "https://mooc1.chaoxing.com/ananas/status/" + id
}

func ObjectIDToStatus(id string) (*CloudfileStatusJson, error) {
	req, _ := http.NewRequest("GET", StatusURLFromObjectID(id), nil)
	req.Header = GlobalHeader.Clone()
	req.Header.Set("Referer", "https://mooc1.chaoxing.com")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}
	body, _ := io.ReadAll(res.Body)
	data := &CloudfileStatusJson{}
	return data, json.Unmarshal(body, data)
}
