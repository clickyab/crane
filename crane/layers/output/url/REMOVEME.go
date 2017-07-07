package url

func init() {
	var _ = click{}
	c := make(chan Data)
	Register(c)

}
