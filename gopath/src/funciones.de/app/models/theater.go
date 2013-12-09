package models

import (
	"strings"
	"fmt"
)

type Theater struct {
	IdTheater int `json:"idTheater"`
	Name string `json:"name"`	
	ActiveRecord
}

func NewTheater(IdTheater int,Name string) (Theater){
	return Theater{IdTheater,Name, ActiveRecord{nil}}
}

func (this *Theater) Save() {
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "idTheater", StructToMap(this))
}


