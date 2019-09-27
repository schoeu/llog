package agent

import (
	"github.com/schoeu/nma/util"
	"regexp"
)

func IsInclude(text string, regs []string) bool {
	for _, v := range regs {
		r, err := regexp.Compile(v)
		util.ErrHandler(err)
		if r.MatchString(text) {
			return true
		}
	}
	return false
}
