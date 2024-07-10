package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"
)

const CodeHostingPlatformUserID = "gitee_id"

// SigInfo struct.
type SigInfo struct {
	Name         string        `yaml:"name,omitempty"`
	Description  string        `yaml:"description,omitempty"`
	MailingList  string        `yaml:"mailing_list,omitempty"`
	MeetingURL   string        `yaml:"meeting_url,omitempty"`
	MatureLevel  string        `yaml:"mature_level,omitempty"`
	Mentors      []*PersonInfo `yaml:"mentors,omitempty"`
	Maintainers  []*PersonInfo `yaml:"maintainers,omitempty"`
	Repositories []*RepoAdmin  `yaml:"repositories,omitempty"`
}

// RepoAdmin struct.
type RepoAdmin struct {
	Repo         []string      `yaml:"repo,omitempty"`
	Admins       []*PersonInfo `yaml:"admins,omitempty"`
	Committers   []*PersonInfo `yaml:"committers,omitempty"`
	Contributors []*PersonInfo `yaml:"contributor,omitempty"`
}

// PersonInfo struct.
type PersonInfo struct {
	PlatformID   string `yaml:"gitee_id,omitempty"`
	Name         string `yaml:"name,omitempty"`
	Organization string `yaml:"organization,omitempty"`
	Email        string `yaml:"email,omitempty"`
}

func loadCacheFormGitPlatform(url string) error {
	commonUrl := "https://gitee.com/api/v5/repos/"
	urlQuery := "?access_token=&ref=master"

	req, err := http.NewRequest(http.MethodGet, commonUrl+"ibforuorg/community-test/raw/sig/Test/sig-info.yaml"+urlQuery, nil)
	if err != nil {
		fmt.Println("asdsadasdasdasd")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
	res, err := do(req)
	if err != nil {
		return err
	}
	if err == nil {
		val, err1 := io.ReadAll(res.Body)
		if err1 == nil {
			var s SigInfo
			err = yaml.Unmarshal(val, &s)
			if err != nil {
				return err
			}
			buf := new(bytes.Buffer)
			enc := gob.NewEncoder(buf)
			err := enc.Encode(&s)
			if err != nil {
				return err
			}
			DoorControlCache.Set([]byte("ibforuorg/community-test/raw/sig/Test/sig-info.yaml"), buf.Bytes()) // 设置 K-V
			return nil
		}
	}

	DoorControlCache.Set([]byte("ibforuorg/community-test/raw/sig/Test/sig-info.yaml"), []byte("空的")) // 设置 K-V
	return nil
}

func do(req *http.Request) (resp *http.Response, err error) {
	var tmp = http.DefaultClient
	if resp, err = tmp.Do(req); err == nil {
		return
	}

	maxRetries := 4
	backoff := 100 * time.Millisecond

	for retries := 0; retries < maxRetries; retries++ {
		time.Sleep(backoff)
		backoff *= 2

		if resp, err = tmp.Do(req); err == nil {
			break
		}
	}
	return
}
