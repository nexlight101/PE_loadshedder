package modules

import (
	"fmt"
	"html/template"
	"log"
	"sort"
	"strings"
	"time"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nexlight101/webshedder/v3"
)

// shortform declares short form of standard time
const shortForm = "2006-01-02"

// Controller struct for template controller
type Controller struct {
	tpl *template.Template
}

// Sched caters for all stages + times
type Sched struct {
	Stage string
	Group string
	Times []string
}

var (
	date           string
	area           string
	areas          []webshedder.Area
	areaM          map[string][]string
	schedules      []webshedder.Schedule
	s1, s2, s3, s4 string
)

func init() {
	schedules, areas = webshedder.ReadJSON(webshedder.Filename1, webshedder.Filename2)
	// printSchedule(schedules)
	areaM = webshedder.BuildMap(areas)
}

// NewController provides new controller for template processing
func NewController(t *template.Template) *Controller {
	return &Controller{t}
}

// // printSchedule prints all schedules
// func printSchedule(s []webshedder.Schedule) {
// 	f, _ := os.Create("data/out.txt")
// 	defer f.Close()
// 	var txt string
// 	for index, schedule := range s {
// 		txt += fmt.Sprintf(" Record : %d, Date : %v , Group: %v\n", index, schedule.Date, schedule.Group)
// 	}
// 	f.WriteString(txt)
// }

// Main provides a method for the root page
func (c Controller) Main(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// populate the template struct with empty values
	templateData := struct {
		BGround string
	}{
		BGround: "bg",
	}
	c.tpl.ExecuteTemplate(w, "landing.gohtml", templateData)
}

// Index GET ROUTE dislays form to read a date, stage and suburb.
func (c Controller) Index(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println("Find ROUTE activated")
	// Create the area slice
	suburbs := []string{}
	for _, v := range areas {
		suburbs = append(suburbs, v.AreaName...)
	}
	suburbs = cleanArea(suburbs)
	sort.Strings(suburbs)

	// populate the template struct with empty values
	templateData := struct {
		Areas   []string
		BGround string
	}{
		Areas:   suburbs,
		BGround: "in",
	}
	c.tpl.ExecuteTemplate(w, "index.gohtml", templateData)

}

// Schedule GET ROUTE finds the schedule.
func (c Controller) Schedule(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	var (
		date1   time.Time
		groupX  []string
		schedX1 []Sched
		schedX2 []Sched
	)
	fmt.Println("Schedule ROUTE activated")
	// format time

	date1, err := time.Parse(shortForm, date)
	if err != nil {
		log.Printf("Cannot parse time: %v\n", err)
	}
	// Lookup the group
	groupX = areaM[area]
	// Change value from on to stage n
	if s1 == "on" {
		s1 = "Stage 1"
		sched1 := buildSched(&date1, &s1, groupX, schedules)
		if len(groupX) == 1 {
			schedX1 = append(schedX1, sched1...)
		} else {
			schedX1 = append(schedX1, sched1[0])
			schedX2 = append(schedX2, sched1[1])
		}
	}
	if s2 == "on" {
		s2 = "Stage 2"
		sched2 := buildSched(&date1, &s2, groupX, schedules)
		if len(groupX) == 1 {
			schedX1 = append(schedX1, sched2...)
		} else {
			schedX1 = append(schedX1, sched2[0])
			schedX2 = append(schedX2, sched2[1])
		}
	}
	if s3 == "on" {
		s3 = "Stage 3"
		sched3 := buildSched(&date1, &s3, groupX, schedules)
		if len(groupX) == 1 {
			schedX1 = append(schedX1, sched3...)
		} else {
			schedX1 = append(schedX1, sched3[0])
			schedX2 = append(schedX2, sched3[1])
		}
	}
	if s4 == "on" {
		s4 = "Stage 4"
		sched4 := buildSched(&date1, &s4, groupX, schedules)
		if len(groupX) == 1 {
			schedX1 = append(schedX1, sched4...)
		} else {
			schedX1 = append(schedX1, sched4[0])
			schedX2 = append(schedX2, sched4[1])
		}
	}
	if s1 == "" && s2 == "" && s3 == "" && s4 == "" {
		noDataError := "You have not select all necessary fields!"
		templateData3 := struct {
			FileError string
		}{
			FileError: noDataError,
		}
		c.tpl.ExecuteTemplate(w, "error.gohtml", templateData3)
		return
	}
	templateData1 := struct {
		Area       string
		Date       string
		Group1     string
		Schedules1 []Sched
		BGround    string
	}{
		Area:       area,
		Date:       date,
		Group1:     groupX[0],
		Schedules1: schedX1,
		BGround:    "sc",
	}
	if len(groupX) == 1 {
		c.tpl.ExecuteTemplate(w, "schedule1.gohtml", templateData1)
	} else {
		// populate the template struct with values
		templateData2 := struct {
			Area       string
			Date       string
			Group1     string
			Group2     string
			Schedules1 []Sched
			Schedules2 []Sched
			BGround    string
		}{
			Area:       area,
			Date:       date,
			Group1:     groupX[0],
			Group2:     groupX[1],
			Schedules1: schedX1,
			Schedules2: schedX2,
			BGround:    "sc",
		}
		c.tpl.ExecuteTemplate(w, "schedule2.gohtml", templateData2)
	}
}

// FindArea POST ROUTE finds the schedule.
func (c Controller) FindArea(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println("Find Area ROUTE activated")
	// If POST request
	if req.Method == http.MethodPost {
		// Get fieldsfrom form
		s1 = req.FormValue("Stage 1")
		s2 = req.FormValue("Stage 2")
		s3 = req.FormValue("Stage 3")
		s4 = req.FormValue("Stage 4")

		date = req.FormValue("datepicker")
		area = req.FormValue("areas")
		fmt.Printf("Request: %s Date : %v\n", area, date)
		fmt.Printf("Stage 1: %s Stage 2: %s Stage 3: %s Stage 4: %s \n", s1, s2, s3, s4)

		http.Redirect(w, req, "/schedule", http.StatusSeeOther)
		return
	}

	http.Redirect(w, req, "/index", http.StatusSeeOther)
}

// buildSched builds the schedule struct per stage to pass into schedule template
func buildSched(t *time.Time, st *string, g []string, s []webshedder.Schedule) []Sched {
	groupsX := webshedder.SearchTimes(t, st, g, s)
	schedX := []Sched{}
	for _, v := range groupsX {
		sched1 := Sched{
			Stage: *st,
			Group: v.Group,
			Times: v.Times,
		}
		schedX = append(schedX, sched1)
	}
	// add some humour
	for _, v := range schedX {
		if v.Times[0] == "" && len(v.Times) < 2 {
			v.Times[0] = "Consider yourself lucky!"
		} else {
			for i, j := range v.Times {
				v.Times[i] = strings.TrimPrefix(j, " ")
			}
		}
	}

	return schedX
}

// cleanArea takes out duplicates
func cleanArea(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

// About provides a method for the information page
func (c Controller) About(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// populate the template struct with empty values

	templateData := struct {
		BGround string
	}{
		BGround: "ab",
	}
	c.tpl.ExecuteTemplate(w, "about.gohtml", templateData)
}
