package entities

import "database/sql"

// CTRStat is the ctr stat table
type CTRStat struct {
	Impression1  sql.NullInt64 `db:"imp1"`
	Impression2  sql.NullInt64 `db:"imp2"`
	Impression3  sql.NullInt64 `db:"imp3"`
	Impression4  sql.NullInt64 `db:"imp4"`
	Impression5  sql.NullInt64 `db:"imp5"`
	Impression6  sql.NullInt64 `db:"imp6"`
	Impression7  sql.NullInt64 `db:"imp7"`
	Impression8  sql.NullInt64 `db:"imp8"`
	Impression9  sql.NullInt64 `db:"imp9"`
	Impression10 sql.NullInt64 `db:"imp10"`
	Impression11 sql.NullInt64 `db:"imp11"`
	Impression12 sql.NullInt64 `db:"imp12"`
	Impression13 sql.NullInt64 `db:"imp13"`
	Impression14 sql.NullInt64 `db:"imp14"`
	Impression15 sql.NullInt64 `db:"imp15"`
	Impression16 sql.NullInt64 `db:"imp16"`
	Impression17 sql.NullInt64 `db:"imp17"`
	Impression18 sql.NullInt64 `db:"imp18"`
	Impression19 sql.NullInt64 `db:"imp19"`
	Impression20 sql.NullInt64 `db:"imp20"`
	Impression21 sql.NullInt64 `db:"imp21"`

	Click1  sql.NullInt64 `db:"click1"`
	Click2  sql.NullInt64 `db:"click2"`
	Click3  sql.NullInt64 `db:"click3"`
	Click4  sql.NullInt64 `db:"click4"`
	Click5  sql.NullInt64 `db:"click5"`
	Click6  sql.NullInt64 `db:"click6"`
	Click7  sql.NullInt64 `db:"click7"`
	Click8  sql.NullInt64 `db:"click8"`
	Click9  sql.NullInt64 `db:"click9"`
	Click10 sql.NullInt64 `db:"click10"`
	Click11 sql.NullInt64 `db:"click11"`
	Click12 sql.NullInt64 `db:"click12"`
	Click13 sql.NullInt64 `db:"click13"`
	Click14 sql.NullInt64 `db:"click14"`
	Click15 sql.NullInt64 `db:"click15"`
	Click16 sql.NullInt64 `db:"click16"`
	Click17 sql.NullInt64 `db:"click17"`
	Click18 sql.NullInt64 `db:"click18"`
	Click19 sql.NullInt64 `db:"click19"`
	Click20 sql.NullInt64 `db:"click20"`
	Click21 sql.NullInt64 `db:"click21"`
}
