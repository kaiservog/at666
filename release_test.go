package main

import (
  "testing"
  "bytes"
  "net/http"
  "net/url"
  "io/ioutil"
)

//const ATAPP_SERVER_URL = "https://atapp.herokuapp.com"
const ATAPP_SERVER_URL = "http://localhost:9090"



func doLastId() (*http.Response, error) {
  url := ATAPP_SERVER_URL + "/at/comment/lastId/0.000000/0.000000"
  resp, err := http.Get(url);

  return resp, err
}

func TestLastIdMustReturnSomething(t *testing.T) {
  resp, err := doLastId()

  if resp.Status != "200 OK" {
    t.Error("Return is not 200 OK it is", resp.Status)
  }

  if err != nil {
    t.Error(err)
  }
}

func TestLastIdMustReturnError(t *testing.T) {
  url := ATAPP_SERVER_URL + "/at/comment/lastId/0.000000/huehue"
  resp, err := http.Get(url);
  if resp.Status != "500 Internal Server Error" {
    t.Error("Return is not 500 Internal Server Error it is", resp.Status)
  }

  if err != nil {
    t.Error(err)
  }
}

func TestMustUpdateLastId(t *testing.T) {
    myUrl := ATAPP_SERVER_URL + "/at/comment/lastId/0.000000/0.000000"
    resp, err := http.Get(myUrl);

    body, err := ioutil.ReadAll(resp.Body)
    lastId := string(body)

    form := url.Values{}
    form.Add("lat", "0.000000")
    form.Add("lon", "0.000000")
    form.Add("nick", "GO LANG")
    form.Add("text", "Hi from test!")

    req, err := http.NewRequest("PUT",
        ATAPP_SERVER_URL + "/at/comment",
        bytes.NewBufferString(form.Encode()))

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work

    client := http.Client{}
    resp, err = client.Do(req)

    if err != nil {
      t.Error(err)
    }

    if resp.Status != "200 OK" {
      t.Error("status is not 200 OK it is", resp.Status)
    }

    resp, err = doLastId()

    if err != nil {
      t.Error(err)
    }

    if resp.Status != "200 OK" {
      t.Error("status is not 200 OK it is", resp.Status)
    }

    body, err = ioutil.ReadAll(resp.Body)
    lastIdUpdate := string(body)

    if lastIdUpdate == lastId {
      t.Error("LastId Not updated", lastIdUpdate, lastId)
    }
}

func doPutPeople() (*http.Response, error) {
    atserverUrl := ATAPP_SERVER_URL + "/at/people"
    form := url.Values{}
    form.Add("lat", "0.000000")
    form.Add("lon", "0.000000")
    form.Add("nick", "GO LANG")

    req, err := http.NewRequest("PUT", atserverUrl, bytes.NewBufferString(form.Encode()))
    if err != nil {
      return nil, err
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
    client := http.Client{}
    return client.Do(req)
}

func TestPutPeople(t *testing.T) {
  resp, err := doPutPeople()

  if err != nil {
    t.Error(err)
  }

  if resp.Status != "200 OK" {
    t.Error("status is not 200 OK it is", resp.Status)
  }
}

func TestGetPeople(t *testing.T) {
  doPutPeople()

  url := ATAPP_SERVER_URL + "/at/people/0.000000/0.000000"
  resp, err := http.Get(url);
  if resp.Status != "200 OK" {
    t.Error("Return is not 200 OK it is", resp.Status)
  }

  if err != nil {
    t.Error(err)
  }

}
