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

const SERVER_PATH = "http://localhost:8080"
const ACCOUNTS_API_PATH = "/v1/organisation/accounts"

func Fetch(id string) (*AccountData, error){
	url_path, _ := url.Parse(fmt.Sprintf("%s%s/%s", SERVER_PATH, ACCOUNTS_API_PATH, id))

	res, err := http.Get(url_path.String())
	if err != nil {
		log.Fatal(err)
		return nil, &InternalServerError{ErrorMessage: "An error occurred when trying to fetch an Account"}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return &JsonParseError{
			ErrorMessage: "An error occurred when trying to parse request to json",
			StackTrace: err.Error(),
		}
	}

	res, err := http.Post(url_path.String(), "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
		return &InternalServerError{ErrorMessage: "An error occurred when trying to create an new Account"}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
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

func Delete(id string, version int) error {
	url_path, _ := url.Parse(fmt.Sprintf("%s%s/%s", SERVER_PATH, ACCOUNTS_API_PATH, id))
	params := url.Values{}
	params.Add("version", strconv.Itoa(version))

	url_path.RawQuery = params.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url_path.String(), nil)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return &InternalServerError{ErrorMessage: "An error occurred when trying to delete an Account"}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
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