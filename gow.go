package main

import (
	"fmt"
	//"bufio"
	"encoding/json"
	//"fmt"
	//"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
	//"net/url"
	"os"
	//"os/exec"
	//"os/user"
	//"path/filepath"
	"strconv"
	"time"
)

const (
	appID     = "6ea5cf4906d41b9e4108e66624861864"
	url       = "http://api.openweathermap.org/data/2.5/"
	units     = "metric"
	metric    = "°C"
	mode      = "json"
	spaces    = 2
	dayLen    = 3 + spaces
	minMaxLen = 2 + spaces
	delim     = " "
)

var (
	app = kingpin.New("gow", "A command-line open weather map application written in Golang.")

	nowCmd                    = app.Command("now", "Get the current weather data").Default()
	nowArgCity                = nowCmd.Arg("city", "city to get the weather").String()
	nowFlagShowCoords         = nowCmd.Flag("show-coords", "show coordinates of the location given").Default("true").Bool()
	nowFlagShowWeatherMeta    = nowCmd.Flag("show-weather-meta", "show weather strings, icons, etc.").Default("true").Bool()
	nowFlagShowMain           = nowCmd.Flag("show-main", "show main data abouth the weather (temp, pressure, humidity").Default("true").Bool()
	nowFlagShowWind           = nowCmd.Flag("show-wind", "show data about the wind").Default("true").Bool()
	nowFlagShowClouds         = nowCmd.Flag("show-clouds", "show data about the clouds").Default("true").Bool()
	nowFlagShowAdditionalData = nowCmd.Flag("show-additional-data", "show data about the country, sunrise, sunset, message").Default("true").Bool()

	forecastCmd         = app.Command("forecast", "Get the forecast with weather lowest and highest temperatures")
	forecastArgCity     = forecastCmd.Arg("city", "the citty to get the weather of.").String()
	forecastFlagDays    = forecastCmd.Flag("days", "the citty to get the weather of.").Default("7").Short('d').Int()
	forecastFlagMinimal = forecastCmd.Flag("minimal", "the citty to get the weather of.").Short('m').Bool()
	forecastFlagWeekday = forecastCmd.Flag("weekday", "the citty to get the weather of.").Short('w').Bool()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case nowCmd.FullCommand():
		printNow(
			getNow(*nowArgCity),
			nowFlagShowCoords,
			nowFlagShowWeatherMeta,
			nowFlagShowMain,
			nowFlagShowWind,
			nowFlagShowClouds,
			nowFlagShowAdditionalData,
		)

	case forecastCmd.FullCommand():
		if *forecastFlagMinimal {
			printForecastMinimal(
				getForecast(*forecastArgCity, *forecastFlagDays),
				*forecastFlagMinimal,
				*forecastFlagWeekday,
			)
		} else {
			printForecast(
				getForecast(*forecastArgCity, *forecastFlagDays),
				*forecastFlagMinimal,
				*forecastFlagWeekday,
			)
		}
	}
}

func getNow(city string) weather {
	var w weather
	json.Unmarshal(
		query(
			makeRequest("weather", city, ""),
			""),
		&w,
	)
	return w
}

func getForecast(city string, days int) forecast {
	var f forecast
	json.Unmarshal(
		query(
			makeRequest("forecast/daily", city, "&cnt="+strconv.Itoa(days)),
			""),
		&f,
	)
	return f
}

func query(URL, method string) []byte {
	client := &http.Client{}
	if method == "" {
		method = "GET"
	}
	request, _ := http.NewRequest(method, URL, nil)
	res, _ := client.Do(request)
	ret, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return ret
}

func makeRequest(method, city, param string) string {
	req := url
	req += method
	req += "?q=" + city
	req += "&mode=" + mode
	req += "&units=" + units
	req += param
	req += "&appid=" + appID
	return req
}

func printForecast(f forecast, minimal, weekday bool) {
	forecastLine := ""
	dtformat := "2"
	if weekday {
		dtformat = "Mon"
	}
	for _, key := range f.List {
		forecastLine += "datetime:	" + time.Unix(int64(key.Dt), 0).Format(dtformat) + "\n"
		forecastLine += "morning:	" + strconv.FormatFloat(key.Temp.Morn, 'f', -1, 64) + "\n"
		forecastLine += "day:		" + strconv.FormatFloat(key.Temp.Day, 'f', -1, 64) + "\n"
		forecastLine += "evening:	" + strconv.FormatFloat(key.Temp.Eve, 'f', -1, 64) + "\n"
		forecastLine += "night:		" + strconv.FormatFloat(key.Temp.Night, 'f', -1, 64) + "\n"
		forecastLine += "\n"
	}
	fmt.Println(forecastLine)
}

func printForecastMinimal(f forecast, minimal, weekday bool) {
	forecastLine := ""
	dtformat := "2"
	if weekday {
		dtformat = "Mon"
	}
	for _, key := range f.List {
		forecastLine += time.Unix(int64(key.Dt), 0).Format(dtformat) + ":	"
		forecastLine += strconv.FormatFloat(key.Temp.Morn, 'f', -1, 64) + "	"
		forecastLine += strconv.FormatFloat(key.Temp.Day, 'f', -1, 64) + "	"
		forecastLine += strconv.FormatFloat(key.Temp.Eve, 'f', -1, 64) + "	"
		forecastLine += strconv.FormatFloat(key.Temp.Night, 'f', -1, 64) + "	"
		for _, value := range key.Weather {
			forecastLine += value.Main + "	"
		}
		forecastLine += "\n"
	}
	fmt.Println(forecastLine)
}

func printNow(
	w weather,
	nowFlagShowCoords *bool,
	nowFlagShowWeatherMeta *bool,
	nowFlagShowMain *bool,
	nowFlagShowWind *bool,
	nowFlagShowClouds *bool,
	nowFlagShowAdditionalData *bool,
) {
	nowLine := ""
	if *nowFlagShowCoords {
		nowLine += "lat-lng: " + strconv.FormatFloat(w.Coord.Lat, 'f', -1, 64) + " " + strconv.FormatFloat(w.Coord.Lon, 'f', -1, 64) + "\n"
	}
	if *nowFlagShowWeatherMeta {
	}
	if *nowFlagShowMain {
		nowLine += "humidity: " + strconv.Itoa(w.Main.Humidity) + "\n"
		nowLine += "pressure: " + strconv.Itoa(w.Main.Pressure) + "\n"
		nowLine += "temp: " + strconv.Itoa(w.Main.Temp) + "\n"
		nowLine += "temp-min: " + strconv.Itoa(w.Main.TempMin) + "\n"
		nowLine += "temp-max: " + strconv.Itoa(w.Main.TempMax) + "\n"
	}
	if *nowFlagShowWind {
		nowLine += "wind-speed: " + strconv.FormatFloat(w.Wind.Speed, 'f', -1, 64) + "\n"
		nowLine += "wind-degrees: " + strconv.Itoa(w.Wind.Deg) + "\n"
	}
	if *nowFlagShowClouds {
		nowLine += "clouds-all: " + strconv.Itoa(w.Clouds.All) + "\n"
	}
	if *nowFlagShowAdditionalData {
		nowLine += "type: " + strconv.Itoa(w.Sys.Type) + "\n"
		nowLine += "message: " + strconv.FormatFloat(w.Sys.Message, 'f', -1, 64) + "\n"
		nowLine += "conutry: " + w.Sys.Country + "\n"
		nowLine += "sunrise: " + time.Unix(int64(w.Sys.Sunrise), 0).Format("15:04") + "\n"
		nowLine += "sunset: " + time.Unix(int64(w.Sys.Sunset), 0).Format("15:04") + "\n"
	}

	fmt.Println(nowLine)
}

type weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     int `json:"temp"`
		Pressure int `json:"pressure"`
		Humidity int `json:"humidity"`
		TempMin  int `json:"temp_min"`
		TempMax  int `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

type forecast struct {
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
	} `json:"city"`
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Temp struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity int     `json:"humidity"`
		Weather  []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Speed  float64 `json:"speed"`
		Deg    int     `json:"deg"`
		Clouds int     `json:"clouds"`
		Snow   float64 `json:"snow"`
		Rain   float64 `json:"rain,omitempty"`
	} `json:"list"`
}
