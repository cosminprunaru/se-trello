package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/cosminprunaru/se-trello/models"
	"github.com/cosminprunaru/se-trello/mongodao"
	"github.com/cosminprunaru/se-trello/seconfig"
	"github.com/gorilla/mux"
)

var config = seconfig.Config{}
var dao = mongodao.BoardsDAO{}

/**************************************
Boards methods
**************************************/
func getAllBoardsEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received getAllBoards request...")
	boards, err := dao.FindAll()
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(boards)
}

func getBoardByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["boardId"]

	log.Println("Received getBoardByID " + id + " request...")

	board, err := dao.FindByID(id)

	if err != nil {
		json.NewEncoder(w).Encode(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(board)
}

func createBoardEndPoint(w http.ResponseWriter, r *http.Request) {
	var board models.Board
	log.Println("Received createBoard request...")

	_ = json.NewDecoder(r.Body).Decode(&board)

	board.ID = bson.NewObjectId()

	err := dao.Insert(board)

	if err != nil {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	boards, _ := dao.FindAll()
	json.NewEncoder(w).Encode(boards)
}

func updateBoardEndPoint(w http.ResponseWriter, r *http.Request) {
	var board models.Board
	log.Println("Received updateBoard request...")

	_ = json.NewDecoder(r.Body).Decode(&board)

	err := dao.Update(board)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	boards, _ := dao.FindAll()
	json.NewEncoder(w).Encode(boards)
}

func deleteBoardEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["boardId"]

	// please only add/delete boards with unique names
	log.Println("Received deleteBoard with name: " + name + "...")
	err := dao.DeleteByID(name)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	boards, _ := dao.FindAll()
	json.NewEncoder(w).Encode(boards)
}

/**************************************
Lists methods
**************************************/
func getBoardByID(r *http.Request) models.Board {
	params := mux.Vars(r)
	boardID := params["boardId"]

	board, err := dao.FindByID(boardID)
	if err != nil {
		log.Fatal("Could not find board with ID" + boardID)
		return models.Board{}
	}
	return board
}

func getAllListsEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received getAllLists request...")
	board := getBoardByID(r)
	json.NewEncoder(w).Encode(board.Lists)
}

func getListByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	board := getBoardByID(r)

	params := mux.Vars(r)
	listID := params["listId"]

	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(board.Lists[i].ID, listID) {
			json.NewEncoder(w).Encode(board.Lists[i])
			return
		}
	}
	json.NewEncoder(w).Encode(http.StatusNotFound)
}

func createListEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received createList request...")
	var list models.BoardList
	listID := -1
	board := getBoardByID(r)

	_ = json.NewDecoder(r.Body).Decode(&list)

	for i := 0; i < len(board.Lists); i++ {
		id, _ := strconv.Atoi(board.Lists[i].ID)
		if listID < id {
			listID = id
		}
	}

	list.ID = strconv.Itoa(listID + 1)
	board.Lists = append(board.Lists, list)

	err := dao.Update(board)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(board)
}

func updateListEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received updateList request...")
	board := getBoardByID(r)
	var list models.BoardList

	_ = json.NewDecoder(r.Body).Decode(&list)

	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(board.Lists[i].ID, list.ID) {
			board.Lists[i] = list
			dao.Update(board)
			json.NewEncoder(w).Encode(board)
			return
		}
	}
	json.NewEncoder(w).Encode(board)
}

func deleteListEndPoint(w http.ResponseWriter, r *http.Request) {
	board := getBoardByID(r)

	params := mux.Vars(r)
	id := params["listId"]

	log.Println("Received deleteList request " + id + "...")
	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(board.Lists[i].ID, id) {
			board.Lists = append(board.Lists[:i], board.Lists[i+1:]...)
			dao.Update(board)
			json.NewEncoder(w).Encode(board)
			return
		}
	}
	json.NewEncoder(w).Encode(board)
}

/**************************************
Cards methods
**************************************/
func getListByID(r *http.Request) models.BoardList {
	params := mux.Vars(r)
	boardID := params["boardId"]
	listID := params["listId"]

	board, err := dao.FindByID(boardID)
	if err != nil {
		log.Fatal("Could not find board with ID" + boardID)
		return models.BoardList{}
	}

	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(listID, board.Lists[i].ID) {
			return board.Lists[i]
		}
	}
	return models.BoardList{}
}

func getAllCardsEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received deleteAllCards request...")
	list := getListByID(r)
	json.NewEncoder(w).Encode(list.Cards)
}

func getCardByIDEndPoint(w http.ResponseWriter, r *http.Request) {
	list := getListByID(r)

	params := mux.Vars(r)
	cardID := params["cardId"]

	for i := 0; i < len(list.Cards); i++ {
		if reflect.DeepEqual(list.Cards[i].ID, cardID) {
			json.NewEncoder(w).Encode(list.Cards[i])
			return
		}
	}
	json.NewEncoder(w).Encode(http.StatusNotFound)
}

func createCardEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received createCard request...")
	var card models.Card
	cardID := -1
	board := getBoardByID(r)
	list := getListByID(r)

	_ = json.NewDecoder(r.Body).Decode(&card)

	for i := 0; i < len(list.Cards); i++ {
		listID, _ := strconv.Atoi(list.ID)
		id, _ := strconv.Atoi(board.Lists[listID].Cards[i].ID)
		if cardID < id {
			cardID = id
		}
	}

	card.ID = strconv.Itoa(cardID + 1)
	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(board.Lists[i].ID, list.ID) {
			board.Lists[i].Cards = append(board.Lists[i].Cards, card)
		}
	}

	err := dao.Update(board)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(board)
}

func updateCardEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received updateCard request...")
	board := getBoardByID(r)
	list := getListByID(r)

	var card models.Card

	_ = json.NewDecoder(r.Body).Decode(&card)

	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(board.Lists[i].ID, list.ID) {
			for j := 0; j < len(list.Cards); j++ {
				if reflect.DeepEqual(list.Cards[j].ID, card.ID) {
					board.Lists[i].Cards[j] = card
					dao.Update(board)
					json.NewEncoder(w).Encode(board)
					return
				}
			}
		}
	}
	json.NewEncoder(w).Encode(board)
}

func deleteCardEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Received deleteCard request...")
	board := getBoardByID(r)
	list := getListByID(r)

	params := mux.Vars(r)
	cardID := params["cardId"]

	for i := 0; i < len(board.Lists); i++ {
		if reflect.DeepEqual(board.Lists[i].ID, list.ID) {
			for j := 0; j < len(board.Lists[i].Cards); j++ {
				if reflect.DeepEqual(board.Lists[i].Cards[j].ID, cardID) {
					board.Lists[i].Cards = append(board.Lists[i].Cards[:j], board.Lists[i].Cards[j+1:]...)
					dao.Update(board)
					json.NewEncoder(w).Encode(board)
					return
				}
			}
		}
	}
	json.NewEncoder(w).Encode(board)
}

func getOptions(w http.ResponseWriter, r *http.Request) {
	log.Println("Received getOptions for CORS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST GET DELETE PUT OPTIONS")

	json.NewEncoder(w).Encode(http.StatusOK)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	r := mux.NewRouter()

	log.Println("Starting REST API on localhost:3000")

	r.HandleFunc("/boards", getOptions).Methods("OPTIONS")
	r.HandleFunc("/boards/{boardId}", getOptions).Methods("OPTIONS")
	r.HandleFunc("/boards/{boardId}/lists", getOptions).Methods("OPTIONS")
	r.HandleFunc("/boards/{boardId}/lists/{listId}", getOptions).Methods("OPTIONS")
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards", getOptions).Methods("OPTIONS")
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards/{cardId}", getOptions).Methods("OPTIONS")

	// Handlers for boards
	r.HandleFunc("/boards", getAllBoardsEndPoint).Methods("GET")
	r.HandleFunc("/boards", createBoardEndPoint).Methods("POST")
	r.HandleFunc("/boards", updateBoardEndPoint).Methods("PUT")
	r.HandleFunc("/boards/{boardId}", deleteBoardEndPoint).Methods("DELETE")
	r.HandleFunc("/boards/{boardId}", getBoardByIDEndPoint).Methods("GET")

	// Handlers for lists
	r.HandleFunc("/boards/{boardId}/lists", getAllListsEndPoint).Methods("GET")
	r.HandleFunc("/boards/{boardId}/lists", createListEndPoint).Methods("POST")
	r.HandleFunc("/boards/{boardId}/lists", updateListEndPoint).Methods("PUT")
	r.HandleFunc("/boards/{boardId}/lists/{listId}", deleteListEndPoint).Methods("DELETE")
	r.HandleFunc("/boards/{boardId}/lists/{listId}", getListByIDEndPoint).Methods("GET")

	// Handlers for cards
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards", getAllCardsEndPoint).Methods("GET")
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards", createCardEndPoint).Methods("POST")
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards", updateCardEndPoint).Methods("PUT")
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards/{cardId}", deleteCardEndPoint).Methods("DELETE")
	r.HandleFunc("/boards/{boardId}/lists/{listId}/cards/{cardId}", getCardByIDEndPoint).Methods("GET")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
