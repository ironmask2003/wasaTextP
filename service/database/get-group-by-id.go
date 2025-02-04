package database

var queryFindGroupById = `SELECT GroupId, GroupName FROM group_t WHERE GroupId = ?;`

func (db *appdbimpl) GetGroupById(groupId int) (Group, error) {
	var group Group
	err := db.c.QueryRow(queryFindGroupById, groupId).Scan(&group.GroupId, &group.GroupName)
	return group, err
}
