package models

import (
	"strings"
	"fmt"
	"strconv"
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

func (this *Theater) Save() {
	this.Key = strconv.Itoa(this.IdTheater)+"-"+this.Name
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "key", StructToMap(this))
}


