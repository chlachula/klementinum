/* Example how to create, read and remove a PKLM sample csv file
 */
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	k "github.com/chlachula/klementinum"
	"github.com/chlachula/klementinum/data"
)

var tStat k.TStat

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

	page1 := `<!DOCTYPE html>
<html>
 <head>
  <style>
   .center {
    text-align: center
   }
  </style>
 </head>
<body>
<div class="center">
<a href=\"/\">Home</a><hr/>
<h2>%s</h2>
 <svg width="800" height="800">
  %s
 </svg>
</div>
</body>
</html>`
	svg1 := `<circle cx="400" cy="400" r="370" stroke="green" stroke-width="40" fill="yellow" />`
	svg1 = makeDodecagon(400, 400, 200.0)
	minmax := minMaxString()
	page := fmt.Sprintf(page1, minmax, svg1)
	fmt.Fprintf(w, page)
}
func exitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>Home<br/>
	<a href="/hello">Hello</a><br/>
	<a href="/temp">Temperature</a><br/>
	<a href="/years_avg_temps">Average Temperatures</a><br/>
	<a href="/exit">Exit</a></h1>`)
}
func y_avg_tempsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<a href=\"/\">Home</a><hr/><h1>Average temperatures</h1>")
	fmt.Fprintf(w, avgTempsString())
	fmt.Fprintf(w, k.svgAverage())

}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<a href=\"/\">Home</a><hr/><h1>Hello World!</h1>")
}

func main() {
	fmt.Println("Start of program")
	tStat = k.TemperatureStatistics(data.TemperatureRecords())

	//mux := http.NewServeMux()
	http.HandleFunc("/years_avg_temps", y_avg_tempsHandler)
	http.HandleFunc("/temp", temperatureHandler)
	http.HandleFunc("/exit", exitHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("End of program")
}
