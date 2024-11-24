package functions

import "fmt"

func SetDefaultPhoto(userId int) string {
	return fmt.Sprintf("./storage/%d/default_profile_photo.png", userId)
}
