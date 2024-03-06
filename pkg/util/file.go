package util

import (
	"fmt"
	"os"

	"k8s.io/klog/v2"
)

func init() {
	_pullDir := fmt.Sprintf("./logs")
	if exist, _ := pathExists(_pullDir); !exist {
		err := os.Mkdir(_pullDir, os.ModePerm)
		if err != nil {
			klog.Errorf("mkdir failed!", err)
		} else {
			klog.Infof("mkdir success!")
		}
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
