package models

import (
	"strings"
	"fmt"
)

type Movie struct {
	Title string `json:"title"`
	ActiveRecord
}

func NewMovie(Title string) (Movie){
	return Movie{Title, ActiveRecord{nil}}
}

func (this *Movie) Save() {
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "title", StructToMap(this))
}

