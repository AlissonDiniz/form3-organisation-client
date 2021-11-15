package accounts

import (
	"fmt"
	"log"
	"bytes"
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

var SERVER_PATH = ""
var ACCOUNTS_API_PATH = "/v1/organisation/accounts"

func Setup(server_path string) {
	SERVER_PATH = server_path
}

func Fetch(id string) (*AccountData, error){
	url_path, _ := url.Parse(fmt.Sprintf("%s%s/%s", SERVER_PATH, ACCOUNTS_API_PATH, id))

	res, err := http.Get(url_path.String())
	if err != nil {
		log.Println(err.Error())
		return nil, &InternalServerError{ErrorMessage: "An error occurred when trying to fetch an Account"}
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return nil, &ReadResponseBodyError{
			ErrorMessage: "An error occurred when trying to read the response body",
			StackTrace: err.Error(),
		}
	}

	if res.StatusCode == 404 {
		var err NotFoundError
		json.Unmarshal(body, &err)
		return nil, &err
	}

	var data ResponseAccount
	json.Unmarshal(body, &data)
	return data.Data, nil
}

func Create(account_data *AccountData) error {
	url_path, _ := url.Parse(fmt.Sprintf("%s%s", SERVER_PATH, ACCOUNTS_API_PATH))

	request := RequestCreateAccount {
		Data: account_data,
	}
	json_data, err := json.Marshal(request)
	if err != nil {
		log.Println(err.Error())
		return &JsonParseError{
			ErrorMessage: "An error occurred when trying to parse request to json",
			StackTrace: err.Error(),
		}
	}

	res, err := http.Post(url_path.String(), "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Println(err.Error())
		return &InternalServerError{ErrorMessage: "An error occurred when trying to create an new Account"}
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return &ReadResponseBodyError{
			ErrorMessage: "An error occurred when trying to read the response body",
			StackTrace: err.Error(),
		}
	}

	if res.StatusCode != 201 {
		if res.StatusCode == 409 {
			var err ConflictError
			json.Unmarshal(body, &err)
			return &err
		} else {
			var err InternalServerError
			json.Unmarshal(body, &err)
			return &err
		}
	}
	return nil
}

func Delete(id string, version int64) error {
	url_path, _ := url.Parse(fmt.Sprintf("%s%s/%s", SERVER_PATH, ACCOUNTS_API_PATH, id))
	params := url.Values{}
	params.Add("version", strconv.FormatInt(version, 10))

	url_path.RawQuery = params.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url_path.String(), nil)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return &InternalServerError{ErrorMessage: "An error occurred when trying to delete an Account"}
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return &ReadResponseBodyError{
			ErrorMessage: "An error occurred when trying to read the response body",
			StackTrace: err.Error(),
		}
	}

	if res.StatusCode != 204 {
		if res.StatusCode == 404 {
			var err NotFoundError
			json.Unmarshal(body, &err)
			return &err
		} else {
			var err InternalServerError
			json.Unmarshal(body, &err)
			return &err
		}
	}
	return nil
}