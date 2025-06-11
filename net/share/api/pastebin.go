package api

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// CentOS Pastebin Service
const PastebinURL = "https://paste.centos.org"

// From fpaste https://pagure.io/fpaste/blob/main/f/fpaste#_32
const PastebinAPIKey = "5uZ30dTZE1a5V0WYhNwcMddBRDpk6UzuzMu-APKM38iMHacxdA0n4vCqA34avNyt"

func PastebinShareCreate(title, text string) (string, error) {
	form := url.Values{}
	form.Set("text", text)
	form.Set("title", title)
	form.Set("expire", "60")
	res, err := http.PostForm(PastebinURL+"/api/create?apikey="+url.QueryEscape(PastebinAPIKey), form)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	p, _ := io.ReadAll(res.Body)
	body := string(p)

	if strings.HasPrefix(body, "Error:") {
		return "", errors.New(body)
	} else if p := strings.Split(strings.TrimSpace(body), "/"); len(p) < 1 {
		return "", errors.New("unexpected response: " + body)
	} else {
		return p[len(p)-1], nil
	}
}

func PastebinShareGet(id string) (string, error) {
	res, err := http.Get(PastebinURL + "/view/raw/" + id)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}

	if p, err := io.ReadAll(res.Body); err != nil {
		return "", err
	} else {
		return string(p), nil
	}
}
