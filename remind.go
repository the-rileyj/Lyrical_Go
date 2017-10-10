package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type datAuth struct {
	Sid    string `json:"sid"`
	Token  string `json:"token"`
	Number string `json:"number"`
}

ype 

func main() {
	/*var dat datAuth
	bdata, err := ioutil.ReadFile("keys.json")
	if err != nil {
		log.Fatal("Could not read data properly")
	}
	err = json.Unmarshal(bdata, &dat)
	if err != nil {
		log.Fatal("Error Unmarshalling the data")
	}
	client := twilio.NewClient(dat.Sid, dat.Token, nil)
	//_, err = client.Messages.SendMessage(dat.Number, "+17013186330", "Sent via go :) ✓", nil)
	uri, err := url.Parse("http://therileyjohnson.com/public/files/mp3/MaskOff.mp3")
	_, err = client.Calls.MakeCall(dat.Number, "+17013186330", uri)
	if err != nil {
		log.Fatal(err)
	}*/
	res, err := http.Get("http://api.musixmatch.com/ws/1.1/track.search?apikey=f3a615206f8902d77d6ef5aba179a496&q_artist=Future&q_track=Mask%20Off")
	if err != nil {
		log.Fatal("Error opening response")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Printf("%s", string(body))
}
