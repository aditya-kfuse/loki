package pattern

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type node interface {
	fmt.Stringer
}

type Expr []node

func (e Expr) hasCapture() bool {
	return e.captureCount() != 0
}

func (e Expr) validate() error {
	if !e.hasCapture() {
		return ErrNoCapture
	}
	// Consecutive captures are not allowed.
	for i, n := range e {
		if i+1 >= len(e) {
			break
		}
		if _, ok := n.(capture); ok {
			if _, ok := e[i+1].(capture); ok {
				return fmt.Errorf("found consecutive capture '%s': %w", n.String()+e[i+1].String(), ErrInvalidExpr)
			}
		}
	}

	caps := e.captures()
	uniq := map[string]struct{}{}
	for _, c := range caps {
		if _, ok := uniq[c]; ok {
			return fmt.Errorf("duplicate capture name (%s): %w", c, ErrInvalidExpr)
		}
		uniq[c] = struct{}{}
	}
	return nil
}

func (e Expr) captures() (captures []string) {
	for _, n := range e {
		if c, ok := n.(capture); ok && !c.isUnamed() {
			captures = append(captures, c.Name())
		}
	}
	return
}

func (e Expr) escapeString(i string) string {
	// note: \ should be at the beginning of this string, so that we don't escape escaped characters again.
	s := `\[](){}.+*?|^$`
	for _, c := range s {
		if strings.Contains(i, string(c)) {
			i = strings.ReplaceAll(i, string(c), `\`+string(c))
		}
	}
	return i
}

func (e Expr) CapturesMapWithRegex() (string, map[string]int) {
	var sb strings.Builder
	m := map[string]int{}
	if len(e) == 0 {
		return "", m
	}
	gid := 1
	lnCapture := false
	lUnamed := false
	sb.WriteString("^")
	for i, n := range e {
		if c, ok := n.(capture); ok && !c.isUnamed() {
			m[c.Name()] = gid
			sb.WriteString("(.+?)")
			gid += 1
			if i == len(e)-1 {
				lnCapture = true
			}
			lUnamed = false
		} else if c.isUnamed() {
			if lUnamed {
				continue
			}
			sb.WriteString(".*?")
			lUnamed = true
		}
		if l, ok := n.(literals); ok {
			sb.WriteString(e.escapeString(string(l)))
			lUnamed = false
		}
	}

	if lnCapture {
		sb.WriteString("$")
	} else {
		sb.WriteString(".*?$")
	}

	return sb.String(), m
}

func (e Expr) captureCount() (count int) {
	return len(e.captures())
}

type capture string

func (c capture) String() string {
	return "<" + string(c) + ">"
}

func (c capture) Name() string {
	return string(c)
}

func (c capture) isUnamed() bool {
	return string(c) == underscore
}

type literals []byte

func (l literals) String() string {
	return string(l)
}

func runesToLiterals(rs []rune) literals {
	res := make([]byte, len(rs)*utf8.UTFMax)
	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(res[count:], r)
	}
	res = res[:count]
	return res
}
