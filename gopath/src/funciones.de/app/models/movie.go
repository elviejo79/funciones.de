package models

import (
	"strings"
	"fmt"
	"github.com/jmcvetta/neoism"
	"log"
)

type Movie struct {
	Title string `json:"title"`
	ActiveRecord
}

func NewMovie(Title string) (Movie){
	return Movie{Title, ActiveRecord{nil}}
}

func MoviesByTheaters(ts []Theater) (results []Movie){
	var ts_keys []string
	for _,t := range ts {
		ts_keys = append(ts_keys,t.GenKey())
	}
	in := strings.Join(ts_keys,"', '")

	cq := neoism.CypherQuery{
                Statement: `MATCH (t:Theater)--(s:Showtime)--(m:Movie) WHERE t.key IN ['`+in+`'] RETURN m.title as title, count(s) ORDER BY count(s) DESC`,
                Parameters: map[string]interface{}{},
                Result: &results,
        }
	
	log.Printf("Movies By Theater %#v \n CQ %#v",results,cq.Statement)
        GlobalDb().Cypher(&cq)
	return
}

func (this *Movie) Save() {
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "title", StructToMap(this))
}

func (this *Movie) GetShowtimesInTheaters(ts []Theater) (results []Showtime){
	var ts_keys []string
	for _,t := range ts {
		ts_keys = append(ts_keys,t.GenKey())
	}
	in := strings.Join(ts_keys,"', '")

	var cypher_return string
	for k,_ := range StructToMap(new(Showtime)) {
		cypher_return = cypher_return +fmt.Sprintf("s.%s as %s, ",k,k)
	}
	cypher_return=cypher_return[:len(cypher_return)-2]

	cq := neoism.CypherQuery{
		Statement: `start m=node:node_auto_index(title='`+this.Title+
			`') MATCH (t:Theater)--(s:Showtime)--(m) WHERE t.key IN ['`+in+`'] RETURN DISTINCT `+cypher_return+` ORDER BY s.time ASC`,
		Parameters: map[string]interface{}{},
		Result: &results,
	}
	
	log.Printf("Showtimes For Movie %#v \n CQ %#v",results,cq.Statement)
	err := GlobalDb().Cypher(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return
}
