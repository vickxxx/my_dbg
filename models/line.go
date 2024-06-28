package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gookit/color"
)

// T@9:   viosocket.cc:   140:    3: | | >vio_read

type Line struct {
	SessionID int64  `json:"-"`
	FileName  string `json:"-"`
	LineNO    int64  `json:"-"`
	Loc       string `json:"loc,omitempty"`
	Depth     int64  `json:"-"`
	Type      int64  `json:"-"` // 0: raw 1: in 2: out, 3: data
	Title     string
	Raw       string  `json:"-"`
	Detail    []*Line `json:"detail,omitempty"`
}

func ParseLine(lineStr string) *Line {
	color.Redln(lineStr)
	re := regexp.MustCompile(`.*@(\d+):\s+(\S+):\s+(\d+):\s+(\d+):[\s\|]*(.+)`)
	re2 := regexp.MustCompile(`.*@(\d+):\s+(\S+):\s+(\d+):[\s\|]*(.+)`)

	matches := re.FindStringSubmatch(lineStr)
	line := Line{
		Raw: lineStr,
	}
	if matches != nil {
		line.SessionID, _ = strconv.ParseInt(matches[1], 10, 64)
		line.FileName = matches[2]
		line.LineNO, _ = strconv.ParseInt(matches[3], 10, 64)
		line.Loc = fmt.Sprintf("%s:%d", line.FileName, line.LineNO)
		line.Depth, _ = strconv.ParseInt(matches[4], 10, 64)
		line.Title = matches[5]
		if strings.HasPrefix(line.Title, ">") {
			line.Type = 1
			line.Title = strings.TrimLeft(line.Title, ">")
		} else if strings.HasPrefix(line.Title, "<") {
			line.Type = 2
			line.Title = strings.TrimLeft(line.Title, "<")
		} else {
			line.Type = 3
		}
		color.Yellowln(spew.Sdump(line))
		color.Yellowln("-------------------")
	}

	if matches == nil {
		matches = re2.FindStringSubmatch(lineStr)
		if matches != nil {
			line.SessionID, _ = strconv.ParseInt(matches[1], 10, 64)
			line.FileName = matches[2]
			line.Loc = line.FileName
			// line.LineNO, _ = strconv.ParseInt(matches[3], 10, 64)
			line.Depth, _ = strconv.ParseInt(matches[3], 10, 64)
			line.Title = matches[4]
			if strings.HasPrefix(line.Title, ">") {
				line.Type = 1
				line.Title = strings.TrimLeft(line.Title, ">")
			} else if strings.HasPrefix(line.Title, "<") {
				line.Type = 2
				line.Title = strings.TrimLeft(line.Title, "<")
			} else {
				line.Type = 3
			}
			color.Yellowln(spew.Sdump(line))
			color.Yellowln("-------------------")
		}
	}

	return &line
}
