package models

import (
	"strings"
	"fmt"
)

type Showtime struct {
	IdTheater string `json:"idTheater"`
	IdMovie   string `json:"idMovie"`
	Language  string `json:"language"`
	RoomType  string `json:"roomType"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	ActiveRecord
}

func NewShowtime(IdTheater string,IdMovie   string,Language  string,RoomType  string,Date      string,Time      string) (Showtime){
	return Showtime{IdTheater,IdMovie,Language,RoomType,Date,Time, "",ActiveRecord{nil}}
}

func (this *Showtime) Save() {
	this.Name = strings.Join([]string{this.IdTheater,this.IdMovie,this.Language,this.RoomType,this.Date,this.Time},"-")
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "name",StructToMap(this))
}
