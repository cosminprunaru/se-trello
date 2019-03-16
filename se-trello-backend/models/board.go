package models

import "gopkg.in/mgo.v2/bson"

/*
Board structure defining what a board looks like
*/
type Board struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Lists []BoardList   `bson:"lists" json:"lists"`
}
