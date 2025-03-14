package cache_service

import (
	"github.com/youngking/gin-blog/pkg/e"
	"strconv"
	"strings"
)

type Tags struct {
	ID   int
	Name string

	State    int
	PageNum  int
	PageSize int
}

func (t *Tags) GetTagsKey() string {
	keys := []string{e.CACHE_TAG, "LIST"}
	if t.ID > 0 {
		keys = append(keys, strconv.Itoa(t.ID))
	}
	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}
	return strings.Join(keys, "_")
}
