package models

import (
	"strings"
	"fmt"
	"strconv"
	"github.com/jmcvetta/neoism"
)

type Theater struct {
	IdTheater int `json:"idTheater"`
	Name string `json:"name"`	
	Key string `json:"key"`	
	ActiveRecord
}

func NewTheater(IdTheater int,Name string) (Theater){
	return Theater{IdTheater,Name, "", ActiveRecord{nil}}
}

func TheatersByCity(city City) (results []Theater){
        cq := neoism.CypherQuery{
                Statement: `MATCH (c:City)--(t:Theater) WHERE c.name="`+city.Name+`" RETURN t.idTheater as idTheater, t.name as name ORDER BY t.name`,
                Parameters: map[string]interface{}{"name":city.Name},
                Result: &results,
        }

        GlobalDb().Cypher(&cq)
	return
}

func (this *Theater) GenKey() string{
	return strconv.Itoa(this.IdTheater)+"-"+this.Name
}

func (this *Theater) Save() {
	this.Key = this.GenKey()
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "key", StructToMap(this))
}


