/* Example how to create, read and remove a PKLM sample csv file
 */
package klementinum

import (
	"fmt"

	"github.com/chlachula/klementinum/data"
)

//all := data.TemperatureRecords()

type TStat = struct {
	Year1    int
	YearEnd  int
	MinT     data.TempRecord
	MaxT     data.TempRecord
	YearTavg []float32
}

func Equal_TStat(a, b TStat) bool {
	if a.Year1 != b.Year1 {
		return false
	}
	if a.YearEnd != b.YearEnd {
		return false
	}
	if a.MinT.D != b.MinT.D || a.MinT.M != b.MinT.M || a.MinT.Y != b.MinT.Y {
		return false
	}
	if a.MaxT.D != b.MaxT.D || a.MaxT.M != b.MaxT.M || a.MaxT.Y != b.MaxT.Y {
		return false
	}
	if len(a.YearTavg) != len(b.YearTavg) {
		return false
	}
	for i, x := range a.YearTavg {
		if x != b.YearTavg[i] {
			return false
		}
	}

	return true
}

type TAverageDifferences = struct {
	Average float32
	Diff    []float32
}

func DaysInYear(y int) int {
	yDays := 365
	if y%4 == 0 && y != 2000 {
		yDays = 366
	}
	return yDays
}
func TemperatureStatistics(allData []data.TempRecord) TStat {
	var s TStat
	s.Year1 = allData[0].Y
	s.YearEnd = allData[len(allData)-1].Y
	s.YearTavg = make([]float32, s.YearEnd-s.Year1+1)
	s.MinT.T = 100.0
	s.MaxT.T = -273.0

	for _, r := range allData {
		s.YearTavg[r.Y-s.Year1] += r.T
		if r.T < s.MinT.T {
			s.MinT = r
		}
		if r.T > s.MaxT.T {
			s.MaxT = r
		}
	}
	for y, _ := range s.YearTavg {
		year := y + s.Year1
		s.YearTavg[y] = s.YearTavg[y] / float32(DaysInYear(year))
	}

	return s
}

func SVG_average(average float32, diffs []float32, tempRange float32) string {
	width := 600
	height := 300
	s := ""
	y0 := float32(height) / 2.0
	deltaX := float32(width) / (float32(len(diffs)) + 1)
	for i, deltaT := range diffs {
		color := "red"
		if deltaT < 0 {
			color = "blue"
		}
		p := float32(height) / tempRange
		deltaY := deltaT * p
		s += fmt.Sprintf("<path d=\"M%.1f,%.1f v%.1fh%.1fv%.1fz \" fill=\"%s\" stroke=\"none\" stroke-width=\"0\" />\n",
			float32(i)*deltaX, y0, -deltaY, deltaX, deltaY, color)
	}
	svgFormat := `<svg
 xmlns="http://www.w3.org/2000/svg" 
 xmlns:xlink="http://www.w3.org/1999/xlink" 
 viewBox="0 0 %d %d" 
 width="%d" height="%d">
 <title>Yearly differences of the Klementinum temperature records</title>
 <defs>
	<style>
	   svg { background-color: lightgray; }
	</style>
	<pattern id="bg_image" patternUnits="userSpaceOnUse" width="100" height="100">
	   <image href="/embeded_single_Klementinum_image" x="0" y="0" width="30" height="40" />
	</pattern>
 </defs>
 <g id="main">
    <rect x="0" y="0" width="600" height="300" fill="none" stroke="none" />
    <text x="10" y="30" >Average year temperature: %.1f°C</text>
 </g>
 %s
</svg>
`
	return fmt.Sprintf(svgFormat, width, height, width, height, average, s)
}
