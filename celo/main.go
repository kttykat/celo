package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	funcs "celo/main/celo/utils"
)



func main() {
	cli := os.Args[1:]
	if len(cli) == 0 {
		fmt.Println("celo <github_user>")
	} else {
		name := cli[0]
		url := fmt.Sprintf(`http://localhost:5124/%v`, name)

		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data := []funcs.Repo{}
		json.Unmarshal(body, &data)
		funcs.RepoList(data, name)
	}
}
