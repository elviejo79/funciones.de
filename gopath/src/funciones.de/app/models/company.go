package models

import (
	"strings"
	"fmt"
)

type Company struct {
	Url string `json:"url"`
	Name string `json:"name"`
	ActiveRecord
}

func NewCompany(Url,Name string) (Company){
	return Company{Url,Name, ActiveRecord{nil}}
}

func (this *Company) Save() {
	aType := strings.Split(fmt.Sprintf("%T",this),".")
	this.ActiveRecord.Save(aType[len(aType)-1], "name", StructToMap(this))
}
