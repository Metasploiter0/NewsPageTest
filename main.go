package main

import (
  "net/http"
  "html/template"
  "os"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "time"
  "strings"

  "github.com/sirupsen/logrus"
)

type DataStruct struct {
  Title string 
  Descr string
  URLtoImage string
  ReadMore string
}

type Response struct {
    Status string `json:"status"`
    TotalResults int `json:"totalResults"`
    Articles[] struct {
        Source struct {
            ID string `json:"id"`
            Name string `json:"name"`
        } `json:"source"`
        Author string `json:"author"`
        Title string `json:"title"`
        Description string `json:"description"`
        URL string `json:"url"`
        URLToImage string `json:"urlToImage"`
        PublishedAt time.Time `json:"publishedAt"`
        Content any `json:"content"`
    } `json:"articles"`
}

const (
  url string = "https://newsapi.org/v2/top-headlines?country=de&category=business&apiKey=f687d1c28b2046d6a3f4850db11a3a0f"
  port string = ":8080"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
  switch r.URL.Path {
    case "/":
      tmpl, err := template.ParseFiles("index.html")
      if err != nil {
        logrus.Error("Error parsing template:", err)
        return
      }
    
      dataGet := func() Response {
        data, err := getNewsData(url)
        if err != nil {
          logrus.Error(err)
          return Response{}
        }
        return data
      }

      RawData := dataGet()
    
      intermedData := []DataStruct{}
      for i := 0; i < len(RawData.Articles); i++ {
        intermedData = append(intermedData, DataStruct{RawData.Articles[i].Title, RawData.Articles[i].Description, fmt.Sprintf("%v", RawData.Articles[i].URLToImage), RawData.Articles[i].URL})
      }
      Data := map[string][]DataStruct{
        "Articles": intermedData,
      }
      err = tmpl.Execute(w, Data)
      if err != nil {
        logrus.Error("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
      }
  
    case "/post-own":
      tmpl, err := template.ParseFiles("post.html")
      if err != nil {
        logrus.Error("Error parsing template:", err)
        return
      }
      err = tmpl.Execute(w, nil)
      if err != nil {
        logrus.Error("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
      }
    default:
      http.NotFound(w, r)
      return
  }
}

func probeURL(url string) []string {
  resp, err := http.Get(url)
  if err != nil {
    return []string{}
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return []string{}
  }
  defer resp.Body.Close()
  fmt.Println(string(body))
  sepStr := strings.Split(string(body), ">")
  return sepStr
}

func main() {
  file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
    logrus.Fatal("Error opening log file:", err)
    return
  }
  defer file.Close()
  logrus.SetOutput(file)

  logrus.Infoln("Starting server...")
  http.HandleFunc("/", requestHandler)
  http.ListenAndServe(port, nil)
}

func getNewsData(url string) (Response, error) {
  resp, err := http.Get(url)
  if err != nil {
    return Response{}, fmt.Errorf(
      "Error when getting data from %s: %v", url, err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return Response{}, fmt.Errorf(
      "Error when reading response body from %s: %v", url, err)
  }
  defer resp.Body.Close()
  var result Response
  if err := json.Unmarshal(body, &result); err != nil {
    return Response{}, fmt.Errorf(
      "Error when unmarshalling response body from %s: %v", url, err)
  }
  return result, nil
}