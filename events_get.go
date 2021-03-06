package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
	"time"
)

var (
	dDayT         *template.Template
	EventResponse string
)

func init() {
	format := `📅 학교 주요 일정이에요!
{{ range . }}
{{ .MMDD }} {{ .Name }}
{{if not .IsDDAY}}D{{ .LeftDays }}{{else}}D-DAY 🎉{{end}}
{{ end }}`

	dDayT = template.Must(template.New("format").Parse(format))
}

func getEvents() {
	// events.json 파일 가져오기
	file, err := ioutil.ReadFile("events.json")
	if err != nil {
		log.Println(err)
		return
	}

	var events []Event
	// json 해독
	err = json.Unmarshal(file, &events)
	if err != nil {
		log.Println("Error while unmarshal events.json:", err)
		return
	}

	validEvents := make([]Event, 0, 8)

	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	for _, value := range events {
		// yyyy-mm-dd 에서 time.Time 파싱
		parsedDate, err := time.Parse("2006/01/02", value.DateString)
		if err != nil {
			log.Println(err)
			continue
		}

		if value.After == 0 {
			value.MMDD = value.DateString[5:]
		} else {
			value.MMDD = value.DateString[5:] + " - " + parsedDate.Local().AddDate(0, 0, value.After).Format("2006/01/02")[5:]
		}

		value.Date = parsedDate.Local().Add(time.Hour * -9)

		// 지금 마이너스 그날
		left := value.Date.Sub(midnight).Hours()
		if left <= 0 && int(left/24) >= -value.After {
			value.IsDDAY = true
		} else if (left < 0 && int(left/24) < -value.After) || (int(left/24) > 70) {
			continue
		}
		value.LeftDays = -int(left / 24)

		validEvents = append(validEvents, value)

		if len(validEvents) >= 7 {
			break
		}
	}

	if len(validEvents) == 0 {
		EventResponse = "📅 등록되어 있는 일정이 없어요!\\n나중에 다시 확인해주세요."
		return
	}

	var tpl bytes.Buffer
	err = dDayT.Execute(&tpl, validEvents)
	if err != nil {
		log.Println("Error while executing dday get...:", err)
		return
	}

	EventResponse = strings.Replace(tpl.String(), "\n", "\\n", -1)

}
