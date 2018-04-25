package action

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const DOUBAN_API_URL string = "https://frodo.douban.com/api/v2/user/%s/interests?apikey=%s&type=%s&status=%s&count=%s&start=%s&_sig=%s&_ts=%s"
const DOUBAN_API_CODE string = "0dad551ec0f84ed02907ff5c42e8ec70"
const DOUBAN_API_UID string = "2539792"
const DOUBAN_API_SIG string = "Dcu02wdHVuAHjQjV6mr5ZVduSQg=" //签名
const DOUBAN_API_TS string = "1524554244"                    //签名时间

var douban_types = []string{"app", "book", "game", "movie", "music"}

type ResultObj struct {
	Start    int32     `json:"start"`
	Count    int32     `json:"count"`
	Total    int32     `json:"total"`
	Subjects []Subject `json:"subjects"`
}

type Subject struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Cover  string  `json:"cover"`
	Url    string  `json:"url"`
	Rating int32   `json:"rating"` //打分
	Star   float32 `json:"star"`   //星级
	Type   string  `json:"type"`
	Status string  `json:"status"`
}

type Data struct {
	Count     int32 `json:"count"`
	Start     int32 `json:"start"`
	Total     int32 `json:"total"`
	Interests []struct {
		Rating struct {
			Value int32 `json:"value"` //打分
		} `json:"rating"`
		Subject struct {
			Id    string `json:"id"`    //唯一ID
			Title string `json:"title"` //标题
			Type  string `json:"type"`  //类型
			Url   string `json:"url"`   //豆瓣地址
			Pic   struct {
				Normal string `json:"normal"` //封面
			}
			Rating struct {
				Value float32 `json:"value"` //星级
			} `json:"rating"`
		}
		Status string `json:"status"` //状态
	} `json:"interests"`
}

func Douban(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ptype := vars["type"]
	conn := true
	for _, douban_type := range douban_types {
		if douban_type == ptype {
			conn = false
			break
		}
	}
	// 是否禁止下一步
	if conn {
		http.Error(w, "404 - Page Not Found", 404)
		return
	}

	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}

	count := r.FormValue("count")
	if count == "" {
		count = "10"
	}

	status := r.FormValue("status")

	result := getDoubanData(ptype, count, start, status)
	if result == "" {
		http.Error(w, "500 - Internal Server Error", 500)
	} else {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		fmt.Fprintln(w, handleDoubanData(result))
	}
}

func getDoubanData(ptype string, count string, start string, status string) string {
	client := &http.Client{}
	urls := fmt.Sprintf(DOUBAN_API_URL, DOUBAN_API_UID, DOUBAN_API_CODE, ptype, status, count, start, DOUBAN_API_SIG, DOUBAN_API_TS)
	u, _ := url.Parse(urls)
	q := u.Query()
	u.RawQuery = q.Encode() //转义
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("User-Agent", "com.douban.frodo/5.24.0(132)")
	req.Header.Set("Host", "frodo.douban.com")
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return string(body)
	}
	return ""
}

func handleDoubanData(content string) string {
	var data Data

	result := ResultObj{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	result.Total = data.Total
	result.Count = data.Count
	result.Start = data.Start

	length := len(data.Interests)
	var subjects []Subject

	for i := 0; i < length; i++ {
		var interest = data.Interests[i]
		sub := Subject{}
		sub.Id = interest.Subject.Id
		sub.Title = interest.Subject.Title
		sub.Cover = interest.Subject.Pic.Normal
		sub.Url = interest.Subject.Url
		sub.Star = interest.Subject.Rating.Value
		sub.Rating = interest.Rating.Value
		sub.Type = interest.Subject.Type
		if strings.EqualFold(sub.Type, "tv") {
			sub.Type = "movie"
		}
		sub.Status = interest.Status

		subjects = append(subjects, sub)
	}
	result.Subjects = subjects
	json, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(json)
}
