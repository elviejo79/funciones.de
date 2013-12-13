package main

import (
	"fmt"
	"../imdb"
	_ "time"
)



func main(){
	busquedas := []string{"EL HOBBIT 2","EL HOBBIT LA DESOLACIÓN", "THOR", "thor 2", "thor 2", "Thor Un Mundo Oscuro"}
	for _,b := range busquedas {
		results,_ := imdb.Search(b+" T")
		fmt.Printf("busqueda: %s, resultados %#v \n", b, results)

	}

}
