package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ohko/hst"
)

const filechunk = 1024 * 4

var filePaths []string // 存储所有需要分析的文件列表
var hashs map[string][]string
var repeats map[string][]string
var ww *hst.HST
var doing int64

func main() {
	// log.SetFlags(log.Flags() | log.Lshortfile)
	readAllFiles()

	ww = hst.NewHST(nil)
	ww.HandleFunc("/", wwRoot)
	ww.HandleFunc("/refresh", wwRefresh)
	ww.HandleFunc("/files", wwFiles)
	ww.HandleFunc("/remove", wwRemove)
	ww.ListenHTTP(":8080")
}
func wwRoot(c *hst.Context) {
	http.FileServer(assetFS()).ServeHTTP(c.W, c.R)
}
func wwRefresh(c *hst.Context) {
	go readAllFiles()
	c.JSON2(0, "ok")
}
func wwFiles(c *hst.Context) {
	if len(repeats) == 0 {
		c.JSON2(0, "无重复文件发现")
	} else {
		c.JSON2(0, repeats)
	}
}

func wwRemove(c *hst.Context) {
	if !atomic.CompareAndSwapInt64(&doing, 0, 1) {
		c.JSON2(1, "其它任务执行中...")
		return
	}
	c.R.ParseForm()
	fs := c.R.Form["fs[]"]
	go func() {
		count := len(fs)
		for k, v := range fs {
			log.Printf("[%d/%d]移除文件: %s\n", k+1, count, v)
			if err := os.Remove(v); err != nil {
				log.Println("移除错误：", err)
			}
		}
		log.Println("移除完成")
		repeats = map[string][]string{}
		doing = 0
	}()
	c.JSON2(0, "OK")
}

func readAllFiles() {
	if !atomic.CompareAndSwapInt64(&doing, 0, 1) {
		return
	}
	filePaths = []string{}
	hashs = make(map[string][]string)
	repeats = map[string][]string{}
	log.Println("开始分析文件...")

	hashAllFiles("./")
	for k, v := range filePaths {
		log.Printf("[%d/%d] 分析: %s\n", k+1, len(filePaths), v)
		hash := getHash(v)
		hashs[hash] = append(hashs[hash], v)
	}

	for k, v := range hashs {
		if len(v) > 1 {
			repeats[k] = append(repeats[k], v...)
		}
	}
	log.Println("分析完成，可刷新浏览器了，文件总数:", len(filePaths), "个 / 重复文件：", len(repeats), "组")
	doing = 0
}
func hashAllFiles(path string) {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range fi {
		if v.IsDir() {
			hashAllFiles(path + v.Name())
		} else {
			filePath := path + v.Name()
			filePaths = append(filePaths, filePath)
		}
	}
}

func getHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
