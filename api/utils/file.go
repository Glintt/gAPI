package utils

import (
	"io/ioutil"
)

func LoadJsonFile(location string) ([]byte, error) {
	return ioutil.ReadFile(location)
}
