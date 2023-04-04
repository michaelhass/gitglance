package file

import (
	"fmt"

	"github.com/michaelhass/gitglance/internal/git"
)

type ListItem struct {
	FileStatus git.FileStatus
	Path       string
	Accessory  string
}

func (item ListItem) String() string {
	if len(item.Accessory) == 0 {
		return item.Path
	}
	return fmt.Sprintf("%s %s", item.Accessory, item.Path)
}

func NewListItem(fileStatus git.FileStatus) ListItem {
	var (
		path, accessory string
	)

	path = fileStatus.Path
	if len(fileStatus.Extra) > 0 {
		path = fmt.Sprintf("%s → %s", path, fileStatus.Extra)
	}

	accessory = fmt.Sprintf("[%s]", string(fileStatus.Code))

	return ListItem{
		FileStatus: fileStatus,
		Path:       path,
		Accessory:  accessory,
	}
}
