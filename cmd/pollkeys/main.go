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
	target   string
	keysFile string
)

func init() {
	flag.StringVar(&target, "target", "", "Target file to keep track of")
	flag.StringVar(&keysFile, "path", "~/.ssh/authorized_keys", "File to save newest version of authorized_keys")
	flag.StringVar(&tempFile, "temp", "/tmp/keys.txt", "Temporary save location of downloaded version of authorized_keys")

	flag.Parse()
}

func main() {
	log.Printf("Polling for changes to %s\n", target)
	err := DownloadFile(tempFile, target)
	if err != nil {
		log.Fatalf("Failed to download %s\n", target)
	}
	log.Printf("Successfully downloaded %s\n", target)

	checksum, err := getMD5(tempFile)
	if err != nil {
		log.Fatalf("Failed to get checksum of %s\n", tempFile)
	}
	keysChecksum, err := getMD5(keysFile)
	if err != nil {
		log.Fatalf("Failed to get checksum of %s\n", keysFile)
	}
	if bytes.Compare(checksum[:], keysChecksum[:]) == 0 {
		log.Println("Checksums match, no changes necessary")
		os.Exit(0)
	}
	log.Printf("Changes detected, overriding %s\n", keysFile)
	input, err := ioutil.ReadFile(tempFile)
	if err != nil {
		panic(err)
	}

	// authorized_keys should have 0644 permissions
	err = ioutil.WriteFile(keysFile, input, 0644)
	if err != nil {
		panic(err)
	}
	log.Printf("%s overwritten\n", keysFile)
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

func getMD5(path string) ([16]byte, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return [16]byte{}, err
	}
	return md5.Sum(contents), nil
}
