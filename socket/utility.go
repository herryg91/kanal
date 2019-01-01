package socket

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func Get(url string, queryString string, additionalHeader map[string]string) (map[string]interface{}, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", url, queryString), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	for k, v := range additionalHeader {
		req.Header.Set(k, v)
	}
	response, err := client.Do(req)

	if err != nil {
		log.Println("[HTTP GET] ", err)
		return make(map[string]interface{}), err
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("[HTTP GET] error reading bytes : ", err)
		return make(map[string]interface{}), err
	}

	var jsonData map[string]interface{}

	if err := jsoniter.ConfigFastest.Unmarshal(responseData, &jsonData); err != nil {
		log.Println("[HTTP GET] error reading json response : ", err)
		return make(map[string]interface{}), err
	}

	return jsonData, nil
}
