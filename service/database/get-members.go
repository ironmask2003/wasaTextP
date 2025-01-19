package database

// Query used to take all member from a group
var queryTakeMembers = `SELECT UserId FROM user_group WHERE GroupId = ?`

// Function used to get all members of a group with the user_group table
func (db appdbimpl) GetMembers(groupId int) ([]User, error) {
	// Get the members of a group
	rows, err := db.c.Query(queryTakeMembers, groupId)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	// All members
	var users []User

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var user User
		err = rows.Scan(&user.UserId)
		if err != nil {
			return nil, err
		}

		user, err = db.GetUserById(user.UserId)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
