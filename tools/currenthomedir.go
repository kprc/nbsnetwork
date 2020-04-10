package tools

import (
	"runtime"
	"os"
	"bytes"
	"os/exec"
	"strings"
	"os/user"
	"errors"
	"log"
	"github.com/kprc/nbsdht/nbserr"
	"io/ioutil"
)

var filenotfind = nbserr.NbsErr{ErrId:nbserr.FILE_NOT_FOUND,Errmsg:"File Not Found"}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}


func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
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

func Save2File(data []byte,filename string) error {

	f,err:=os.OpenFile(filename,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,0755)
	if err!=nil{
		log.Fatal(err)
	}
	defer f.Close()

	if _,err:=f.Write(data);err!=nil{
		f.Close()
		log.Fatal(err)
	}

	return nil
}

func OpenAndReadAll(filename string) (data []byte,err error)  {
	if !FileExists(filename){
		return nil,filenotfind
	}

	f,err:=os.OpenFile(filename,os.O_RDONLY,0755)
	if err!=nil{
		log.Fatal(err)
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}