# gopexels

Gopexels is a simple [Pexels](https://pexels.com) image downloader based on it's internal API.
It has a nice progress bar, and is multi-threaded !

## Installation

```bash
go get -u github.com/jesuiscamille/gopexels
cd $GOPATH/src/github.com/jesuiscamille/gopexels
make
```

## Usage

```bash
Usage of ./gopexels:
  -amount int
    	The amount of images to download (default 100)
  -pageAmount int
    	The amount of pages to fetch (default 10)
  -query string
    	The pexels search term to be used
  -threads int
    	The amount of threads to use to download the images (default 3)
```

### Note

Note that, if you want a really high number of images, you must increase the default number of pages fetched, by adding `-pageAmount 50` ( per example ) in your command line.

## Example

### Invoking 
( while on $GOPATH/src/github.com/jesuiscamille/gopexels, after having executed `make` )

```bash
./gopexels -query cats -amount 1000 -pageAmount 20 -threads 5
```

### Output

```
reinit:gopexels camille$ ./gopexels -query cats -amount 1000 -pageAmount 20 -threads 1
2018/07/10 13:02:09 Getting results pages for the query...
2018/07/10 13:02:11 Search results pages downloaded: 20.
2018/07/10 13:02:11 Starting the downloads... Threads: 1
 23 / 1000 [===>------------------------------------------------------------------------------------------------------------------------------------------]   2.30% 5h55m48s^C
reinit:gopexels camille$ ls images/
pexels-photo-1034832.jpeg pexels-photo-127028.jpeg  pexels-photo-193255.jpeg  pexels-photo-326875.jpeg  pexels-photo-617278.jpeg  pexels-photo-774731.jpeg
pexels-photo-1077449.jpeg pexels-photo-140134.jpeg  pexels-photo-209037.jpeg  pexels-photo-399647.jpeg  pexels-photo-730537.jpeg  pexels-photo-96428.jpeg
pexels-photo-1096091.jpeg pexels-photo-170304.jpeg  pexels-photo-236230.jpeg  pexels-photo-416160.jpeg  pexels-photo-730896.jpeg  pexels-photo-96938.jpeg
pexels-photo-126407.jpeg  pexels-photo-177809.jpeg  pexels-photo-248280.jpeg  pexels-photo-533083.jpeg  pexels-photo-749212.jpeg  pexels-photo-979503.jpeg
```

Here, the program sucessfully fetched 1000 image urls out of twenty pages downloaded.
23 images got entirely downloaded by _one_ thread. ( My connection is too poor to launch multiple threads )
You can see at the end the list of the images that got downloaded.
( there may be more than 23 because the program may got closed while writting a file to the disk, so before incrementing the progress bar )

## Contributing

You are encouraged to contribute to this project, by making an API out of it, per example ( I wanted to create a go package, but I'm a bit lazy right now )
Any pull request or issue will be warmly welcomed.

## Credits

Thanks to [@cheggaaa](https://github.com/cheggaaa) for the [pb project](https://github.com/cheggaaa/pb).

## LICENSE

`¯\_(ツ)_/¯,`, just don't do bad things and respect the image licenses :)
Oh, and respect the [pb license](https://github.com/cheggaaa/pb/blob/master/LICENSE) !
If there are any licensing issue, contact me, and I'll fix that.
