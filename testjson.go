package main

/*import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	type ColorGroup struct {
		ID     int
		Name   string `json:"name"`
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "red",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := json.MarshalIndent(group, " ", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	var jsonBlob = []byte(`[
		{"name": "latypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string `json:"name"`
		Order string
	}
	var animals []Animal
	uerr := json.Unmarshal(jsonBlob, &animals)
	if uerr != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)

}*/

import (
	"encoding/json"
	"fmt"
	// simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"os"
	// "path/filepath"
	// "regexp"
)

type Config struct {
	ServerIp     string `json:"serverip"`
	ServerPort   string `json:"serverport"`
	Separator    string `json:"separator"`
	KeyIndex     string `json:"keyindex"`
	RedisType    string `json:"redistype"`
	RedisKeyName string `json:"rediskeyname"`
	Format       string `json:"format"`
}

func main() {

	/*	js, _ := simplejson.NewJson([]byte(`{
			"test12": {
				"string_array": ["asdf", "ghjk", "zxcv"],
				"array": [1, "2", 3],
				"arraywithsubs": [{"subkeyone": 1},
				{"subkeytwo": 2, "subkeythree": 3}],
				"int": 10,
				"float": 5.150,
				"bignum": 9223372036854775807,
				"string": "simplejson",
				"bool": true
			}
		}`))
		b, _ := js.EncodePretty()
		os.Stdout.Write(b)*/

	/*	conf := Config{
			ServerIp:     "127.0.0.1",
			ServerPort:   "8080",
			Separator:    "|",
			KeyIndex:     "0,3",
			RedisType:    "hset",
			RedisKeyName: "testh",
			Format:       "first||and second ||!",
		}

		ParaToFile("paratest", conf)*/
	var cfg Config
	FileToPara("paratest", &cfg)

	fmt.Printf("%+v", cfg)

}

func ParaToFile(filename string, cfg Config) {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		fmt.Println(err)
	}
	by, marshalerr := json.MarshalIndent(cfg, "", "    ")
	if marshalerr != nil {
		fmt.Println(marshalerr)
	}
	f.Write(by)
	f.Close()
}

func FileToPara(filename string, cfg *Config) {
	var b []byte
	f, err := os.Open(filename)

	defer f.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		b, _ = ioutil.ReadAll(f)
		json.Unmarshal(b, cfg)
	}
}
