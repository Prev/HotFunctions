package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/client"
)

var cli *client.Client
var cachedImages map[string]int64
var cachedImagesMutex *sync.Mutex
var imageIsBuilding map[string]bool
var imageIsBuildingMutex *sync.Mutex
var logger *log.Logger

const IMAGE_CACHE_NUM = 8
const REQUEST_API_KEY = "CS530"

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
	runtime.GOMAXPROCS(4)

	cachedImages = make(map[string]int64)
	cachedImagesMutex = new(sync.Mutex)
	imageIsBuilding = make(map[string]bool)
	imageIsBuildingMutex = new(sync.Mutex)

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	var err error
	cli, err = client.NewClientWithOpts(client.WithVersion("1.40"))
	if err != nil {
		panic(err)
	}

	port := 8222
	if len(os.Args) >= 2 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	logger.Printf("server listening at :%d\n", port)
	http.Handle("/", new(FrontHandler))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		panic(err)
	}
}

type FrontHandler struct {
	http.Handler
}

type FailResponse struct {
	Error   bool
	Message string
}

type SuccessResponse struct {
	Result                lambdaResponseResult
	IsWarm                bool
	ExecutionTime         int64
	InternalExecutionTime int64
}

func writeFailResponse(w *http.ResponseWriter, message string) {
	resp := FailResponse{true, message}
	bytes, _ := json.Marshal(resp)
	(*w).Write(bytes)
}

func (h *FrontHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	nameParam := q["name"]
	apiKeyParam := q["key"]

	w.Header().Add("Content-Type", "application/json")

	if len(apiKeyParam) == 0 || apiKeyParam[0] != REQUEST_API_KEY {
		writeFailResponse(&w, "param 'key' is not given or not valid")
		return
	}

	if len(nameParam) == 0 {
		writeFailResponse(&w, "param 'name' is not given")
		return
	}
	startTime := makeTimestamp()

	functionName := nameParam[0]
	imageName := imageTagName(functionName)

	logger.Println("function requested:", functionName)

	cacheExists := false
	cachedImagesMutex.Lock()
	if _, exists := cachedImages[imageName]; exists == true {
		cacheExists = true
	}
	cachedImagesMutex.Unlock()

BUILD_IMAGE:
	// Build image if cache not exist
	if cacheExists == false {
		imageIsBuildingMutex.Lock()
		if imageIsBuilding[imageName] == true {
			// Wait until image is built
			logger.Printf("Image '%s' not found. Wait for build compeletion...\n", imageName)

			for {
				imageIsBuildingMutex.Unlock()
				time.Sleep(time.Second / 20)
				imageIsBuildingMutex.Lock()

				if imageIsBuilding[imageName] == false {
					break
				}
			}
			imageIsBuildingMutex.Unlock()

		} else {
			// Build image
			imageIsBuilding[imageName] = true
			imageIsBuildingMutex.Unlock()

			logger.Printf("Image '%s' not found. Image build start.\n", imageName)
			if err := buildImage(functionName); err != nil {
				writeFailResponse(&w, err.Error())
				return
			}

			imageIsBuildingMutex.Lock()
			imageIsBuilding[imageName] = false
			imageIsBuildingMutex.Unlock()

			logger.Printf("Image '%s' build fin.\n", imageName)
		}
	}

	out, err := runContainer(imageName)
	if err != nil {
		// writeFailResponse(&w, err.Error())
		logger.Println(err.Error(), "Retry...")
		cacheExists = false
		goto BUILD_IMAGE
	}

	endTime := makeTimestamp()
	logger.Println("fin", functionName)

	resp := SuccessResponse{
		out.Result,
		cacheExists,
		endTime - startTime,
		out.EndTime - out.StartTime,
	}
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)

	go updateCachedImages(imageName)
}

func updateCachedImages(imageName string) {
	cachedImagesMutex.Lock()
	cachedImages[imageName] = time.Now().Unix()
	removeOldImages()
	cachedImagesMutex.Unlock()
}
