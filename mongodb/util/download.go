package util

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/terraform/command"
)

var randLock sync.Mutex
var duplicateDownloadLocks map[string]*sync.Mutex
var randShared *rand.Rand
var downloadStorage string

func init() {
	// ensure all generated numbers are unique
	randLock.Lock()
	defer randLock.Unlock()
	if randShared == nil {
		randShared = rand.New(rand.NewSource(
			time.Now().UnixNano() * int64(os.Getpid())))
	}

	// init the mutexes
	duplicateDownloadLocks = make(map[string]*sync.Mutex)

	// ensure we have a downloads directory
	downloadStorage = filepath.Join(command.DefaultDataDir, "downloads")
	if err := os.MkdirAll(downloadStorage, 0750); err != nil {
		log.Fatal(err)
	}
}

// DownloadFile downloads the specified url to a temporary location
func DownloadFile(url string) (destFile *os.File, finalErr error) {
	filename := filepath.Clean(filepath.Base(url))
	dest := filepath.Join(downloadStorage, filename)

	var mtx *sync.Mutex
	var ok bool

	// create a mutex if one for the current key was not already defined
	if mtx, ok = duplicateDownloadLocks[filename]; !ok {
		mtx = &sync.Mutex{}
		duplicateDownloadLocks[filename] = mtx
	}

	// prevent duplicate downloads of the same filename
	mtx.Lock()
	defer mtx.Unlock()

	// return early if the file already exists
	if _, err := os.Stat(dest); os.IsExist(err) {
		destFile, finalErr = os.Open(dest)
		if finalErr != nil {
			finalErr = MergeErrors("os.DownloadFile: could not open existing file "+dest, finalErr, nil)
		}
		return
	}

	// download the file
	httpClient := &http.Client{Timeout: DownloadTimeout}
	resp, err := httpClient.Get(url)
	if err != nil {
		finalErr = MergeErrors("os.DownloadFile: could not download from "+url, err, nil)
		return
	}
	defer LogError(resp.Body.Close)

	// check the server's response, ensure it's 200/OK
	if resp.StatusCode != http.StatusOK {
		finalErr = fmt.Errorf("os.DownloadFile: got bad HTTP status code: %s", resp.Status)
		return
	}

	// copy all data to a temporary file
	log.Printf("[DEBUG] downloading %s", url)
	file, err := ReadAllIntoTempFile(resp.Body, filename)
	if err != nil {
		if file != nil {
			filename = file.Name()
		}
		finalErr = MergeErrors("os.DownloadFile: could not save to temporary file "+filename, err, nil)
		return
	}
	log.Printf("[DEBUG] created local filename: %s", file.Name())

	// Move the file to its final location
	err = os.Rename(file.Name(), dest)
	if err != nil {
		finalErr = MergeErrors("os.DownloadFile: could not move file to its final location "+dest, err, nil)
		return
	}
	log.Printf("[DEBUG] moved downloaded data to its final destination: %s", dest)

	// return the resulting file
	destFile, finalErr = os.Open(dest)
	if finalErr != nil {
		finalErr = MergeErrors("os.DownloadFile: could not open file "+dest, finalErr, nil)
	}
	return
}
