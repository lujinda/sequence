package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
  "net/url"
  "encoding/json"
  "sync"
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
    sourceBytes, err := ioutil.ReadFile(filename)
    if err != nil {
      panic(err);
    }
    wg.Add(1)
    fmt.Printf("%s --> %s.png\n", filename, filename)
    go sequence(string(sourceBytes), fmt.Sprintf("%s.png", filename), &wg)
  }

  wg.Wait()
}

func sequence(source string, destination string, wg *sync.WaitGroup) (err error){
  var jsonResponse = Response{}

  defer wg.Done()

  res, err := http.PostForm(fmt.Sprintf("%s/index.php", baseUrl), url.Values{"style": {"default"}, "message": {source}, "apiVersion": {"1"}, "format": {"png"}})
  if err != nil {
    return err;
  }

  err = json.NewDecoder(res.Body).Decode(&jsonResponse)
  if err != nil {
    return err;
  }

  res, err = http.Get(fmt.Sprintf("%s/%s", baseUrl, jsonResponse.ImageUrl))
  if err != nil {
    return err;
  }

  defer res.Body.Close()
  file, err := os.Create(destination)
  if err != nil {
      return err
  }

  _, err = io.Copy(file, res.Body)
  if err != nil {
      return err
  }
  file.Close()

  return nil
}
