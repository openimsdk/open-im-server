package im_mysql_model

func (tb *Group) Create(groups []*Group) error {
	return nil
}
func (tb *Group) Take(groupIDs []string) (*Group, error) {
	return nil, nil
}
func (tb *Group) Get(groupIDs []string) (*Group, error) {
	return nil, nil
}
func (tb *Group) Update(groups []*Group) error {
	return nil
}
func (tb *Group) GetByName(groupName string, pageNumber, showNumber int32) ([]GroupWithNum, int64, error) {

}
func (tb *Group) GetGroups(pageNumber, showNumber int) ([]GroupWithNum, error) {
}
func (tb *Group) OperateGroupStatus(groupId string, groupStatus int32) error {
}

func (tb *Group) GetCountsNum(groupIDs []string) ([]int32, error) {

}

func (tb *Group) UpdateDefaultZero(groupID string, args map[string]interface{}) error {
}

func (tb *Group) GetGroupIDsByGroupType(groupType int) ([]string, error) {

}
