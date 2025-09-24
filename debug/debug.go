package debug

import "os"

func Debug() bool {
	return os.Getenv("DEBUG") == "1"
}
