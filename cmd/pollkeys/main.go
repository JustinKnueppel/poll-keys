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
	"strconv"
)

var (
	tempFile    string
	remoteUrl   string
	localTarget string
	permissions os.FileMode
)

func init() {
	var permissionsString string
	flag.StringVar(&remoteUrl, "remote", "", "Remote url of file to synchronize")
	flag.StringVar(&localTarget, "path", "", "Local path of file to synchronize")
	flag.StringVar(&permissionsString, "perms", "0644", "Permissions for saved file")
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

	p, err := strconv.ParseUint(permissionsString, 10, 32)
	if err != nil {
		fmt.Println("Failed to convert permissions to file mode")
		os.Exit(1)
	}
	permissions = os.FileMode(p)
}

func main() {
	Synchronize(localTarget, remoteUrl)
}

func Synchronize(path, remote string) error {
	log.Printf("Polling for changes to %s\n", remote)
	err := DownloadFile(tempFile, remote)
	if err != nil {
		log.Printf("Failed to download %s\n", remote)
		return err
	}
	log.Printf("Successfully downloaded %s\n", remote)

	checksum, err := getMD5(tempFile)
	if err != nil {
		log.Printf("Failed to get checksum of %s\n", tempFile)
		return err
	}
	keysChecksum, err := getMD5(path)
	if err != nil {
		log.Printf("Failed to get checksum of %s\n", path)
		return err
	}
	if bytes.Equal(checksum[:], keysChecksum[:]) {
		log.Println("Checksums match, no changes necessary")
		return nil
	}
	log.Printf("Changes detected, overriding %s\n", path)
	input, err := ioutil.ReadFile(tempFile)
	if err != nil {
		return err
	}

	// authorized_keys should have 0644 permissions
	err = ioutil.WriteFile(path, input, permissions)
	if err != nil {
		return err
	}
	log.Printf("%s overwritten\n", path)
	return nil
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
