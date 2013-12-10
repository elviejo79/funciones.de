package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"../cinemex"
	_ "time"
)
type info struct{
	address string
	phone string
}

type cinema struct {
	Area_id  int
	Id       int
	Info     info
	Lat      float32
	Lng      float32
	Name     string
	Platinum bool
	State_id int
}

type area struct {
	Cinemas  []cinema
	Id       int
	Name     string
	State_id int
}

func TestDecoding() {
	jsonArea, _ := ioutil.ReadFile("./areas_pretty.json")
	var m []area
	json.Unmarshal(jsonArea, &m)

	fmt.Printf("json %#v \n", m)
}

const SHOWTIMES_URL := "http://cinemex.com/partials/sidebarCinema/%d/date-%s"
func TheaterShowtimes(t models.Theater) ([]models.Showtime){
}

func main(){
	theaters := cinemex.ExtractTheaters("http://cinemex.com/")
	//t := time.Now().Format("20060102")
	fmt.Printf("theater: %#v \n", theaters)
	showtimes,_ :=cinemex.ExtractMovies(

	fmt.Printf("showtimes %#v\n\n",showtimes)
}
