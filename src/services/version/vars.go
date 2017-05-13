package version

import "os"

var (
	hash  = "ba34c59446315bd5273178118d8e94040bb1bb8c"
	short = "ba34c59"
	date  = "04-08-17-22-10-22"
	build = "04-11-17-12-12-57"
	count = "115"
)

func init() {
	if o := os.Getenv("LONGHASH"); o != "" {
		hash = o
	}

	if o := os.Getenv("SHORTHASH"); o != "" {
		short = o
	}

	if o := os.Getenv("COMMITDATE"); o != "" {
		date = o
	}

	if o := os.Getenv("COMMITCOUNT"); o != "" {
		count = o
	}

}
