package up

import (
	"crypto/rand"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/reader"
	"pnxlr.eu.org/roll/net/up/api"
	netUtil "pnxlr.eu.org/roll/net/util"
	"pnxlr.eu.org/roll/util/log"
)

func Upload(file *os.File, options *UploadOptions) (*UploadResult, error) {
	var fname string
	var fsize int
	if finfo, err := file.Stat(); err != nil {
		log.Errf("Upload error: %v", err)
		return nil, err
	} else {
		fname = finfo.Name()
		if options.Verbose {
			fsize = int(finfo.Size())
			log.Infof("Upload size: %vB, name: '%v'.\nUpload start...",
				finfo.Size(), fname)
		}
	}
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	fh := header.NewFileHeaderFromFile(file)
	fh.CompSect.Algo = options.Compress.Algo
	fh.EncSect.Algo = options.Encrypt.Algo
	br := reader.NewBlockReader(file, fsize,
		fh.ToBytes())
	go upload(br, pw, w, options)
	uploader := api.NewRobotUploader()
	req, _ := http.NewRequest("POST", uploader.URL, pr)
	req.Header = uploader.Headers.Clone()
	req.Header.Set("Content-Type", w.FormDataContentType())

	timeStart := time.Now()
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Errf("Upload error: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errf("Upload failed: %v", res.Status)
		return nil, errors.New(res.Status)
	}
	timeElapsed := time.Since(timeStart)

	body, _ := io.ReadAll(res.Body)
	json, err := uploader.Json(body)
	if err != nil {
		log.Errf("Upload error: %v", err)
		return nil, err
	}
	if !uploader.Success(json) {
		log.Errf("Upload failed: %v", json)
	}
	id := uploader.ObjectID(json)
	url := netUtil.ShareURLFromObjectID(uploader.ObjectID(json))
	if options.Verbose {
		log.Infof("Uploaded to '%v', %.2fs elapsed.", url,
			timeElapsed.Seconds())
	}
	return &UploadResult{ObjectID: id, URL: url}, nil
}

func upload(br *reader.BlockReader, pw *io.PipeWriter, w *multipart.Writer,
	options *UploadOptions) {
	defer br.Close()
	defer pw.Close()
	defer w.Close()
	form, err := w.CreateFormFile("file",
		time.Now().Format("20060102150405")+".png")
	if err != nil {
		pw.CloseWithError(err)
		return
	}

	var r io.Reader = br
	var c io.Closer = br
	defer c.Close()
	if options.Compress.On {
		switch options.Compress.Algo {
		case header.CompressionAlgoZSTD:
			log.Infof("Compression enabled with ZSTD")
			zr := reader.NewZSTDEncoder(br)
			r, c = zr, zr
		}
	} else if options.Encrypt.On {
		switch options.Encrypt.Algo {
		case header.EncryptionAlgoAES256GCM:
			log.Infof("Encryption enabled with AES-256-GCM")
			var kiv [32 + 12]byte
			rand.Read(kiv[:])
			log.Infof("AES-256-GCM key and IV: %x", kiv)
			ar := reader.NewAESGCMEncoder(br, kiv[:32], kiv[32:])
			defer ar.Close()
			r, c = ar, ar
		}
	}

	_, err = io.Copy(form, r)
	if err != nil {
		pw.CloseWithError(err)
		return
	}
}
