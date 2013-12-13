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
	BuyLink   string `json:"buyLink"`
	Key       string `json:"key"`
	ActiveRecord
}

func NewShowtime(IdTheater,IdMovie,Language,RoomType,Date,Time,BuyLink string) (Showtime){
	return Showtime{IdTheater,IdMovie,Language,RoomType,Date,Time,BuyLink,"",ActiveRecord{nil}}
}

func ShowtimeByKey(key string) (Showtime){
	var results []Showtime
	var cypher_return string
	for k,_ := range StructToMap(new(Showtime)) {
		cypher_return = cypher_return +fmt.Sprintf("s.%s as %s, ",k,k)
	}
	cypher_return=cypher_return[:len(cypher_return)-2]

	cq := neoism.CypherQuery{
		Statement: `MATCH (s:Showtime) WHERE s.key = '`+key+`'  RETURN DISTINCT `+cypher_return,
		Parameters: map[string]interface{}{},
		Result: &results,
	}
	log.Printf("ShowtimeBy Key %#v",cq)
	GlobalDb().Cypher(&cq)
	return results[0]
}

func ShowtimesForTheaters(ts []Theater) (results []Showtime){
	var ts_keys []string
	for _,t := range ts {
		ts_keys = append(ts_keys,t.GenKey())
	}
	in := strings.Join(ts_keys,"', '")

	var cypher_return string
	for k,_ := range StructToMap(new(Showtime)) {
		if k != "key" {
			cypher_return = cypher_return +fmt.Sprintf("s.%s as %s, ",k,k)
		}
	}
	cypher_return=cypher_return[:len(cypher_return)-2]

	cq := neoism.CypherQuery{
		Statement: `MATCH (t:Theater)--(s:Showtime)--(m) WHERE t.key IN ['`+in+`']  RETURN DISTINCT `+cypher_return+` ORDER BY s.time,s.idMovie,s.idTheater ASC`,
		Parameters: map[string]interface{}{},
		Result: &results,
	}
	
	log.Printf("Theaters %v Showtimes For Theaters %#v \n CQ %#v",ts,results,cq.Statement)
	GlobalDb().Cypher(&cq)
	return
}

func (this *Showtime) GenKey() string {
	return strings.Join([]string{this.IdTheater,this.IdMovie,this.Language,this.RoomType,this.Date,this.Time,},"-")
}

func (this *Showtime) Save() {
	this.Key = this.GenKey()
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "key",StructToMap(this))
}
