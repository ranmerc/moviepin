package utils

import (
	"log"
	"os"
)

var Logger *log.Logger
var ErrorLogger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "[ERROR]:", log.Ldate|log.Ltime|log.Lshortfile)
}
