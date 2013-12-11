package models

import (
	"strings"
	"fmt"
	"github.com/jmcvetta/neoism"
	"log"
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

func ShowtimesForTheaters(ts []Theater) (results []Showtime){
	var ts_keys []string
	for _,t := range ts {
		ts_keys = append(ts_keys,t.GenKey())
	}
	in := strings.Join(ts_keys,"', '")

	var cypher_return string
	for k,_ := range StructToMap(NewShowtime("","","","","","")) {
		cypher_return = cypher_return +fmt.Sprintf("s.%s as %s, ",k,k)
	}
	cypher_return=cypher_return[:len(cypher_return)-2]

	cq := neoism.CypherQuery{
		Statement: `MATCH (t:Theater)--(s:Showtime)--(m) WHERE t.key IN ['`+in+`'] RETURN `+cypher_return+` ORDER BY s.time ASC`,
		Parameters: map[string]interface{}{},
		Result: &results,
	}
	
	log.Printf("Theaters %v Showtimes For Theaters %#v \n CQ %#v",ts,results,cq.Statement)
	GlobalDb().Cypher(&cq)
	return
}


func (this *Showtime) Save() {
	this.Name = strings.Join([]string{this.IdTheater,this.IdMovie,this.Language,this.RoomType,this.Date,this.Time},"-")
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "name",StructToMap(this))
}
