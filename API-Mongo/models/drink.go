package models

import "gopkg.in/mgo.v2/bson"

type Drink struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Price       int           `bson:"price" json:"price"`
        Size        string        `bson:"size" json:"size"`
	Description string        `bson:"description" json:"description"`
}
