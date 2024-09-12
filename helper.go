package main

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
		return strconv.FormatFloat(fileSizeFloat/(1024.0*1024.0), 'f', 2, 64) + "MiB"
	default:
		return strconv.FormatFloat(fileSizeFloat/(1024.0*1024.0*1024.0), 'f', 2, 64) + "GiB"
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
