package apputil

import "os"

func Version() string {
	return os.Getenv("BUILD_VERSION")
}
