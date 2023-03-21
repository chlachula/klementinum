/*
 * Web application to display Klementinum temperature records,
 * especially avarage temperatures
 */
package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	k "github.com/chlachula/klementinum"
	"github.com/chlachula/klementinum/data"
)

var tStat k.TStat
var homeLink = "<a href=\"/\">Home</a><hr/>\n"
var formExit string = `<form action="/" method="post" name="exit"><input type="submit" value="Exit"></form>`
var page1 string = `<!DOCTYPE html>
<html>
 <head>
 <title>Klementinum %s</title>
 <link rel="icon" href="/favicon.svg" type="image/svg+xml">
 <link rel="icon" href="/embeded_favicon1" type="ico">
  <style>
   .center {
    text-align: center
   }
  </style>
 </head>
<body>
<div class="center">
%s
<h2>%s</h2> 
  %s
</div>
</body>
</html>`

// func avgTemps(tStat TStat) TAverageDifferences {
func getAverageTemperature(tStat k.TStat) float32 {
	count := 0
	var sum float32
	for _, avgT := range tStat.YearTavg {
		count += 1
		sum += avgT
	}
	return sum / float32(count)
}
func averageTemperatureString() string {
	s := ""
	for i, avgT := range tStat.YearTavg {
		year := i + tStat.Year1
		s += fmt.Sprintf("%d:%.1f ", year, avgT)
	}
	s += fmt.Sprint("\n<br/>\n")
	return s
}
func TempDiffsToAverage(average float32, tStat k.TStat) []float32 {
	diffs := make([]float32, len(tStat.YearTavg))
	for i, avgT := range tStat.YearTavg {
		delta := avgT - average
		diffs[i] = delta
	}
	return diffs
}
func tempsDiffsString(diffs []float32) string {
	s := ""
	for i, delta := range diffs {
		year := i + tStat.Year1
		if delta > 0.0 {
			s += fmt.Sprintf("<b>%d:%.1f</b> ", year, delta)
		} else {
			s += fmt.Sprintf("%d:%.1f ", year, delta)
		}
	}
	return s
}
func avgTempsString() string {
	s := ""
	count := 0
	var sum float32
	for i, avgT := range tStat.YearTavg {
		count += 1
		sum += avgT
		year := i + tStat.Year1
		s += fmt.Sprintf("%d:%.1f ", year, avgT)
	}
	average := sum / float32(count)
	s += fmt.Sprintf("\n<br/><h2>Average temperature %.2f°C<br/>Relative differences in years %d .. %d</h2><br/>\n", average, tStat.Year1, tStat.YearEnd)
	for i, avgT := range tStat.YearTavg {
		year := i + tStat.Year1
		delta := avgT - average
		if delta > 0.0 {
			s += fmt.Sprintf("<b>%d:%.1f</b> ", year, delta)
		} else {
			s += fmt.Sprintf("%d:%.1f ", year, delta)
		}

	}
	return s
}

func minMaxString() string {
	return fmt.Sprintf("%d years since %d to %d: Min %v Max %v", tStat.YearEnd-tStat.Year1+1, tStat.Year1, tStat.YearEnd, tStat.MinT, tStat.MaxT)
}
func makeDodecagon(xs int, ys int, r float64) string {
	s := ""
	c1 := `<circle cx="%d" cy="%d" r="%.1f" stroke="yellow" stroke-width="8" fill="none" />`
	c0 := `<circle cx="%d" cy="%d" r="%.1f" stroke="green" stroke-width="3" fill="none" />`
	s += fmt.Sprintf(c0, xs, ys, r)
	s += fmt.Sprintf(c1, xs, ys, r*0.6)
	s += fmt.Sprintf(c1, xs, ys, r*1.4)
	monthNames := [12]string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}
	///c :=`<circle cx="%.1f" cy="%.1f" r="10" stroke="green" stroke-width="1" fill="yellow" />`
	t := `<text x="%.1f" y="%.1f" fill="red"  dominant-baseline="middle" text-anchor="middle" font-size="xx-large">%s</text>
`
	r1 := r * 1.65
	for i := 0; i < 12; i++ {
		a := float64(math.Pi) * float64(i) / 6.0
		x := float64(xs) + r1*math.Sin(a)
		y := float64(ys) - r1*math.Cos(a)
		//s += fmt.Sprintf(c,x,y)
		s += fmt.Sprintf(t, x, y, monthNames[i])

	}
	return s
}
func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	svg1 := `\n<svg width="800" height="800">\n` + makeDodecagon(400, 400, 200.0) + ` \n</svg>\n`
	minmax := minMaxString()
	page := fmt.Sprintf(page1, "temperature", homeLink, minmax, svg1)
	fmt.Fprint(w, page)
}

//go:embed to_embed/img/Klementinum2023-0112-0953-1465-600x800dpi72q40.jpg
var embededSingleImage []byte

//go:embed to_embed/img/favicon1.ico
var embededFavicon1 []byte

//go:embed to_embed/img/*.jpg
var embededImgDir embed.FS

func embedHandler(w http.ResponseWriter, r *http.Request) {
	img := "\n<img src=\"/embeded_single_Klementinum_image\" />"
	page := fmt.Sprintf(page1, "picture", homeLink, "Klementinum tower in Prague, Czechia", img)
	fmt.Fprint(w, page)
}

func embededSingleImageHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Write(bytes.NewBufferString("Content-Type: image/jpeg"))
	rw.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(embededSingleImage))))
	rw.Write(embededSingleImage)
}
func embededFavicon1Handler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Write(bytes.NewBufferString("Content-Type: image/ico"))
	rw.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(embededFavicon1))))
	rw.Write(embededFavicon1)
}
func getembededSingleImage(rw http.ResponseWriter, r *http.Request) {

	fname := r.URL.Query().Get("name")

	data, err := embededImgDir.ReadFile("images" + fname)
	if err != nil {
		rw.Write([]byte("Error occured : " + err.Error()))
		return
	}

	rw.Header().Write(bytes.NewBufferString("Content-Type: image/jpg"))
	rw.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(data))))
	rw.Write(data)
}
func faviconSvgHandler(w http.ResponseWriter, r *http.Request) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
	viewBox="0 0 16 16" width="16" height="16">
	<title>Favicon</title>
	<circle cx="8" cy="8" r="7" stroke="green" stroke-width="3" fill="yellow" />
	</svg>`
	w.Header().Write(bytes.NewBufferString("Content-Type: image/svg+xml .svg .svgz"))
	w.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(svg))))
	fmt.Fprint(w, svg)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		print("... exiting via web / exit button")
		os.Exit(0)
	} else {
		page := fmt.Sprintf(page1, "home", "Home<hr/>", "",
			`      <h1><a href="/embeded">Klementinum tower picture</a><br/>
		<a href="/temp">Temperature</a><br/>
		<a href="/years_avg_temps">Average Temperatures</a></h1>
		`+formExit)
		fmt.Fprint(w, page)
	}
}
func y_avg_tempsHandler(w http.ResponseWriter, r *http.Request) {
	average := getAverageTemperature(tStat)
	diffs := TempDiffsToAverage(average, tStat)
	page := fmt.Sprintf(page1, "averate temperatures", homeLink,
		fmt.Sprintf("Average temperatures in years %d .. %d", tStat.Year1, tStat.YearEnd),
		averageTemperatureString()+
			fmt.Sprintf("\n<br/><h2>Average temperature %.2f°C<br/>Relative differences in years %d .. %d</h2><br/>\n", average, tStat.Year1, tStat.YearEnd)+
			tempsDiffsString(diffs)+
			"<br/>\n"+k.SVG_average(1000, 400, average, diffs, tStat.MaxT.T-tStat.MinT.T))
	fmt.Fprint(w, page)
}

func main() {
	fmt.Println("Start of program")
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user
	tStat = k.TemperatureStatistics(data.TemperatureRecords())

	//mux := http.NewServeMux()
	http.HandleFunc("/years_avg_temps", y_avg_tempsHandler)
	http.HandleFunc("/embeded", embedHandler)
	//http.Handle("/embeded", http.FileServer(http.FS(embededImgDir)))
	http.HandleFunc("/embeded_single_Klementinum_image", embededSingleImageHandler)
	http.HandleFunc("/embeded_favicon1", embededFavicon1Handler)
	http.HandleFunc("/favicon.svg", faviconSvgHandler)
	http.HandleFunc("/temp", temperatureHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Printf("Starting server at port 8080 ...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("End of program")
}
