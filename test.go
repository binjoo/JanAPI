package main

import (
	"encoding/json"
	"fmt"
)

type ResultObj struct {
	Start    int32     `json:"start"`
	Count    int32     `json:"count"`
	Total    int32     `json:"total"`
	Subjects []Subject `json:"subjects"`
}

type Subject struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	Url    string `json:"url"`
	Rating float32 `json:"rating"`
	Star   float32 `json:"star"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type Data struct {
	Count     int32 `json:"count"`
	Start     int32 `json:"start"`
	Total     int32 `json:"total"`
	Interests []struct {
		Subject struct {
			Id    string `json:"id"`    //唯一ID
			Title string `json:"title"` //标题
			Type  string `json:"type"`  //类型
			Url   string `json:"url"`   //豆瓣地址
			Pic   struct {
				Normal string `json:"normal"` //封面
			}
			Rating struct {
				Star_count float32 `json:"star_count"` //观看打分
				Value      float32 `json:"value"`      //平均分
			} `json:"rating"`
		}
		Status string `json:"status"` //状态
	} `json:"interests"`
}

func handleDoubanData(content string) ResultObj {
	var data Data
	result := ResultObj{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println(err.Error())
		return result
	}
	fmt.Println(data.Total)
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
		sub.Rating = interest.Subject.Rating.Star_count
		sub.Type = interest.Subject.Type
		sub.Status = interest.Status

		subjects = append(subjects, sub)
	}
	result.Subjects = subjects
	fmt.Println(json.Marshal(result))
	return result
}
