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

	if len(item.FileStatus.Extra) > 0 {
		path = fmt.Sprintf("%s → %s", item.Path, item.Extra)
	} else {
		path = item.Path
	}

	if len(item.Accessory) == 0 {
		return path
	}

	return fmt.Sprintf("%s %s", item.Accessory, item.Path)
}

func NewItem(fileStatus git.FileStatus) Item {
	return Item{
		FileStatus: fileStatus,
		Accessory:  fmt.Sprintf("[%s]", string(fileStatus.Code)),
	}
}
