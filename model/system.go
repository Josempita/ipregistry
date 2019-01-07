package model

import (
	//_ "github.com/minus5/gofreetds"

	"io/ioutil"

	//_ "github.com/a-palchikov/sqlago"

	"encoding/json"
	"fmt"
	"os"
	//_ "github.com/denisenkom/go-mssqldb"
)

func (m *Model) GetConfig() Config {
	count := len(os.Args[1:])
	var arg string
	if count == 0 {
		arg = ""
	} else {
		arg = os.Args[1]
		if arg != "" {
			arg = arg + "\\"
		}
	}
	raw, err := ioutil.ReadFile(arg + "config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Config
	json.Unmarshal(raw, &c)
	return c
}
