package action

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const XINGZHE_API_URL string = "http://www.imxingzhe.com/api/v4/%s"
const XINGZHE_API_URL_INFO string = "account/get_user_info?userid=%s"
const XINGZHE_API_URL_WORKOUT string = "get_user_workouts?user_id=%s&start=%s&count=%s"
const XINGZHE_API_UID string = "20937"

/*
输出对象格式
*/
type X_OUT_INFO struct {
	Username string `json:"username"`  //昵称
	Ulevel   int32  `json:"ulevel"`    //行者等级
	Avatar   string `json:"avatar"`    //头像
	Credits  int32  `json:"credits"`   //积分
	MainTeam string `json:"main_team"` //主俱乐部
	City     string `json:"city"`      //所在城市

	TotalDistance      int32 `json:"total_distance"`       //总里程数（米）
	TotalValidDistance int32 `json:"total_valid_distance"` //有效总里程
	MonthDistance      int32 `json:"month_distance"`       //月里程数
	MonthValidDistance int32 `json:"month_valid_distance"` //有效月里程

	TotalHots      int32 `json:"total_hots"`       //总热度
	MonthHots      int32 `json:"month_hots"`       //当月热度
	CycleMonthHots int32 `json:"cycle_month_hots"` //当月骑行热度
	RunMonthHots   int32 `json:"run_month_hots"`   //当月跑步热度
	WalkMonthHots  int32 `json:"walk_month_hots"`  //当月徒步热度
}

type X_IN_INFO struct {
	Username string `json:"username"`  //昵称
	Ulevel   int32  `json:"ulevel"`    //行者等级
	Avatar   string `json:"avatar"`    //头像
	Credits  int32  `json:"credits"`   //积分
	MainTeam string `json:"main_team"` //主俱乐部
	City     string `json:"city"`      //所在城市

	TotalDistance      int32 `json:"total_distance"`       //总里程数（米）
	ValidDistance      int32 `json:"valid_distance"`       //有效总里程
	MonthDistance      int32 `json:"month_distance"`       //月里程数
	MonthValidDistance int32 `json:"month_valid_distance"` //有效月里程

	Hots           int32 `json:"hots"`
	HotsMonth      int32 `json:"hots_month"`
	HotsMonthCycle int32 `json:"hots_month_cycle"`
	HotsMonthRun   int32 `json:"hots_month_run"`
	HotsMonthWalk  int32 `json:"hots_month_walk"`
}

type X_OUT_WORKOUT struct {
	Id        int32   `json:"id"`         //轨迹编号
	Title     string  `json:"title"`      //标题
	Sport     string  `json:"sport"`      //运动类型，2跑步，3骑行
	Calories  int32   `json:"calories"`   //卡路里
	Distance  int32   `json:"distance"`   //里程
	AvgSpeed  float32 `json:"avg_speed"`  //平均速度
	StartTime int64   `json:"start_time"` //运动开始时间
	EndTime   int64   `json:"end_time"`   //运动结束时间

	AvgAltitude int32 `json:"avg_altitude"` //平均海拔
	MaxAltitude int32 `json:"max_altitude"` //最高海拔

	MinGrade int32 `json:"min_grade"` //最小坡度
	MaxGrade int32 `json:"max_grade"` //最大坡度

	ElevationLoss int32 `json:"elevation_loss"` //累计下降
	ElevationGain int32 `json:"elevation_gain"` //累计上升
}

type X_IN_WORKOUT struct {
	Id        int32   `json:"id"`         //轨迹编号
	Title     string  `json:"title"`      //标题
	Sport     int32   `json:"sport"`      //运动类型，2跑步，3骑行
	Calories  int32   `json:"calories"`   //卡路里
	Distance  int32   `json:"distance"`   //里程
	AvgSpeed  float32 `json:"avg_speed"`  //平均速度
	StartTime int64   `json:"start_time"` //运动开始时间
	EndTime   int64   `json:"end_time"`   //运动结束时间

	AvgAltitude int32 `json:"avg_altitude"` //平均海拔
	MaxAltitude int32 `json:"max_altitude"` //最高海拔

	MinGrade int32 `json:"min_grade"` //最小坡度
	MaxGrade int32 `json:"max_grade"` //最大坡度

	ElevationLoss int32 `json:"elevation_loss"` //累计下降
	ElevationGain int32 `json:"elevation_gain"` //累计上升
}

func Xingzhe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ptype := vars["type"]

	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}

	count := r.FormValue("count")
	if count == "" {
		count = "10"
	}

	var result string
	if ptype == "info" {
		result = handleUserInfoData(getUserInfo())
	} else if ptype == "workout" {
		result = handleUserWorkoutData(getUserWorkout(start, count))
	} else {
		http.Error(w, "404 - Page Not Found", 404)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	fmt.Fprintln(w, result)
}

func getUserInfo() string {
	client := &http.Client{}
	url := fmt.Sprintf(XINGZHE_API_URL_INFO, XINGZHE_API_UID)
	url = fmt.Sprintf(XINGZHE_API_URL, url)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return string(body)
	}
	return ""
}

func handleUserInfoData(content string) string {
	var data X_IN_INFO
	result := X_OUT_INFO{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	result.Username = data.Username
	result.Ulevel = data.Ulevel
	result.Avatar = data.Avatar
	result.Credits = data.Credits
	result.MainTeam = data.MainTeam
	result.City = data.City

	result.TotalDistance = data.TotalDistance
	result.MonthDistance = data.MonthDistance
	result.TotalValidDistance = data.ValidDistance
	result.MonthValidDistance = data.MonthValidDistance

	result.TotalHots = data.Hots
	result.MonthHots = data.HotsMonth
	result.CycleMonthHots = data.HotsMonthCycle
	result.RunMonthHots = data.HotsMonthRun
	result.WalkMonthHots = data.HotsMonthWalk

	json, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(json)
}

func getUserWorkout(start string, count string) string {
	client := &http.Client{}
	url := fmt.Sprintf(XINGZHE_API_URL_WORKOUT, XINGZHE_API_UID, start, count)
	url = fmt.Sprintf(XINGZHE_API_URL, url)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return string(body)
	}
	return ""
}

func handleUserWorkoutData(content string) string {
	var data []X_IN_WORKOUT
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	length := len(data)
	var workouts []X_OUT_WORKOUT

	for i := 0; i < length; i++ {
		var item = data[i]
		wo := X_OUT_WORKOUT{}

		wo.Id = item.Id
		wo.Title = item.Title
		var sport = item.Sport
		if sport == 2 {
			wo.Sport = "跑步"
		} else if sport == 3 {
			wo.Sport = "骑行"
		} else {
			wo.Sport = "未知"
		}
		wo.Calories = item.Calories
		wo.Distance = item.Distance
		wo.AvgSpeed = item.AvgSpeed
		wo.StartTime = item.StartTime
		wo.EndTime = item.EndTime

		wo.AvgAltitude = item.AvgAltitude
		wo.MaxAltitude = item.MaxAltitude

		wo.MinGrade = item.MinGrade
		wo.MaxGrade = item.MaxGrade

		wo.ElevationLoss = item.ElevationLoss
		wo.ElevationGain = item.ElevationGain

		workouts = append(workouts, wo)
	}

	json, err := json.Marshal(workouts)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(json)
}
