package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzQWERTYUIOPLKJHGFDSAZXCVBNM"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// 生成随机值
func GetRandomString(length int) string {
	b := make([]byte, length)
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// 获取当前路径
func PWD() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return pwd, nil
}

// 获取程序所在路径
func AbsolutePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	pwd, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		//sLog("The eacePath failed: %s\n", err.Error())
		return "", err
	}

	return pwd, nil
}

// 文件是否存在
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// 判断是文件夹还是文件
func IsDir(f string) bool {
	x, _ := os.Stat(f)
	return x.IsDir()
}

// 创建路径
func MakeDir(_dir string) error {
	if IsExist(_dir) == true {
		return nil
	}
	err := os.MkdirAll(_dir, os.ModePerm)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 判断文件夹内是否没有文件
func DirIsNull(dirPath string) bool {
	flist, err := ioutil.ReadDir(dirPath)
	if err != nil || len(flist) == 0 {
		return false
	}
	return true
}

func GetDirName(path string) string {
	var dirs []string
	dirs = strings.Split(path, `\`)
	if len(dirs) == 1 {
		dirs = strings.Split(path, `/`)
	}
	return dirs[len(dirs)-1]
}

// 将 src 文件, 重写到 dst 文件
func Rewrite(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		err = fmt.Errorf("Rewrite.dst: %w", err)
		return err
	}
	defer func() { _ = r.Close() }()

	w, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf("Rewrite.src: %w", err)
		return err
	}
	defer func() { _ = w.Close() }()

	_, err = io.Copy(w, r)
	return err
}

// 复制 src 文件， 到 dst 位置
func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		err = fmt.Errorf("CopyFile.read: %w", err)
		return err
	}
	defer func() { _ = r.Close() }()

	w, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf("CopyFile.write: %w", err)
		return err
	}
	defer func() { _ = w.Close() }()

	_, err = io.Copy(w, r)
	return nil
}

// 将 src 文件，追加到 dst 位置
func AppendFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		err = fmt.Errorf("CopyFile.read: %w", err)
		return err
	}
	defer func() { _ = r.Close() }()

	w, err := os.OpenFile(dst, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("CopyFile.write: %w", err)
		return err
	}
	defer func() { _ = w.Close() }()

	_, err = io.Copy(w, r)
	return nil
}
