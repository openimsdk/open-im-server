package utilsv2

import "Open_IM/pkg/common/db/table"

func demo() {

	groups := []*table.GroupModel{}

	groups = DuplicateRemovalAny(groups, func(t *table.GroupModel) string {
		return t.GroupID
	})

}
