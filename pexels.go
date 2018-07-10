package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/cheggaaa/pb"
	"net/http"
)

const (
	path string = "./images/"
)

type Image struct {
	URL         string
	Filename    string
	downloaded  bool
	downloading bool
}

func (i *Image) Downloading() bool {
	return i.downloading
}

func (i *Image) SetDownloading() {
	i.downloading = true
}

func (i *Image) DownloadToDir(dirPath string) (err error) {
	// download the image
	resp, err := http.Get(i.URL)
	if err != nil {
		return
	}

	// open the file
	file, err := os.Create(path + i.Filename)

	if err != nil {
		return
	}

	defer file.Close()

	// create the reader and writters we'll connect
	respReader := bufio.NewReader(resp.Body)
	fileWriter := bufio.NewWriter(file)

	// write the content of the resp to the file
	_, err = respReader.WriteTo(fileWriter)

	if err != nil {
		return
	}

	// flag the image as downloaded
	i.downloaded = true

	return
}

func (i *Image) Downloaded() bool {
	return i.downloaded
}

func NewImage(url string, filename string) *Image {
	return &Image{
		URL:      url,
		Filename: filename,
	}
}

// a page contains a number of images
type Page struct {
	Number int
	Images []Image
}

// AddImage adds an image to a page
func (p *Page) AddImage(i Image) {
	p.Images = append(p.Images, i)
}

func makeResultDir() (err error) {

	if _, rrr := os.Stat(path); os.IsNotExist(rrr) {
		err = os.Mkdir(path, os.ModePerm)
	}

	return
}

// Query query pexels for pages matching a provided query string
//	and return them
// BUG(@jesuiscamille): It may worth it to implement a way to select the number of
//	pages to query and return
func Query(queryString string, amount int) (pages []Page, err error) {
	// check if the search term exist on pexels
	var scanner *bufio.Scanner

	noMatch := regexp.MustCompile("/innerHTML = '';")
	// get the page for the search term
	resp, err := http.Get(fmt.Sprintf("https://www.pexels.com/search/%s/?page=9999&format=js", queryString))

	if err != nil {
		err = errors.New("Error while downloading the pages list: " + err.Error())
		return
	}

	scanner = bufio.NewScanner(resp.Body)
	var tmpPage1 string

	for {
		continu := scanner.Scan()
		tmpPage1 += scanner.Text()
		if !continu {
			break
		}
	}

	// the search term does not exist on pexels
	if noMatch.MatchString(tmpPage1) {
		err = errors.New("The search term did not return anything on pexels.")
		return
	}

	imageRegex := regexp.MustCompile(`photos/[0-9]{1,10}/pexels-photo-[0-9]{1,10}\.jpeg`)
	filenameRegex := regexp.MustCompile(`pexels-photo-[0-9]{0,9}\.jpeg`)

	pageNb := 0

	var tmpPage2 string

	log.Println("Getting results pages for the query...")
	for {
		if pageNb > amount {
			break
		}

		resp, err = http.Get(fmt.Sprintf("https://www.pexels.com/search/%s/?page=%d?format=js", queryString, pageNb))
		scanner = bufio.NewScanner(resp.Body)

		for {
			continu := scanner.Scan()
			tmpPage2 += scanner.Text()
			if !continu {
				break
			}
		}

		//photoURLs := imageRegex.FindAllString(scanner.Text(), -1)
		photoURLs := imageRegex.FindAllString(tmpPage2, -1)

		// if there are no photos left, stop
		if len(photoURLs) <= 0 {
			break
		}

		tmpPage := Page{
			Number: pageNb,
		}

		for _, partImageUrl := range photoURLs {
			tmpPage.AddImage(Image{
				URL:      "https://images.pexels.com/" + partImageUrl,
				Filename: filenameRegex.FindString(partImageUrl),
			})
		}

		pages = append(pages, tmpPage)
		tmpPage2 = ""

		pageNb++
	}

	log.Printf("Search results pages downloaded: %d.\n", pageNb-1)

	return
}

func main() {
	var query string
	var amount int
	var threads int
	var pageAmount int

	flag.StringVar(&query, "query", "", "The pexels search term to be used")
	flag.IntVar(&amount, "amount", 100, "The amount of images to download")
	flag.IntVar(&threads, "threads", 3, "The amount of threads to use to download the images")
	flag.IntVar(&pageAmount, "pageAmount", 10, "The amount of pages to fetch")

	flag.Parse()

	if query == "" {
		log.Fatal("Please select a query using -query")
	}

	if err := makeResultDir(); err != nil {
		log.Fatal("Error while making the result dir: " + err.Error())
	}

	// get the pages for the query

	var pages []Page
	pages, err := Query(query, pageAmount)

	if err != nil {
		log.Fatal("Error while fectching pages: " + err.Error())
	}

	// download all the images of all the pages
	log.Printf("Starting the downloads... Threads: %d\n", threads)

	downloadChan := make(chan Image)
	stopChan := make(chan int)

	// get the total number of images
	total := 0

	for _, page := range pages {
		total += len(page.Images)
	}

	// if the total exceeds our number of images, cut the total to it
	if total > amount {
		total = amount
	}

	// create the progress bar
	var bar *pb.ProgressBar = pb.StartNew(total)
	var wg sync.WaitGroup

	for i := 0; i != threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case img := <-downloadChan:
					if err := img.DownloadToDir(path); err != nil {
						log.Printf("Error while downloading image: %s\n", err.Error())
					}
					bar.Increment()
				case <-stopChan:
					return
				}
			}
		}()
	}

	// send all the images to the download chans
	downloaded := 0

	var sent []string

	for _, page := range pages {
		for _, image := range page.Images {
			if !image.Downloaded() && !image.Downloading() {
				var okcontinue bool = true
				for _, sentUrl := range sent {
					if image.URL == sentUrl {
						okcontinue = false
						break
					}
				}

				if okcontinue {
					sent = append(sent, image.URL)
					downloaded += 1
					image.SetDownloading()
					downloadChan <- image
					if downloaded >= amount {
						break
					}
				}
			}
		}

		if downloaded >= amount {
			break
		}
	}

	// quit all the goroutines
	for i := 0; i != threads; i++ {
		stopChan <- 0
	}

	bar.FinishPrint("All images got downloaded !")
	log.Println("Waiting for the goroutines to exit...")
	wg.Wait()
}
