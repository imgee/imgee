// Package utils provides ...
package utils

import (
	"regexp"
	"strconv"
	"strings"
)

//Flag parses the command arguments syntax:
// -flag=x
// -flag x
// help
func Flag(args string) (string, map[string]interface{}) {
	var (
		r      = make(map[string]interface{}, 10)
		err    error
		target string
	)

	// in case we have args without target
	args = " " + args

	// range
	re := regexp.MustCompile(`(?i)\s-([a-z|0-9]+)[=|\s]{0,1}(\d+\-\d+)`)
	f := re.FindAllStringSubmatch(args, -1)
	for _, kv := range f {
		r[kv[1]] = kv[2]
		args = strings.Replace(args, kv[0], "", 1)
	}
	// none-boolean flags
	for _, rgx := range []string{
		`(?i)\s{1}-([a-z]+)[=|\s](-[0-9]+)`, // negative number
		`(?i)\s{1}-([a-z|0-9]+)[=|\s]([0-9|a-z|'"{}:\.\/_@!#$%^&*)(\+]+)`} {
		re = regexp.MustCompile(rgx)
		f = re.FindAllStringSubmatch(args, -1)
		for _, kv := range f {
			if len(kv) > 1 {
				// trim extra characters (' and ") from value
				kv[2] = strings.Trim(kv[2], "'")
				kv[2] = strings.Trim(kv[2], `"`)
				r[kv[1]], err = strconv.Atoi(kv[2])
				if err != nil {
					r[kv[1]] = kv[2]
				}
				args = strings.Replace(args, kv[0], "", 1)
			}
		}
	}
	// boolean flags
	re = regexp.MustCompile(`(?i)\s-([a-z|0-9]+)`)
	f = re.FindAllStringSubmatch(args, -1)
	for _, kv := range f {
		if len(kv) == 2 {
			r[kv[1]] = ""
			args = strings.Replace(args, kv[0], "", 1)
		}
	}
	// target
	re = regexp.MustCompile(`(?i)^[^-][\S|\w\s]*`)
	t := re.FindStringSubmatch(args)
	if len(t) > 0 {
		target = strings.TrimSpace(t[0])
	}
	// help
	if m, _ := regexp.MatchString(`(?i)help$`, args); m {
		r["help"] = true
	}

	return target, r
}

// SetFlag returns command option(s)
func SetFlag(flag map[string]interface{}, option string, v interface{}) interface{} {
	if sValue, ok := flag[option]; ok {
		switch v.(type) {
		case int:
			if v, ok := sValue.(int); ok {
				return v
			}
		case uint:
			if v, ok := sValue.(uint); ok {
				return v
			}
		case string:
			switch sValue.(type) {
			case string:
				if v, ok := sValue.(string); ok {
					return v
				}
			case int:
				str := strconv.FormatInt(int64(sValue.(int)), 10)
				return str
			case float64:
				str := strconv.FormatFloat(sValue.(float64), 'f', -1, 64)
				return str
			}
		case bool:
			return !v.(bool)
		default:
			return sValue.(string)
		}
	}
	return v
}

func WalkDir(path string) []string {
	return nil
}
