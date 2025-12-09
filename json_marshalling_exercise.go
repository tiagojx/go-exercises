package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
}

func main() {
	u := User{ID: 122, FirstName: "Cleiton", LastName: "Rasta"}

	fmt.Printf("Convertendo struct em JSON...\n")
	jsonConverted, err := json.Marshal(u)
	if err != nil {
		fmt.Println("error: ", err.Error())
	} else {
		// polishing the string to make it more readable.
		jsonSrting := string(jsonConverted)
		rmOpenBracket := strings.ReplaceAll(jsonSrting, "{", "")
		rmClosedBracket := strings.ReplaceAll(rmOpenBracket, "}", "")
		splitedData := strings.Split(rmClosedBracket, ",")
		// printing each item of json (separated by ',')
		for i := range splitedData {
			fmt.Printf("%s\n", splitedData[i])
		}
	}

	fmt.Printf("\nConvertendo JSON em struct...\n")
	u2 := User{}
	err = json.Unmarshal(jsonConverted, &u2)
	if err != nil {
		fmt.Println("error: ", err.Error())
	} else {
		fmt.Printf("ID:\t\t%d\nFirst Name:\t%s\nLast Name:\t%s\n", u2.ID, u2.FirstName, u2.LastName)
	}
}
