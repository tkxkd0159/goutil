package io

import (
	"strings"
)

func CleanString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
