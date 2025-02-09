package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	list, err := os.ReadDir(".")
	checkError(err)
	files := make([]string, 0, len(list))
	for _, item := range list {
		switch item.Name() {
		case ".idea", "main.go":
		case "build.bat", "build.sh":
		case "sequencer.exe", "sequencer":
		default:
			files = append(files, item.Name())
		}
	}
	Sort(files)
	numPrefix := classifyNumber(len(files))
	for i := 0; i < len(files); i++ {
		ext := filepath.Ext(files[i])
		name := fmt.Sprintf("%0"+strconv.Itoa(numPrefix+1)+"d%s", i+1, ext)
		err = os.Rename(files[i], name)
		checkError(err)
	}
}

func classifyNumber(num int) int {
	if num < 10 {
		return 1
	}
	return int(math.Log10(float64(num))) + 1
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// reference: https://github.com/facette/natsort

type stringSlice []string

func (s stringSlice) Len() int {
	return len(s)
}

func (s stringSlice) Less(a, b int) bool {
	return Compare(s[a], s[b])
}

func (s stringSlice) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

var chunkifyRegexp = regexp.MustCompile(`(\d+|\D+)`)

func chunkify(s string) []string {
	return chunkifyRegexp.FindAllString(s, -1)
}

// Sort sorts a list of strings in a natural order
func Sort(l []string) {
	sort.Sort(stringSlice(l))
}

// Compare returns true if the first string precedes the second one according to natural order
func Compare(a, b string) bool {
	chunksA := chunkify(a)
	chunksB := chunkify(b)

	nChunksA := len(chunksA)
	nChunksB := len(chunksB)

	for i := range chunksA {
		if i >= nChunksB {
			return false
		}

		aInt, aErr := strconv.Atoi(chunksA[i])
		bInt, bErr := strconv.Atoi(chunksB[i])

		// If both chunks are numeric, compare them as integers
		if aErr == nil && bErr == nil {
			if aInt == bInt {
				if i == nChunksA-1 {
					// We reached the last chunk of A, thus B is greater than A
					return true
				} else if i == nChunksB-1 {
					// We reached the last chunk of B, thus A is greater than B
					return false
				}

				continue
			}

			return aInt < bInt
		}

		// So far both strings are equal, continue to next chunk
		if chunksA[i] == chunksB[i] {
			if i == nChunksA-1 {
				// We reached the last chunk of A, thus B is greater than A
				return true
			} else if i == nChunksB-1 {
				// We reached the last chunk of B, thus A is greater than B
				return false
			}

			continue
		}

		return chunksA[i] < chunksB[i]
	}

	return false
}
