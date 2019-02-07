package entities

// List for retargeing
type List struct {
	ID     int32  `db:"rl_id"`
	UID    int32  `db:"u_id"`
	KEY    string `db:"rl_key"`
	Domain string `db:"rl_domain"`
	Name   string `db:"rl_name"`
}

// CheckList to be sure target list exists
func CheckList(s string) (*List, error) {
	l := &List{}
	return l, NewManager().GetRDbMap().SelectOne(l, "select rl_id,u_id,rl_key,rl_domain,"+
		"rl_name from retargeting_list where rl_key = ?", s)
}
