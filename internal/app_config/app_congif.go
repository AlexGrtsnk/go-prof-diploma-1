package appconfig

import (
	"flag"
)

type Config struct {
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	Home                 string `env:"HOME"`
	RunAddress           string `env:"RUN_ADDRESS"`
}

// smth
func ParseFlags() (a string, b string, f string) {
	var flagRunAddr string
	var fileName string
	var dataBaseAddress string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&fileName, "r", "idk", "Address of acccural system")
	flag.StringVar(&dataBaseAddress, "d", "localhost", "databaseport")
	flag.Parse()
	return flagRunAddr, fileName, dataBaseAddress
}
