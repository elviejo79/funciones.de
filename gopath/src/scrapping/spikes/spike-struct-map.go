package main

import (
	"funciones.de/app/models"
	_ "github.com/jmcvetta/neoism"
	"fmt"
	_ "log"
)

func main(){
	showtime := models.NewShowtime("idT","idM","L","rT","d","t")
	company := models.NewCompany("www.something.com","something")

	fmt.Printf("map from funciton %#v \n", models.StructToMap(showtime))

	/*db, err := neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		log.Fatal("failed to open default database", err)
	}*/

	fmt.Printf("el nodo antes   de save: %#v \n", company.Node())
	company.Save()
	fmt.Printf("el nodo después de save: %#v \n", company.Node())
//	fmt.Printf("el todo el objeto      : %#v \n", company)
//	newCo, created,err := db.GetOrCreateNode("Company","name",structToMap(company))
//	newCo.AddLabel("Company")
//	fmt.Printf("CreateUniqueCompany: %#v\n created: %#v \n err: %#v",newCo,created,err)
/*
	newCo, err := db.CreateNode(structToMap(company))
	newCo.AddLabel("Company")
	if err != nil {
		log.Fatal("failed to create node:", err)
	}
	fmt.Printf("map from funciton %#v \n", newCo)
	// leer datos de la base de datos
 /*
	res := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}

        //res := []models.Company{}
        cq := neoism.CypherQuery{
                Statement: "MATCH (n:Company) where n.name = 'something'  RETURN n",
                Parameters: map[string]interface{}{
                        "name": "something",
                },
                Result: &res,
        }

        err = db.Cypher(&cq)
        if err != nil {
		fmt.Printf("error: %#v \n", err)
        }
	fmt.Printf("\n ---- from neo4j %#v \n", res[0].N.Data)

	var readCo models.Company
	mapToStruct(res[0].N.Data, &readCo)

	fmt.Printf("\n --- Company? %#v \n\n\n", readCo)

	var resType []models.Company
        cq1 := neoism.CypherQuery{
                Statement: "MATCH (n:Company) where n.name = 'something'  RETURN n.url as url, n.name as name",
                Parameters: map[string]interface{}{
                        "name": "something",
                },
                Result: &resType,
        }

        err = db.Cypher(&cq1)
        if err != nil {
		fmt.Printf("error: %#v \n", err)
        }
	fmt.Printf("\n ---- from neo4j %#v \n", resType)

*/
}
