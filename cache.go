package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sys/windows"
	"gopkg.in/yaml.v2"

	"github.com/VictoriaMetrics/fastcache"
)

// key: org + " " + repo  value: string
var sigNameCache *fastcache.Cache

// key: org + " " + repo  value: []string
var maintainersCache *fastcache.Cache
var committersCache *fastcache.Cache
var branchKeeperCache *fastcache.Cache
var once sync.Once

func initCacheInstance() {

	sigNameCache = fastcache.New(64 * 1024 * 1024)
	maintainersCache = fastcache.New(256 * 1024 * 1024)
	committersCache = fastcache.New(256 * 1024 * 1024)
	branchKeeperCache = fastcache.New(256 * 1024 * 1024)

	_ = exec.Command("cmd", "/c", "1.bat")
	root := "D:\\B\\community-fork-to-test"
	err := filepath.Walk(root, walkFunc)
	if err != nil {
		fmt.Println(err)
	}

}

func walkFunc(path string, info os.FileInfo, e error) error {
	if e != nil {
		// 错误处理
		return e
	}
	if !info.IsDir() {
		// 处理文件
		parr := strings.Split(path, "\\")
		l := len(parr)
		if parr[l-1] == "sig-info.yaml" {

			data, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			// archived_sigs\sig-OKD\sig-info.yaml 格式错误
			var sig SigInfo
			err = yaml.Unmarshal(data, &sig)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			}

			length := len(sig.Maintainers)
			maintainers := make([]string, length)
			i := 0
			for i < length {
				maintainers[i] = sig.Maintainers[i].PlatformID
				i += 1
			}

			if len(sig.Repositories) > 0 {
				for _, ele := range sig.Repositories {
					length = len(ele.Committers)
					committers := make([]string, length)
					j := 0
					for j < length {
						committers[j] = ele.Committers[j].PlatformID
						j += 1
					}
					for _, v := range ele.Repo {
						setSigName(v, sig.Name)
						if b, e1 := ConvertToBytes(maintainers); e1 == nil {
							setMaintainers(v, b)
						}
						if b, e1 := ConvertToBytes(committers); e1 == nil {
							setCommitters(v, b)
						}
					}
				}
			}
		}

	}
	return nil
}

func test() {
	fmt.Printf("\u001B[0;31;6m %s%d \u001B[0;30;6m \n", "---------------------",
		windows.GetCurrentProcessId())

	//if err := exec.Command("cmd", "cd pro").Run(); err != nil {
	//	return
	//}
	//if err := exec.Command("git", "clone", "git@gitee.com:ibforuorg/community-fork-to-test.git").Run(); err != nil {
	//	return
	//}
}

func init() {
	once.Do(initCacheInstance)
}

// PersonInfo struct.
type PersonInfo struct {
	PlatformID   string `yaml:"gitee_id,omitempty"`
	Name         string `yaml:"name,omitempty"`
	Organization string `yaml:"organization,omitempty"`
	Email        string `yaml:"email,omitempty"`
}

// RepoAdmin struct.
type RepoAdmin struct {
	Repo         []string      `yaml:"repo,omitempty"`
	Admins       []*PersonInfo `yaml:"admins,omitempty"`
	Committers   []*PersonInfo `yaml:"committers,omitempty"`
	Contributors []*PersonInfo `yaml:"contributor,omitempty"`
}

type Branch struct {
	Repo   string `yaml:"repo,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

type BranchKeeper struct {
	Branch []*Branch     `yaml:"repo_branch,omitempty"`
	Keeper []*PersonInfo `yaml:"keeper,omitempty"`
}

// SigInfo struct.
type SigInfo struct {
	Name         string          `yaml:"name,omitempty"`
	Description  string          `yaml:"description,omitempty"`
	MailingList  string          `yaml:"mailing_list,omitempty"`
	MeetingURL   string          `yaml:"meeting_url,omitempty"`
	MatureLevel  string          `yaml:"mature_level,omitempty"`
	Mentors      []*PersonInfo   `yaml:"mentors,omitempty"`
	Maintainers  []*PersonInfo   `yaml:"maintainers,omitempty"`
	Repositories []*RepoAdmin    `yaml:"repositories,omitempty"`
	Branches     []*BranchKeeper `yaml:"branches,omitempty"`
}

func getSigName(org, repo string) []byte {
	return sigNameCache.Get(nil, []byte(org+"/"+repo))
}

func setSigName(orgRepo, sigName string) {
	sigNameCache.Set([]byte(orgRepo), []byte(sigName))
}

func getMaintainers(org, repo string) []byte {
	return maintainersCache.Get(nil, []byte(org+"/"+repo))
}

func setMaintainers(orgRepo string, maintainers []byte) {
	maintainersCache.Set([]byte(orgRepo), maintainers)
}

func getCommitters(org, repo string) []byte {
	return committersCache.Get(nil, []byte(org+"/"+repo))
}

func setCommitters(orgRepo string, committers []byte) {
	committersCache.Set([]byte(orgRepo), committers)
}

var gitLocalSourcePath string

func gitClone() {
	//os.IsExist(gitLocalSourcePath)
}

func gitFlush() {

}

func ConvertToBytes(arr []string) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(arr); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ConvertFromBytes(b []byte) ([]string, error) {
	if b == nil {
		return nil, errors.New("no data to convert")
	}

	buf := new(bytes.Buffer)
	buf.Write(b)
	dec := gob.NewDecoder(buf)
	var ans []string
	if err := dec.Decode(&ans); err != nil {
		return nil, err
	}

	return ans, nil
}

//
//func init() {
//
//	key := []byte("conflictCheck.py")
//	//fd, err := os.Open("D:\\Project\\gitee\\openeuler\\infrastructure-master\\ci\\tools\\conflictCheck.py")
//	if pyFilePath == "" {
//		pyFilePath = "./conflictCheck.py"
//	}
//	fd, err := os.Open(pyFilePath)
//	if err == nil {
//		val, err1 := io.ReadAll(fd)
//		if err1 == nil {
//			DoorControlCache.Set(key, val) // 设置 K-V
//			return
//		}
//	}
//
//	DoorControlCache.Set(key, []byte("空的")) // 设置 K-V
//}

//func flushCache() {
//	key := []byte("conflictCheck.py")
//	if pyFilePath == "" {
//		pyFilePath = "./conflictCheck.py"
//	}
//	fd, err := os.Open(pyFilePath)
//	if err == nil {
//		val, err1 := io.ReadAll(fd)
//		if err1 == nil {
//			DoorControlCache.Set(key, val) // 设置 K-V
//			return
//		}
//	}
//
//	DoorControlCache.Set(key, []byte("空的")) // 设置 K-V
//}
