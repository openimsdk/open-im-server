package dataver

/*

UserIDs 顺序
前500顺序


1,2,3,4,5,6,7,8,9

1,3,5,7,8,9


1.sdk添加一个表记录 docID(后续换名字), version
2.sdk同步，先计算idHash，api调用参数idHash, docID, version
3.服务器先判断version变更记录，没有直接返回同步成功。
	有变更，先查版本变更记录，在查前500id，变更记录只保留前500id中的
	根据前500id计算idHash，不一致返回会全量id，不反悔删除id
	全量同步有标识，只返回全量id
	变更记录只包含id，不包括详细信息。
4.sdk通过变更记录，同步数据不一致进行重试。
5.先修改db，在自增版本号，外层加事务















*/
