package database

// Query for update the username of the user
var queryUpdateUsername = "UPDATE user SET Username = ? WHERE UserId = ?;"

func (db *appdbimpl) SetMyUsername(UserId int, newUsername string) error {
	// Update the username of the user
	_, err := db.c.Exec(queryUpdateUsername, newUsername, UserId)
	if err != nil {
		return err
	}
	return nil
}
