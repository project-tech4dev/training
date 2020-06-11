package log

import (
	"log"
	"os"
)

func Init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetOutput(os.Stderr)
}
