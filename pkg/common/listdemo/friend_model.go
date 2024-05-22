package listdemo

type friendModel struct {
	db *List[*Friend, *FriendElem]
}
