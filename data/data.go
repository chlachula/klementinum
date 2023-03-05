package data

/*
data source https://www.chmi.cz/files/portal/docs/meteo/ok/files/PKLM_pro_portal.xlsx

The first line with column names was removed
*/
type TempRecord struct {
	y int
	m int
	d int
	t float32
}
