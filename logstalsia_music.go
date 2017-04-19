package main

import (
"flag"
"net/http"
"runtime"
"strconv"
"time"
"os"
"log"
"bufio"
"fmt"
"strings"
)

var url *string
var lyricsFile *string

func main() {

	url = flag.String("url", "https://bots.discord.pw/", "The URL to hit.")
	lyricsFile = flag.String("lyrics", "lyrics.txt", "The lyrics.")
	flag.Parse()

	// Use all procs
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, err := os.Open(*lyricsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	songInfo := make(map[int]string)
	last := -1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lyricParts := strings.Split(scanner.Text(), "|")
		if lyricParts != nil && len(lyricParts) == 2 {
			partTime, err := strconv.Atoi(strings.TrimSpace(lyricParts[1]))
			if err != nil {
				log.Fatal(err)
				continue
			}

			songInfo[partTime] =  strings.Replace( strings.TrimSpace(strings.TrimSpace(lyricParts[0])), " " , "_", -1)
			last = partTime

		}
	}

	fmt.Println("Done input.. Map length: ", len(songInfo))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration("10ms")
	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 0; i < 100; i++ {
		time.Sleep(duration)
		go doLyric("CORI_SYNC")
	}

	var lastLyric string = "+"
	currentTick := 0
	for currentTick <= last {

		time.Sleep(duration)

		lyric, ok := songInfo[currentTick]
		currentTick++
		if !ok {
			if !strings.EqualFold(lastLyric, "+") {
				go doLyric(lastLyric)
			}
			continue
		}

		lastLyric = lyric
		if !strings.EqualFold(lastLyric, "+") {
			fmt.Println( currentTick, " | ", lastLyric )
			go doLyric( lastLyric )
		} else {
			fmt.Println( currentTick, " | ~~PAUSE~~" )
		}

	}

}

func doLyric( lyric string ) {
	http.Get(*url + lyric)
}
