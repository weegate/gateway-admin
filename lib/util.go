//@author wuyong
//@date   2018/1/9
//@desc

package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// filepath.Glob/Walk like this
func GetFileBydir(dir string, files []string, dirs []string) []string {
	fileDirs, _ := ioutil.ReadDir(dir)
	for _, file := range fileDirs {
		if file.IsDir() {
			dirName := file.Name()
			if dirName == "." || dirName == ".." {
				continue
			}
			dirs = append(dirs, dirName)
			GetFileBydir(dir+"/"+file.Name(), files, dirs)
		} else {
			fileName := dir + "/" + file.Name()
			files = append(files, fileName)
		}
	}
	return dirs
}

// check dir is exist
func IsDirExist(path string) bool {
	//dirPath, _ := filepath.Abs(path)
	if fileInfo, err := os.Stat(path); err == nil {
		return fileInfo.IsDir()
	} else {
		return os.IsExist(err)
	}
	panic("not reached")
}

func SimpleHttpRequest(method string, url string, readData string, ua string) ([]byte, error) {
	logrus.Info(fmt.Sprintf("method: %s => url: %s => requestData: %s", method, url, readData))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(readData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	respon, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer respon.Body.Close()
	body, err := ioutil.ReadAll(respon.Body)
	if err != nil {
		return nil, err
	}
	logrus.Info(fmt.Sprintf("responseData: %s", string(body)))

	return body, nil
}
