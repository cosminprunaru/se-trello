package mongodao

import (
	"log"

	"github.com/cosminprunaru/se-trello/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BoardsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	BOARDS = "boards"
	LISTS  = "lists"
	CARDS  = "cards"
)

//Connect -> Establish connection to db
func (m *BoardsDAO) Connect() {
	log.Println("Trying to connect to Mongo DB ...")
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal("Mongo error: " + err.Error())
	}
	session.SetSafe(&mgo.Safe{})
	db = session.DB(m.Database)
}

//FindAll -> Find list of boards
func (m *BoardsDAO) FindAll() ([]models.Board, error) {
	var boards []models.Board
	err := db.C(BOARDS).Find(bson.M{}).All(&boards)
	return boards, err
}

//FindByID -> Find a board by its id
func (m *BoardsDAO) FindByID(id string) (models.Board, error) {
	var board models.Board
	err := db.C(BOARDS).FindId(bson.ObjectIdHex(id)).One(&board)
	return board, err
}

//Insert -> Insert a board into database
func (m *BoardsDAO) Insert(board models.Board) error {
	err := db.C(BOARDS).Insert(&board)
	return err
}

//DeleteByID -> Delete an existing board
func (m *BoardsDAO) DeleteByID(name string) error {
	collection := db.C(BOARDS)
	err := collection.Remove(bson.M{"name": name})
	return err
}

//Update -> Update an existing board
func (m *BoardsDAO) Update(board models.Board) error {
	err := db.C(BOARDS).UpdateId(board.ID, &board)
	return err
}
