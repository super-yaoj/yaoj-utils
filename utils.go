package utils

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HashValue []byte

func HashSum(a []HashValue) (sum HashValue) {
	h := sha256.New()
	for _, v := range a {
		h.Write(v)
	}
	h.Sum(sum)
	return
}

type ByteValue int64

func (r ByteValue) String() string {
	num := float64(r)
	if num < 1000 {
		return fmt.Sprint(int64(num), "B")
	} else if num < 1e6 {
		return fmt.Sprintf("%.1f%s", num/1e3, "KB")
	} else if num < 1e9 {
		return fmt.Sprintf("%.1f%s", num/1e6, "MB")
	} else {
		return fmt.Sprintf("%.1f%s", num/1e9, "GB")
	}
}

func Map[T any, M any](s []T, f func(T) M) []M {
	var a []M = make([]M, len(s))
	for i, v := range s {
		a[i] = f(v)
	}
	return a
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CopyFile(src, dst string) (int64, error) {
	// log.Printf("CopyFile %s %s", src, dst)
	if src == dst {
		return 0, fmt.Errorf("same path")
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ReaderChecksum(reader io.Reader) Checksum {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return Checksum{}
	}
	var b = hash.Sum(nil)
	if len(b) != 32 {
		panic(b)
	}
	return *(*Checksum)(b)
}

// SHA256 hash for file content.
// for any error, return empty hash
func FileChecksum(name string) Checksum {
	f, err := os.Open(name)
	if err != nil {
		return Checksum{}
	}
	defer f.Close()

	return ReaderChecksum(f)
}

// comparable
type Checksum [32]byte

func (r Checksum) String() string {
	s := ""
	for _, v := range r {
		s += fmt.Sprintf("%02x", v)
	}
	return s
}

type LangTag int

const (
	Lcpp LangTag = iota
	Lcpp11
	Lcpp14
	Lcpp17
	Lcpp20
	Lpython2
	Lpython3
	Lgo
	Ljava
	Lc
	Lplain
	Lpython
)

var LangSuf = []string{
	".cpp98.cpp",
	".cpp11.cpp",
	".cpp14.cpp",
	".cpp17.cpp",
	".cpp20.cpp",
	".py2.py",
	".py3.py",
	".go",
	".java",
	".c",
	".txt",
	".py",
}

// 根据字符串推断程序语言
func SourceLang(s string) LangTag {
	if strings.Contains(s, "java") {
		return Ljava
	}
	if strings.Contains(s, "cpp") || strings.Contains(s, "cc") {
		if strings.Contains(s, fmt.Sprint(11)) {
			return Lcpp11
		}
		if strings.Contains(s, fmt.Sprint(14)) {
			return Lcpp14
		}
		if strings.Contains(s, fmt.Sprint(17)) {
			return Lcpp17
		}
		if strings.Contains(s, fmt.Sprint(20)) {
			return Lcpp20
		}
		return Lcpp
	}
	if strings.Contains(s, "py") {
		if strings.Contains(s, fmt.Sprint(2)) {
			return Lpython2
		}
		if strings.Contains(s, fmt.Sprint(3)) {
			return Lpython3
		}
		return Lpython
	}
	if strings.Contains(s, "go") {
		return Lgo
	}
	if strings.Contains(s, "c") {
		return Lc
	}
	return Lplain
}

type CtntType int

const (
	Cplain CtntType = iota
	Cbinary
	Csource
)

func init() {
	rand.Seed(time.Now().Unix())
}

// dependon: whether i dependon j.
// Complexity: O(n^2)
func TopoSort(size int, dependon func(i, j int) bool) (res []int, err error) {
	indegree := make([]int, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i != j && dependon(i, j) {
				indegree[i]++
			}
		}
	}
	res = make([]int, 0, size)
	err = nil
	for {
		pre := len(res)
		for i := 0; i < size; i++ {
			if indegree[i] == 0 {
				res = append(res, i)
				indegree[i] = -1
			}
		}
		if pre == len(res) {
			break
		}
		for id := pre; id < len(res); id++ {
			i := res[id]
			for j := 0; j < size; j++ {
				if i != j && dependon(j, i) {
					if indegree[j] < 0 {
						panic("topo sort error")
					}
					indegree[j]--
				}
			}
		}
	}
	if len(res) != size {
		err = fmt.Errorf("not a DAG")
	}
	return
}

// index of the first element equaling to v, otherwise return -1
func FindIndex[T comparable](array []T, v T) int {
	for i, item := range array {
		if item == v {
			return i
		}
	}
	return -1
}

func SHA256(str string) string {
	tmp := sha256.New()
	tmp.Write([]byte(str))
	return fmt.Sprintf("%X", tmp.Sum(nil))
}

func Reverse[T any](arr []T) {
	l := len(arr)
	for i := 0; i*2 < l; i++ {
		arr[i], arr[l-i-1] = arr[l-i-1], arr[i]
	}
}

//use binary search
func HasInt(srt []int, val int) bool {
	i := sort.SearchInts(srt, val)
	return i < len(srt) && srt[i] == val
}

func HasElement[T comparable](arr []T, val T) bool {
	return FindIndex(arr, val) >= 0
}

func JoinArray[T any](val []T) string {
	s := strings.Builder{}
	for i, j := range val {
		s.WriteString(fmt.Sprint(j))
		if i+1 < len(val) {
			s.WriteString(",")
		}
	}
	return s.String()
}

func If[T any](a bool, b T, c T) T {
	if a {
		return b
	}
	return c
}

func GetTempDir() (string, error) {
	return os.MkdirTemp(os.TempDir(), "")
}

func UnzipMemory(mem []byte) (map[string][]byte, error) {
	//OpenReader will open the Zip file specified by name and return a ReadCloser.
	reader, err := zip.NewReader(bytes.NewReader(mem), int64(len(mem)))
	if err != nil {
		return nil, err
	}
	ret := make(map[string][]byte)
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		w := bytes.NewBuffer(nil)
		_, err = io.Copy(w, rc)
		if err != nil {
			return nil, err
		}
		ret[file.Name] = w.Bytes()
		rc.Close()
	}
	return ret, nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func AtoiDefault(str string, def int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return val
}

/*
Whether a is start with b
*/
func StartsWith(a, b string) bool {
	if len(a) < len(b) {
		return false
	}
	return a[: len(b)] == b
}

func TimeStamp() int64 {
	return time.Now().UnixMilli()
}

func DeepCopy(dst, src any) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(src)
	if err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dst)
}

type Numbers interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

func Max[T Numbers](a, b T) T {
	if a < b {
		return b
	} else {
		return a
	}
}

func Min[T Numbers](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func DeleteSlice[T any](a []T, id int) {
	a = append(a[:id], a[id+1:]...)
}
/*
resort an sorted array after one entry has modified
*/
func ResortEntry[T any](a []T, f func(int, int) bool, id int) {
	t := a[id]
	for id > 0 && !f(id - 1, id) {
		a[id] = a[id - 1]
		a[id - 1] = t
	}
	for id + 1 < len(a) && !f(id, id + 1) {
		a[id] = a[id + 1]
		a[id + 1] = t
	}
}

/*
every entry of val has an id, we will arrange val such that the result matches the given id array
*/
func ArrangeById[T any](id []int, val []T, getId func(*T)int) []T {
	ret := make([]T, len(id))
	mp := make(map[int]int)
	for i := range val {
		mp[getId(&val[i])] = i
	}
	for i := range id {
		ret[i] = val[mp[id[i]]]
	}
	return ret
}