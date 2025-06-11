package down

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/fs/reader"
	fsUtil "pnxlr.eu.org/roll/fs/util"
	netUtil "pnxlr.eu.org/roll/net/util"
	"pnxlr.eu.org/roll/util/log"
)

// TODO: download pause & resume

func Download(id string, options *DownloadOptions) (*DownloadResult, error) {
	// Retrieve information
	client := &http.Client{}
	info, fh, err := retrieveDownloadInfo(client, id)
	if err != nil {
		log.Errf("Download error: %v", err)
		return nil, err
	}

	// Get output path
	outPath, err := fsUtil.ComposePath(options.Path, string(fh.FileSect.Name))
	if err != nil {
		log.Errf("Download error: %v", err)
		return nil, err
	}
	if options.Verbose {
		log.Infof("Download size: %vB, name: %s, output path: %s",
			info.Length-fh.Len(), string(fh.FileSect.Name), outPath)
	}

	// Create temporary file
	tmpPath := outPath + "." + time.Now().Format("150405.000") + ".tmp"
	tmpFile, err := fsUtil.CreateFile(tmpPath, int64(info.Length-fh.Len()))
	if err != nil {
		log.Errf("Download error: %v", err)
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download to temporary file
	if err := downloadFileParallel(client, info.Download, int64(fh.Size),
		int64(info.Length)-1, tmpFile, options.Verbose); err != nil {
		log.Errf("Download error: %v", err)
		return nil, err
	}
	tmpFile.Seek(0, io.SeekStart)

	// Possible decompression and decryption
	isCompressed := fh.CompSect.Algo != header.CompressionAlgoNone
	isEncrypted := fh.EncSect.Algo != header.EncryptionAlgoNone
	if isCompressed || isEncrypted {
		// Create output file
		file, err := fsUtil.CreateFile(outPath, fh.FileSect.FileSize)
		if err != nil {
			log.Errf("Download error: %v", err)
			return nil, err
		}
		defer file.Close()

		// Decompress / decrypt
		if isCompressed {
			switch fh.CompSect.Algo {
			case header.CompressionAlgoZSTD:
				if options.Verbose {
					log.Infoln("Decompressing with ZSTD")
				}
				dec := reader.NewZSTDDecoder(tmpFile)
				defer dec.Close()
				if _, err := io.Copy(file, dec); err != nil {
					log.Errf("Download error: %v", err)
					return nil, err
				}
			default:
				err := fmt.Errorf("unknown compression algorithm: %v",
					fh.CompSect.Algo)
				log.Errf("Download error: %v", err)
				return nil, err
			}
		} else if isEncrypted {
			switch fh.EncSect.Algo {
			case header.EncryptionAlgoAES256GCM:
				if options.Verbose {
					log.Infoln("Decrypting with AES-256-GCM")
				}
				var kivStr string
				log.Info("Input AES-256-GCM key and IV: ")
				fmt.Scanln(&kivStr)
				kiv, err := hex.DecodeString(kivStr)
				if err != nil {
					log.Errf("Download error: %v", err)
					return nil, err
				} else if len(kiv) != 32+12 {
					err := errors.New("invalid key and IV")
					log.Errf("Download error: %v", err)
					return nil, err
				}
				dec := reader.NewAESGCMDecoder(tmpFile, kiv[:32], kiv[32:])
				defer dec.Close()
				if _, err := io.Copy(file, dec); err != nil {
					log.Errf("Download error: %v", err)
					return nil, err
				}
			default:
				err := fmt.Errorf("unknown decryption algorithm: %v",
					fh.EncSect.Algo)
				log.Errf("Download error: %v", err)
				return nil, err
			}
		}
	} else {
		// Move temporary file to output path
		fsUtil.MoveFile(tmpFile.Name(), outPath)
	}

	// Checksum
	file, err := os.Open(outPath)
	if err != nil {
		log.Errf("Download error: %v", err)
		return nil, err
	}
	defer file.Close()
	if eq, hash, err := fh.HashSect.Verify(file); err != nil {
		log.Errf("Download error: %v", err)
		return nil, err
	} else if !eq {
		err := fmt.Errorf("hash mismatch: %x, %x", fh.HashSect.Hash, hash)
		log.Errf("Download error: %v", err)
		return nil, err
	}
	fsUtil.SetFileMTime(file, fh.FileSect.Time)

	return &DownloadResult{Path: outPath}, nil
}

func downloadChunk(client *http.Client, url string,
	start, end int) (io.ReadCloser, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = netUtil.GlobalHeader.Clone()
	req.Header.Set("Referer", "https://mooc1.chaoxing.com")
	req.Header.Set("Range", "bytes="+strconv.Itoa(start)+"-"+strconv.Itoa(end))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != 200 && res.StatusCode != 206 {
		defer res.Body.Close()
		return nil, errors.New(res.Status)
	}
	return res.Body, nil
}

func downloadChunkBytes(client *http.Client, url string,
	start, end int) ([]byte, error) {
	b, err := downloadChunk(client, url, start, end)
	if err != nil {
		return nil, err
	}
	defer b.Close()
	return io.ReadAll(b)
}

func downloadFileHeader(client *http.Client,
	url string) (*header.FileHeader, error) {
	p, err := downloadChunkBytes(client, url, header.PNGSectDataSize,
		header.PNGSectDataSize+4-1)
	if err != nil {
		return nil, err
	}
	p, err = downloadChunkBytes(client, url, 0,
		int(fsUtil.LiteralFromBytes[int32](p))-1)
	if err != nil {
		return nil, err
	}
	fh := &header.FileHeader{}
	return fh, fh.FromBytes(p)
}

func retrieveDownloadInfo(client *http.Client,
	id string) (*netUtil.CloudfileStatusJson, *header.FileHeader, error) {
	id, err := netUtil.ObjectIDFromURL(id)
	if err != nil {
		return nil, nil, err
	}
	json, err := netUtil.ObjectIDToStatus(id)
	if err != nil {
		return nil, nil, err
	}
	fh, err := downloadFileHeader(client, json.Download)
	return json, fh, err
}

func downloadChunkFile(client *http.Client, url string,
	start, end int64, file *os.File, offs int64, mu *sync.Mutex) error {
	rc, err := downloadChunk(client, url, int(start), int(end))
	if err != nil {
		return err
	}
	defer rc.Close()
	mu.Lock()
	defer mu.Unlock()
	if _, err := file.Seek(start+offs, io.SeekStart); err != nil {
		return err
	} else if _, err = io.CopyN(file, rc, end-start+1); err == io.EOF {
		return nil
	} else {
		return err
	}
}

func downloadFileParallel(client *http.Client, url string,
	start, end int64, file *os.File, verbose bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mu, wg := &sync.Mutex{}, &sync.WaitGroup{}
	type Task struct{ start, end int64 }
	chTsk, chErr := make(chan Task), make(chan error, 1)
	chPrg := make(chan int64)

	// Workers
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-chTsk:
					if !ok {
						return
					} else if err := downloadChunkFile(client, url,
						task.start, task.end, file, -start, mu); err != nil {
						select {
						case chErr <- err:
							cancel()
						default:
						}
						return
					}
					if verbose {
						chPrg <- task.end - task.start + 1
					}
				}
			}
		}()
	}

	// Scheduler
	go func() {
		const chunkSize = 4 * fsUtil.MiB
		defer close(chTsk)
		for offs := start; offs <= end; offs += chunkSize {
			select {
			case <-ctx.Done():
				return
			case chTsk <- Task{start: offs, end: min(offs+chunkSize-1, end)}:
			}
		}
	}()

	// Print progress
	if verbose {
		go func() {
			var totalDown, lastDown, lastPrg int64
			for {
				select {
				case <-ctx.Done():
					cancel()
				case down := <-chPrg:
					totalDown += down
					if prg := 100 * totalDown / (end + 1 - start); lastPrg < prg &&
						(lastDown+fsUtil.MiB < totalDown || prg == 100) {
						lastDown, lastPrg = totalDown, prg
						log.Infof("Downloaded: %d%%", prg)
					}
				}
			}
		}()
	}

	wg.Wait()
	select {
	case err := <-chErr:
		return err
	default:
		return nil
	}
}
