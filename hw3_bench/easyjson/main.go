package main

import "fmt"

//easyjson:json
type Browsers struct {
	//Browsers []string
	Email string
	Name  string
}

func main() {
	d := &Browsers{}
	d.UnmarshalJSON([]byte(`{"company":"Jatri","country":"Kenya","email":"eum_rerum_explicabo@Topiczoom.info","job":"Web Developer #{N}","name":"Susan Ellis","phone":"187-70-57"} `))
	fmt.Println(d.Email)
	fmt.Println(d.Name)

}
