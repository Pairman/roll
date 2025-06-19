package main

import (
	"os"

	"github.com/akamensky/argparse"
	"pnxlr.eu.org/roll/fs/header"
	"pnxlr.eu.org/roll/net/down"
	"pnxlr.eu.org/roll/net/up"
	"pnxlr.eu.org/roll/util/log"
)

/*

TODO:

- Add command "cloud" with "login", "logout", "list", "down" and "up" sub-commands
- Store files in a dedicated folder in the cloud disk
*/

func main() {
	parser := argparse.NewParser("roll v"+Version, "Share files on the go")
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Print verbosely",
	})
	versionCmd := parser.NewCommand("version", "Print version and exit")
	uploadCmd := parser.NewCommand("up", "Upload a file")
	uploadFile := uploadCmd.File("f", "file", os.O_RDONLY, 0, &argparse.Options{
		Required: true,
		Help:     "Path to the file",
	})
	uploadCompress := uploadCmd.Flag("c", "compress", &argparse.Options{
		Required: false,
		Help:     "Compress the file (does not work with encryption)",
	})
	uploadEncrypt := uploadCmd.Flag("e", "encrypt", &argparse.Options{
		Required: false,
		Help:     "Encrypt the file (does not work with compression)",
	})
	downloadCmd := parser.NewCommand("down", "Download a file")
	downloadId := downloadCmd.String("i", "id", &argparse.Options{
		Required: true,
		Help:     "URL, object ID or resource ID",
	})
	downloadPath := downloadCmd.String("f", "file", &argparse.Options{
		Required: false,
		Help:     "Path to save the file",
	})

	if err := parser.Parse(os.Args); err != nil {
		log.Err(parser.Usage(err))
		os.Exit(1)
	}

	switch {
	case versionCmd.Happened():
		log.Infof("v" + Version)
		os.Exit(0)
	case uploadCmd.Happened():
		if *uploadCompress && *uploadEncrypt {
			log.Errf(
				"Compression and encryption cannot be specified simultaneously")
			os.Exit(1)
		}

		options := &up.UploadOptions{
			Verbose:  *verbose,
			Compress: up.CompressionOptions{On: *uploadCompress},
			Encrypt:  up.EncryptionOptions{On: *uploadEncrypt},
		}
		if options.Compress.On {
			options.Compress.Algo = header.CompressionAlgoZSTD
		} else if options.Encrypt.On {
			options.Encrypt.Algo = header.EncryptionAlgoAES256GCM
		}
		res, err := up.Upload(uploadFile, options)
		if err == nil {
			if res.ShareKey != "" {
				log.Infof("Uploaded to '%v', share key '%v' expires in 1 hour",
					res.URL, res.ShareKey)
			} else {
				log.Infof("Uploaded to '%v'", res.URL)
			}
		} else {
			log.Errf("Upload error: %v", err)
			os.Exit(1)
		}
	case downloadCmd.Happened():
		res, err := down.Download(*downloadId, &down.DownloadOptions{
			Verbose: *verbose, Path: *downloadPath,
		})
		if err == nil {
			log.Infof("Downloaded to '%v'", res.Path)
		} else {
			log.Errf("Download error: %v", err)
			os.Remove(res.Path)
			os.Exit(1)
		}
	}
}
