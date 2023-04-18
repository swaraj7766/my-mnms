package mnms

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const testdir = "testog"

var testlog = path.Join(testdir, "test.log")

const testmagebyte = 100
const backupnumber = 3
const persize = 200

var testtimeout = time.Minute * 1
var maxsize = testmagebyte * 1024 * 1024 * (backupnumber + 1)

func TestStress_Logrotate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testtimeout)
	defer cancel()
	rand.Seed(time.Now().UnixNano())
	l := getTestLogger()
	defer os.RemoveAll(testdir)
	for {
		select {
		case <-ctx.Done():
			err := l.Close()
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(10 * time.Second)
			err = checkFiles(backupnumber)
			if err != nil {
				t.Fatal(err)
			}
			err = CheckSize(maxsize, testdir)
			if err != nil {
				t.Fatal(err)
			}
			err = CheckGzFileNumbr(backupnumber, testdir)
			if err != nil {
				t.Fatal(err)
			}
			return
		default:
			size := rand.Intn(persize)
			_, err := l.Write([]byte(RandStringRunes(size)))
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}

func getTestLogger() *lumberjack.Logger {
	Logger := &lumberjack.Logger{
		Filename:   testlog,
		MaxSize:    testmagebyte,
		MaxBackups: backupnumber,
		Compress:   true,
		LocalTime:  true,
	}
	return Logger
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func checkFiles(number int) error {

	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		log.Fatal(err)
	}
	number++
	if len(files) != number {
		return fmt.Errorf("files number error,expect:%v,actual:%v", number, len(files))
	}
	return nil
}
func CheckSize(max int, dir string) error {
	n, err := DirSize(dir)
	if err != nil {
		return err
	}
	if n > int64(max) {
		return fmt.Errorf("max size should be %v,actual:%v", maxsize, n)
	}
	return nil
}
func CheckGzFileNumbr(expectnmuber int, dir string) error {
	count := 0
	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), "gz") {
			count++
		}
		return err
	})
	if count != expectnmuber {
		return fmt.Errorf("GzFile,expect:%v,actual:%v", expectnmuber, count)
	}

	return err
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
