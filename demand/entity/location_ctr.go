package entity

import "io"

//LocationCTR interface of creatives ctr per page and seat
type LocationCTR interface {
	//CreativeLocationID return location per creative ID
	CreativeLocationID() int64
	//SeatID return seat id of location
	SeatID() int64
	//PublisherPageID return publisher page id
	PublisherPageID() int64
	//CreativeID return creative id
	CreativeID() int64
	//CreativeSize return creative size of location
	CreativeSize() int64
	//ActiveDays return number of active days of location
	ActiveDays() int64
	//TotalImp return total impression od location
	TotalImp() int64
	//TotalClicks return total clicks od location
	TotalClicks() int64
	//TotalCTR return total CTR od location
	TotalCTR() int64
	//YesterdayImp return yesterday impression od location
	YesterdayImp() int64
	//YesterdayClicks return yesterday clicks od location
	YesterdayClicks() int64
	//YesterdayCTR return yesterday CTR od location
	YesterdayCTR() int64
	//TodayImp return today impression od location
	TodayImp() int64
	//TodayClicks return today clicks od location
	TodayClicks() int64
	//TodayCTR return today CTR od location
	TodayCTR() int64
	// Decode try to decode object from io reader
	Decode(r io.Reader) error
	// Encode is the encode function for serialize object in io writer
	Encode(w io.Writer) error
}
