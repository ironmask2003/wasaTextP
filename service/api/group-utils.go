package api

import (
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
	groupPhoto, err := imageFunctions.ImageToBase64(imageFunctions.SetDefaultPhoto(g.GroupId))
	if err != nil {
		return err
	}
	g.Photo = groupPhoto
	return nil
}
