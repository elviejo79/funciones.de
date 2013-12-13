package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"../cinepolis"
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


func main(){
//	theaters := cinemex.ExtractTheaters("http://cinemex.com/")
	movies :=cinepolis.ExtractMovies()
	fmt.Printf("cuantas peliculas son? %d \n", len(movies))

}
