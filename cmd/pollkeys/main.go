package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	tempFile string
	target string
	outFile string
)

func main () {
	flag.StringVar(&target, "target", "", "Target file to keep track of")
	flag.StringVar(&outFile, "out", "~/.ssh/authorized_keys", "File to save newest version of authorized_keys")
	flag.StringVar(&tempFile, "temp", "/tmp/keys.txt", "Temporary save location of downloaded version of authorized_keys")
	
	flag.Parse()
	
	log.Printf("Polling for changes to %s\n", target)
	err := DownloadFile(tempFile, target)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully downloaded file")
	contents, err := ioutil.ReadFile(tempFile)
	if err != nil {
		panic(err)
	}
	checksum := md5.Sum(contents)
	currentChecksum := currentMD5()
	if bytes.Compare(checksum[:], currentChecksum[:]) == 0 {
		log.Println("Checksums match, no changes necessary")
	} else {
		log.Printf("Changes detected, overriding %s\n", outFile)
		input, err := ioutil.ReadFile(tempFile)
		if err != nil {
			panic(err)
		}

		// authorized_keys should have 0644 permissions
		err = ioutil.WriteFile(outFile, input, 0644)
		if err != nil {
			panic(err)
		}
		log.Printf("%s overwritten\n", outFile)
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func currentMD5() [16]byte {
	contents, err := ioutil.ReadFile(outFile)
	if err != nil {
		panic(err)
	}
	return md5.Sum(contents)
}
