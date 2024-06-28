package main

import (
	"testing"

	"github.com/gookit/color"
	"github.com/vickxxx/my_dbg/models"
)

func TestAddParents(t *testing.T) {
	// addParents(nil, 3)
}

func TestParseLine(t *testing.T) {
	ret := models.ParseLine("T@9:   sql_parse.cc:    3: | | <dispatch_sql_command")
	color.Redln(ret)
}

// ".*@(\d+):\s+(\S+):\s+(\d+):\s+(\d+):[\s\|]*(.+)"
