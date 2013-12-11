package models

import (
	"strings"
	"fmt"
	"regexp"
	"github.com/jmcvetta/neoism"
	"log"
)

type City struct {
	IdCity int `json:"idCity"`
	Name string `json:"name"`
	ActiveRecord
}

func NewCity(IdCity int,Name string) (City){
	re := regexp.MustCompile("(\\s+)")
	Name=strings.TrimSpace(re.ReplaceAllString(Name, " "))
	return City{IdCity,Name, ActiveRecord{nil}}
}

func AllCities()(results []City){
        cq := neoism.CypherQuery{
                Statement: `MATCH (c:City) RETURN c.idCity as idCity, c.name as name ORDER BY c.name`,
                Parameters: map[string]interface{}{},
                Result: &results,
        }
        GlobalDb().Cypher(&cq)
	return
}

func CityByName(Name string) (city City){
	log.Printf("CityByName %#v", Name)
	var cities []City
        cq := neoism.CypherQuery{
                Statement: `MATCH (c:City) WHERE c.name="`+Name+`" RETURN c.idCity as idCity, c.name as name ORDER BY c.name`,
                Parameters: map[string]interface{}{"name":Name},
                Result: &cities,
        }
        GlobalDb().Cypher(&cq)

	if len(cities)>0 {
		city = cities[0]
	}else{
		city = NewCity(0,"")
	}
	return
}

func (this *City) Save() {
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1],"name", StructToMap(this))
}

func (this *City) Abbrev() string{
	return this.Name
}
