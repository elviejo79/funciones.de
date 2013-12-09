package models

import (
	"strings"
	"fmt"
)

type City struct {
	IdCity int `json:"idCity"`
	Name string `json:"name"`
	ActiveRecord
}

func NewCity(IdCity int,Name string) (City){
	return City{IdCity,Name, ActiveRecord{nil}}
}

func (this *City) Save() {
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1],"name", StructToMap(this))
}

