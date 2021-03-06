package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if runtime.GOOS == "windows" {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func FileList(dir, suffix string) (filelist []string, err error) {
	// 获取该路径下的文件列表，并不会像Walk一样遍历子目录
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if fi.IsDir() {
			continue
		} else {
			if strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
				filelist = append(filelist, fi.Name())
			}
		}
	}
	return
}

func LockPid(pid int, dir string) error {
	pids, err := FileList(dir, "lock")
	if err != nil {
		return err
	}

	if len(pids) > 0 {
		return fmt.Errorf("process(%s) existed", strings.Split(pids[0], ".")[0])
	}

	lock, err := os.Create(filepath.Join(dir, fmt.Sprintf("%d.lock", pid)))
	if err != nil {
		return err
	}
	lock.Close()
	return nil
}

func UnLockPid(pid int, dir string) error {
	return os.Remove(filepath.Join(dir, fmt.Sprintf("%d.lock", pid)))
}

func LocalEnv() error {
	workDir, e := os.Getwd()
	if e != nil {
		return e
	}
	if err := godotenv.Load(filepath.Join(workDir, ".env")); err != nil {
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
