package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	tempFile    string
	remoteUrl   string
	localTarget string
)

func init() {
	flag.StringVar(&remoteUrl, "remote", "", "Remote url of file to synchronize")
	flag.StringVar(&localTarget, "path", "", "Local path of file to synchronize")
	flag.StringVar(&tempFile, "temp", "/tmp/keys.txt", "Temporary save location of downloaded file")

	flag.Parse()

	if remoteUrl == "" {
		fmt.Println("Remote url required. Use -remote <url> to specify a remote target")
		os.Exit(1)
	}

	if localTarget == "" {
		fmt.Println("Local file path is required. Use -path <path/to/file> to specify a local file location")
		os.Exit(1)
	}
}

func main() {
	log.Printf("Polling for changes to %s\n", remoteUrl)
	err := DownloadFile(tempFile, remoteUrl)
	if err != nil {
		log.Fatalf("Failed to download %s\n", remoteUrl)
	}
	log.Printf("Successfully downloaded %s\n", remoteUrl)

	checksum, err := getMD5(tempFile)
	if err != nil {
		log.Fatalf("Failed to get checksum of %s\n", tempFile)
	}
	keysChecksum, err := getMD5(localTarget)
	if err != nil {
		log.Fatalf("Failed to get checksum of %s\n", localTarget)
	}
	if bytes.Compare(checksum[:], keysChecksum[:]) == 0 {
		log.Println("Checksums match, no changes necessary")
		os.Exit(0)
	}
	log.Printf("Changes detected, overriding %s\n", localTarget)
	input, err := ioutil.ReadFile(tempFile)
	if err != nil {
		panic(err)
	}

	// authorized_keys should have 0644 permissions
	err = ioutil.WriteFile(localTarget, input, 0644)
	if err != nil {
		panic(err)
	}
	log.Printf("%s overwritten\n", localTarget)
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
