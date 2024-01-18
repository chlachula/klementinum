package data

//go test -v  # Verbose test
import "testing"

func gregorianToJulianDateNumber(Y int, M int, D int) int {
	JDN := (1461*(Y+4800+(M-14)/12))/4 + (367*(M-2-12*((M-14)/12)))/12 - (3*((Y+4900+(M-14)/12)/100))/4 + D - 32075
	return JDN
}
func TestDataLen(t *testing.T) {
	//length test
	want := 90215
	d := TemperatureRecords()
	if got := len(d); got != want {
		t.Errorf("Array length is = %d, want %d", got, want)
	}
}

func TestDataFirstRecord(t *testing.T) {
	//the first record test
	wantY := 1775
	wantM := 1
	wantD := 1
	d := TemperatureRecords()
	r := d[0]
	if r.Y != wantY || r.M != wantM || r.D != wantD {
		t.Errorf("The first record date is = %04d-%02d-%02d, want %04d-%02d-%02d", r.Y, r.Y, r.D, wantY, wantM, wantD)
	}
}
func TestDataLastRecord(t *testing.T) {
	//the first record test
	wantY := 2021
	wantM := 12
	wantD := 31
	d := TemperatureRecords()
	i := len(d) - 1
	r := d[i]
	if r.Y != wantY || r.M != wantM || r.D != wantD {
		t.Errorf("The last record date is = %04d-%02d-%02d, want %04d-%02d-%02d", r.Y, r.M, r.D, wantY, wantM, wantD)
	}
}
func TestDataAllRecords(t *testing.T) {
	d := TemperatureRecords()
	// this is just test of the range functionality
	for i, r := range d {
		if r.Y != d[i].Y || r.M != d[i].M || r.D != d[i].D {
			t.Errorf("Unexpected record date %04d-%02d-%02d with index %d", r.Y, r.M, r.D, i)
			return
		}
	}
	r := d[0]
	//sequence date test
	jdnPrev := gregorianToJulianDateNumber(r.Y, r.M, r.D)
	for i := 1; i < len(d)-1; i++ {
		r = d[i]
		jdn := gregorianToJulianDateNumber(r.Y, r.M, r.D)
		if jdn != jdnPrev+1 {
			t.Errorf("Bronken julian date sequence at date %04d-%02d-%02d, expected: %d, got: %d", r.Y, r.M, r.D, (jdnPrev + 1), jdn)
			return
		}
		jdnPrev = jdn
	}
}
