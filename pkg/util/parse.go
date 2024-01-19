package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const sep = ","

// ParseBool parses a HTTP request bool param, already handles response on error.
func ParseBool(c *fiber.Ctx, in, name string) (value bool, ok bool) {
	value, err := strconv.ParseBool(strings.TrimSpace(in))
	if err == nil {
		return value, true
	}
	m := fmt.Sprintf(
		"Parameter '%s' contains '%s' which is not a boolean value. Enter 'true' or 'false', or leave empty",
		name, in)
	wrap400(c, m)
	return false, false
}

// ParsePosInt parses a HTTP request positive integer param, already handles response on error.
func ParsePosInt(c *fiber.Ctx, in, name string) (int64, bool) {
	value, err := strconv.ParseInt(strings.TrimSpace(in), 10, 64)
	if err == nil && value > 0 {
		return value, true
	}
	m := fmt.Sprintf(
		"Parameter '%s' contains '%s' which is not a positive integer value. Enter eg. 2, or leave empty for 1",
		name, in)
	wrap400(c, m)
	return 0, false
}

// SplitStr splits a string by the given separator.
func SplitStr(in string) []string {
	if len(in) == 0 {
		return nil
	}
	return strings.Split(in, sep)
}
