package filelist

import (
	"fmt"

	"github.com/michaelhass/gitglance/internal/git"
)

type Item struct {
	git.FileStatus
	Accessory string
}

func (item Item) String() string {
	var path string

	if item.IsRenamed() {
		path = fmt.Sprintf("%s → %s", item.Extra, item.Path)
	} else {
		path = item.Path
	}

	if len(item.Accessory) == 0 {
		return path
	}

	return fmt.Sprintf("%s %s", item.Accessory, path)
}
func NewItem(fileStatus git.FileStatus, accessory string) Item {
	return Item{
		FileStatus: fileStatus,
		Accessory:  fmt.Sprintf("[%s]", accessory),
	}
}
