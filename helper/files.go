package helper

import (
	"path"
	"strconv"
	"strings"
)

func HumanSize(fileSize uint64) string {
	fileSizeFloat := float64(fileSize)
	switch {
	case fileSize < 1024.0:
		return strconv.FormatFloat(fileSizeFloat, 'f', 2, 64) + "B"
	case fileSize < 1024.0*1024.0:
		return strconv.FormatFloat(fileSizeFloat/1024.0, 'f', 2, 64) + "KiB"
	case fileSize < 1024.0*1024.0*1024.0:
		return strconv.FormatFloat(fileSizeFloat/float64(1024*1024), 'f', 2, 64) + "MiB"
	default:
		return strconv.FormatFloat(fileSizeFloat/float64(1024*1024*1024), 'f', 2, 64) + "GiB"
	}

}

func AbridgedFileName(fileName string) string {
	const maxFileNameLen = 72
	const numEOFCharsToShow = 8

	switch len(fileName) < maxFileNameLen-numEOFCharsToShow {
	case true:
		return fileName
	default:
		fileExt := path.Ext(fileName)
		fileParts := strings.SplitAfter(fileName, fileExt)
		fileName := fileParts[0]

		return fileName[:maxFileNameLen-numEOFCharsToShow] + "..." + fileName[len(fileName)-numEOFCharsToShow:]
	}
}

// This isn't too hard. We check that our current directory (dir1) is < dir2's length by 1, then we just check that each item in dir1 matches each item in dir2 (to make sure it's actually a subdirectory)
func IsSubdirectory(dir1 []string, dir2 []string) bool {
	if len(dir1)+1 != len(dir2) {
		return false
	}

	for idx := range dir1 {
		if dir1[idx] != dir2[idx] {
			return false
		}
	}

	return true
}

func SliceToDirectory(dir []string) string {
	return "/" + strings.Join(dir, "/")
}

func DirectoryToSlice(dir string) []string {
	if dir == "/" {
		return []string{}
	}
	dir = strings.TrimPrefix(dir, "/")
	slice := strings.Split(dir, "/")
	for idx := range slice {
		if slice[idx] == " " {
			slice[idx] = ""
		}
	}
	return slice
}

func IsDirEqual(dir1 []string, dir2 []string) bool {
	if len(dir1) != len(dir2) {
		return false
	}

	for idx := range len(dir1) {
		if dir1[idx] != dir2[idx] {
			return false
		}
	}

	return true
}
