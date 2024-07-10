package main

import (
	"io"
	"os"

	"github.com/VictoriaMetrics/fastcache"
)

var pyFilePath = ""

// DoorControlCache 初始化一个大小为 32MB 的缓存
var DoorControlCache = fastcache.New(32 * 1024 * 1024)

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

func flushCache() {
	key := []byte("conflictCheck.py")
	if pyFilePath == "" {
		pyFilePath = "./conflictCheck.py"
	}
	fd, err := os.Open(pyFilePath)
	if err == nil {
		val, err1 := io.ReadAll(fd)
		if err1 == nil {
			DoorControlCache.Set(key, val) // 设置 K-V
			return
		}
	}

	DoorControlCache.Set(key, []byte("空的")) // 设置 K-V
}
