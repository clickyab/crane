package workers

const sql = "INSERT INTO ads  () VALUES () ON DUPLICATE KEY UPDATE SET a=5,b=12,c=VALUES(c)+423 "

type tableModel struct {
	Time          int
	Request       int
	Impression    int
	Win           int
	Show          int
	ImpressionBid int64
	ShowBid       int64
	Supplier      string
	Source        string
	Demand        string
}

func flush(supDemSrc map[string]tableModel, supSrc map[string]tableModel) error {
	// TODO insert into both tables
	return nil
}
