/* Example how to create, read and remove a PKLM sample csv file
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
	a := sum / float32(count)
	s += fmt.Sprintf("\n<br/><br/>Average temperature %.2f<br/>\n", a)
	for i, avgT := range tStat.YearTavg {
		year := i + tStat.Year1
		delta := avgT - a
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
	/*	if r.URL.Path != "/hello" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
	*/
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	svg1 := `<circle cx="400" cy="400" r="370" stroke="green" stroke-width="40" fill="yellow" />`
	svg1 = `\n<svg width="800" height="800">\n` + makeDodecagon(400, 400, 200.0) + ` \n</svg>\n`
	minmax := minMaxString()
	page := fmt.Sprintf(page1, "temperature", "<a href=\"/\">Home</a><hr/>", minmax, svg1)
	fmt.Fprintf(w, page)
}

//go:embed to_embed/img/Klementinum2023-0112-0953-1465-600x800dpi72q40.jpg
var embededSingleImage []byte

//go:embed to_embed/img/*.jpg
var embededImgDir embed.FS

func embedHandler(w http.ResponseWriter, r *http.Request) {
	img := "\n<img src=\"/embeded_single_Klementinum_image\" />"
	page := fmt.Sprintf(page1, "picture", homeLink, "Klementinum tower in Prague, Czechia", img)
	fmt.Fprintf(w, page)
}

func embededSingleImageHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Write(bytes.NewBufferString("Content-Type: image/jpeg"))
	rw.Header().Write(bytes.NewBufferString("Content-Length: " + strconv.Itoa(len(embededSingleImage))))
	rw.Write(embededSingleImage)
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
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		print("... exiting via web / exit button")
		os.Exit(0)
	} else {
		page := fmt.Sprintf(page1, "home", "<a href=\"/\">Home</a><hr/>", "Home",
			`      <h1><a href="/embeded">Klementinum tower picture</a><br/>
		<a href="/temp">Temperature</a><br/>
		<a href="/years_avg_temps">Average Temperatures</a></h1>
		`+formExit)
		fmt.Fprintf(w, page)
	}
}
func y_avg_tempsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<a href=\"/\">Home</a><hr/><h1>Average temperatures</h1>")
	fmt.Fprintf(w, avgTempsString())
	fmt.Fprintf(w, "<br/>\n")
	fmt.Fprintf(w, k.SVG_average())

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
	http.HandleFunc("/temp", temperatureHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Printf("Starting server at port 8080 ...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("End of program")
}
