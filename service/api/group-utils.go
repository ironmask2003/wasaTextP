package api

import (
	"regexp"
	"wasa.project/service/api/imageFunctions"
	"wasa.project/service/database"
)

type Group struct {
	GroupId   int    `json:"groupId"`
	GroupName string `json:"groupName"`
	Photo     string `json:"photo"`
}

func (g *Group) ConvertGroupForDB() database.Group {
	return database.Group{
		GroupId:   g.GroupId,
		GroupName: g.GroupName,
	}
}

func (g *Group) ConvertGroupFromDB(groupDB database.Group) error {
	g.GroupId = groupDB.GroupId
	g.GroupName = groupDB.GroupName
	groupPhoto, err := imageFunctions.ImageToBase64(imageFunctions.SetDefaultPhotoGroup(g.GroupId))
	if err != nil {
		return err
	}
	g.Photo = groupPhoto
	return nil
}

// Check if the group respect the regex
func (g *Group) IsValid() bool {
	validUser := regexp.MustCompile(`^[a-z][a-z0-9]{2,13}$`)
	return validUser.MatchString(g.GroupName)
}
