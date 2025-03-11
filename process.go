package main

import (
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func processJob(job *Job) { // processes job asynchronously
	var wg sync.WaitGroup
	errorOccurred := false
	var mu sync.Mutex 

	for i, visit := range job.Visits {
		wg.Add(1)
		go func(i int, visit Visit) {
			defer wg.Done()

			if !StoreExists(visit.StoreID) {
				mu.Lock()
				job.Errors = append(job.Errors, JobError{
					StoreID: visit.StoreID,
					Error:   "store id not found",
				})
				mu.Unlock()
				errorOccurred = true
				return
			}

			var imgWg sync.WaitGroup
			results := make([]ImageResult, len(visit.ImageURLs))
			for j, url := range visit.ImageURLs {
				imgWg.Add(1)
				go func(j int, url string) {
					defer imgWg.Done()
					res, err := processImage(url)
					if err != nil {
						results[j] = ImageResult{
							URL:   url,
							Error: err.Error(),
						}
						mu.Lock()
						job.Errors = append(job.Errors, JobError{
							StoreID: visit.StoreID,
							Error:   err.Error(),
						})
						mu.Unlock()
						errorOccurred = true
					} else {
						results[j] = ImageResult{
							URL:       url,
							Perimeter: res,
						}
					}
				}(j, url)
			}
			imgWg.Wait()

			mu.Lock()
			job.Visits[i].Images = results
			mu.Unlock()
		}(i, visit)
	}
	wg.Wait()

	if errorOccurred {
		job.Status = StatusFailed
	} else {
		job.Status = StatusCompleted
	}
	updateJob(job)
}

func processImage(url string) (int, error) { // downloads an image, calculates its perimeter, and simulates GPU delay
	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.New("failed to download image")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to download image: " + resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, errors.New("failed to decode image")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	perimeter := 2 * (width + height)

	sleepDuration := time.Duration(100+rand.Intn(300)) * time.Millisecond
	time.Sleep(sleepDuration)

	return perimeter, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
