package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	twilio "github.com/saintpete/twilio-go"
)

//Struct to hold the various key data
type datAuth struct {
	Sid    string `json:"sid"`
	Token  string `json:"token"`
	Lkey   string `json:"lyric_key"`
	Number string `json:"number"`
}

type id struct {
	ID int32 `json:"track_id"`
}

type track struct {
	Track id `json:"track"`
}

type trackList struct {
	Tracks []track `json:"track_list"`
}

type fullSearchBody struct {
	TrackList trackList `json:"body"`
}

type fullSearchMessage struct {
	Body fullSearchBody `json:"message"`
}

type lyrics struct {
	Lyrics string `json:"lyrics_body"`
}

type lyricBody struct {
	LyricObject lyrics `json:"lyrics"`
}

type fullLyricBody struct {
	LyricBody lyricBody `json:"body"`
}

type fullLyricMessage struct {
	Body fullLyricBody `json:"message"`
}

//Gets
func getSongID(artist, track, api string) string {
	base := "http://api.musixmatch.com/ws/1.1/track.search?apikey=" + api
	var b bytes.Buffer
	b.WriteString(base)
	if track != "" {
		b.WriteString("&q_track=")
		b.WriteString(url.QueryEscape(track))
	}
	if artist != "" {
		b.WriteString("&q_artist=")
		b.WriteString(url.QueryEscape(artist))
	}
	if b.String() == base {
		log.Fatal("Both 'artist' and 'track' parameters cannot be empty")
	}
	res, err := http.Get(b.String())
	if err != nil {
		log.Fatal("Error fetching response")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response")
	}
	var m fullSearchMessage
	json.Unmarshal(body, &m)
	return fmt.Sprintf("%d", m.Body.TrackList.Tracks[0].Track.ID)
}

func getSongLyrics(ID, api string) string {
	res, err := http.Get(fmt.Sprintf("http://api.musixmatch.com/ws/1.1/track.lyrics.get?apikey=%s&track_id=%s", api, ID))
	if err != nil {
		log.Fatal("Error fetching response")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response")
	}
	var m fullLyricMessage
	json.Unmarshal(body, &m)
	return m.Body.LyricBody.LyricObject.Lyrics
}

func main() {
	var artist, title, keys, span string //Variables for getting the command line args corresponding to the variables' names
	flag.StringVar(&artist, "artist", "", "The name of the artist of the song you want to lookup")
	flag.StringVar(&title, "title", "", "The name of the song you want to lookup")
	flag.StringVar(&keys, "keys", "", "The location of the keys for the Twilio and MusixMatch API's")
	flag.StringVar(&span, "span", "", "The time span over which the lyrics are to be sent every 'span' / 'number of verses' amount of time")
	if keys == "" {
		log.Fatal("Need the location for the API keys, please")
	}
	var dat datAuth
	bdata, err := ioutil.ReadFile(keys)
	if err != nil {
		log.Fatal("Could not read data properly")
	}
	err = json.Unmarshal(bdata, &dat)
	if err != nil {
		log.Fatal("Error Unmarshalling the data")
	}
	lyrics := getSongLyrics(getSongID(artist, title, dat.Lkey), dat.Lkey)
	client := twilio.NewClient(dat.Sid, dat.Token, nil)

	//Assuring the starting date argument isn't empty and if it is defaulting to right now
	if span != "" {
		tTime, err := time.ParseDuration(span) //Parse the start time into the local time
		if err != nil {
			log.Fatal("Error parsing starting time")
		}
		slyrics := strings.Split(lyrics, "\n\n")
		for _, l := range slyrics[:len(slyrics)-1] {
			time.Sleep(time.Duration(tTime.Nanoseconds() / int64(len(slyrics)-1)))
			_, err = client.Messages.SendMessage(dat.Number, "+17013186330", l, nil)
		}
	} else {
		_, err = client.Messages.SendMessage(dat.Number, "+17013186330", lyrics, nil)
	}
}
