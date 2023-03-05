/* Example how to create, read and remove a PKLM sample csv file
 */
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
)

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
	page := fmt.Sprintf(page1, "2022-0815-1028", svg1)
	fmt.Fprintf(w, page)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>Home<br/><a href="/hello">Hello</a><br/><a href="/temp">Temperature</a></h1>`)
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<a href=\"/\">Home</a><hr/><h1>Hello World!</h1>")
}

func main() {
	fmt.Println("Start of program")

	//mux := http.NewServeMux()
	http.HandleFunc("/temp", temperatureHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("End of program")
}
