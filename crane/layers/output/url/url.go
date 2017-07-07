package url

// Register will register a channel that receive click data for processing/storing
func Register(c chan<- Data) {
	worker = c
}

var worker chan<- Data
