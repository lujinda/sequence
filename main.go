package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
  "net/url"
  "encoding/json"
  "sync"
  "log"
  "io"
)

const baseUrl = "http://www.websequencediagrams.com"

type Response struct {
  ImageUrl string `json:"img"`
  Errors []string `json:"errors"`
}

func main()  {
  var wg sync.WaitGroup

  for _, filename := range(os.Args[1:]) {
    wg.Add(1)
    fmt.Printf("%s --> %s.png\n", filename, filename)
    go sequence(filename, fmt.Sprintf("%s.png", filename), &wg)
  }

  wg.Wait()
}

func sequence(source string, destination string, wg *sync.WaitGroup){
  var jsonResponse = Response{}
  defer wg.Done()

  sourceBytes, err := ioutil.ReadFile(source)
  if err != nil {
    log.Fatal(err)
  }

  res, err := http.PostForm(fmt.Sprintf("%s/index.php", baseUrl), url.Values{"style": {"default"}, "message": {string(sourceBytes)}, "apiVersion": {"1"}, "format": {"png"}})
  if err != nil {
    log.Fatal(err)
  }

  err = json.NewDecoder(res.Body).Decode(&jsonResponse)
  if err != nil {
    log.Fatal(err)
  }

  if(len(jsonResponse.Errors) > 0) {
    for _, msg := range(jsonResponse.Errors) {
      log.Fatal(msg)
    }
  }

  res, err = http.Get(fmt.Sprintf("%s/%s", baseUrl, jsonResponse.ImageUrl))
  if err != nil {
    log.Fatal(err)
  }

  defer res.Body.Close()
  file, err := os.Create(destination)
  if err != nil {
    log.Fatal(err)
  }

  _, err = io.Copy(file, res.Body)
  if err != nil {
    log.Fatal(err)
  }
  file.Close()
}
