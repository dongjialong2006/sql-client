package file

import (
	"fmt"
	"os"
	"strings"
)

func WriteFile(path string, data []byte) error {
	var dir string = ""
	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		dir = "./" + path
	} else {
		dir = path[:pos+1]
	}

	_, err := os.Stat(dir)
	if nil != err {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, os.ModePerm)
			if nil != err {
				return err
			}
		} else {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if nil != err {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func CreatePath(path string) error {
	if "" == path {
		return fmt.Errorf("path is empty.")
	}

	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		return nil
	}
	path = path[:pos]
	_, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
		}
	}

	return err
}
