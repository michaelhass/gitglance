package filelist

import (
	"fmt"

	"github.com/michaelhass/gitglance/internal/git"
)

type Item struct {
	FileStatus git.FileStatus
	Path       string
	Accessory  string
}

func (item Item) String() string {
	if len(item.Accessory) == 0 {
		return item.Path
	}
	return fmt.Sprintf("%s %s", item.Accessory, item.Path)
}

func NewItem(fileStatus git.FileStatus) Item {
	var (
		path, accessory string
	)

	path = fileStatus.Path
	if len(fileStatus.Extra) > 0 {
		path = fmt.Sprintf("%s → %s", path, fileStatus.Extra)
	}

	accessory = fmt.Sprintf("[%s]", string(fileStatus.Code))

	return Item{
		FileStatus: fileStatus,
		Path:       path,
		Accessory:  accessory,
	}
}
