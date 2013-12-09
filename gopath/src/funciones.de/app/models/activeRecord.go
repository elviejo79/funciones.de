package models

import (
	"encoding/json"
	"github.com/jmcvetta/neoism"
	"log"
)


var db *neoism.Database
func GlobalDb() *neoism.Database {
	if db == nil {
		db, _ = neoism.Connect("http://localhost:7474/db/data")
	}
	return db
}


type ActiveRecord struct {
	node *neoism.Node
}

func (this *ActiveRecord) Save (label string, key string, obj map[string]interface{}){
	var err error
	this.node,_,err = GlobalDb().GetOrCreateNode(label,key,obj)
	if err != nil {
                log.Fatal("error en GetOrCreateNode %#v",err)
        }

	err = this.node.AddLabel(label)
	if err != nil {
                log.Fatal("error en AddLabel %#v",err)
        }


}

func (this *ActiveRecord) Node() (*neoism.Node){
	return this.node
}

func StructToMap(myStruct interface{}) (mapMyStruct map[string]interface{}) {
	jsonMyStruct, _ := json.Marshal(myStruct)
	json.Unmarshal(jsonMyStruct, &mapMyStruct)
	return
}

func MapToStruct(myMap map[string]interface{}, v interface{}) (err error) {
	jsonMyMap, _ := json.Marshal(myMap)
	return json.Unmarshal(jsonMyMap, v)
}
