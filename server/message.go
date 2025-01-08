package main

import (
	"encoding/json"
	"log"
)

type MsgType = int

const (
	ID_MESSAGE MsgType = iota
	GAME_PARAMS_REQ
	GAME_PARAMS_RESP
	MOVE_MESSAGE
	CLIENT_MOVE_MESSAGE
	JOIN_GAME_REQ
	JOIN_GAME_RESP
	OTHER_CLIENT_JOIN

	ID_MESSAGE_STR          string = "id_message"
	GAME_PARAMS_REQ_STR     string = "game_params_req"
	GAME_PARAMS_RESP_STR    string = "game_params_resp"
	MOVE_MESSAGE_STR        string = "move_message"
	CLIENT_MOVE_MESSAGE_STR string = "client_move_message"
	JOIN_GAME_REQ_STR       string = "join_game_req"
	JOIN_GAME_RESP_STR      string = "join_game_resp"
	OTHER_CLIENT_JOIN_STR   string = "other_client_join"
)

type MessageRaw = struct {
	MessageType string          `json:"message_type"`
	Message     json.RawMessage `json:"message"`
}

type IdMessage = struct {
	Id string `json:"id"`
}

type GameParamsReq = struct {
	Id         string `json:"id"`
	IsDiamonds bool   `json:"is_diamonds"`
	NumRows    uint   `json:"num_rows"`
	NumCols    uint   `json:"num_cols"`
	NumColors  uint   `json:"num_colors"`
}

type GameParamsResp = struct {
	Board  *Board `json:"board"`
	GameId uint64 `json:"game_id"`
}

type MoveMessage = struct {
	Board  *Board `json:"board"`
	MyTurn bool   `json:"my_turn"`
}

type ClientMoveMessage = struct {
	Id    string `json:"id"`
	Color uint   `json:"color"`
}

type JoinGameReq = struct {
	Id     string `json:"id"`
	GameId uint64 `json:"game_id"`
}

type JoinGameResp = struct {
	Success   bool   `json:"success"`
	PlayerNum uint   `json:"player_num"` // 0 or 1. determines their board side
	Board     *Board `json:"board"`
	MyTurn    bool   `json:"my_turn"`
}

type OtherClientJoin = struct {
	MyTurn bool `json:"my_turn"`
}

// Returns the message type and the raw message inside
func Preprocess(message []byte) (MsgType, json.RawMessage) {
	var msgStruct MessageRaw
	err := json.Unmarshal(message, &msgStruct)
	if err != nil {
		log.Println("JSON parsing: ", err)
	}

	var msgType MsgType
	switch msgStruct.MessageType {
	case ID_MESSAGE_STR:
		msgType = ID_MESSAGE
	case GAME_PARAMS_REQ_STR:
		msgType = GAME_PARAMS_REQ
	case GAME_PARAMS_RESP_STR:
		msgType = GAME_PARAMS_RESP
	case MOVE_MESSAGE_STR:
		msgType = MOVE_MESSAGE
	case CLIENT_MOVE_MESSAGE_STR:
		msgType = CLIENT_MOVE_MESSAGE
	case JOIN_GAME_REQ_STR:
		msgType = JOIN_GAME_REQ
	case JOIN_GAME_RESP_STR:
		msgType = JOIN_GAME_RESP
	case OTHER_CLIENT_JOIN_STR:
		msgType = OTHER_CLIENT_JOIN
	default:
		log.Println("JSON parsing: ERROR invalid message type")
	}

	return msgType, msgStruct.Message
}

func NewMessageIdMessage(id string) MessageRaw {
	idMessage, err := json.Marshal(IdMessage{
		id,
	})
	if err != nil {
		log.Println("Marshal: ", err)
	}
	message := MessageRaw{
		ID_MESSAGE_STR,
		json.RawMessage(idMessage),
	}
	return message
}

func NewMessageGameParamsResp(board *Board, gameId uint64) MessageRaw {
	gameParamsResp, err := json.Marshal(GameParamsResp{
		board,
		gameId,
	})
	if err != nil {
		log.Println("Marshal: ", err)
	}
	message := MessageRaw{
		GAME_PARAMS_RESP_STR,
		json.RawMessage(gameParamsResp),
	}
	return message
}

func NewMessageMoveMessage(board *Board, myTurn bool) MessageRaw {
	gameParamsResp, err := json.Marshal(MoveMessage{
		board,
		myTurn,
	})
	if err != nil {
		log.Println("Marshal: ", err)
	}
	message := MessageRaw{
		MOVE_MESSAGE_STR,
		json.RawMessage(gameParamsResp),
	}
	return message
}

func NewMessageJoinGameResp(success bool, playerNum uint, board *Board, myTurn bool) MessageRaw {
	joinGameResp, err := json.Marshal(JoinGameResp{
		success,
		playerNum,
		board,
		myTurn,
	})
	if err != nil {
		log.Println("Marshal: ", err)
	}
	message := MessageRaw{
		JOIN_GAME_RESP_STR,
		json.RawMessage(joinGameResp),
	}
	return message
}

func NewMessageOtherClientJoin(myTurn bool) MessageRaw {
	otherClientJoin, err := json.Marshal(OtherClientJoin{
		myTurn,
	})
	if err != nil {
		log.Println("Marshal: ", err)
	}
	message := MessageRaw{
		OTHER_CLIENT_JOIN_STR,
		json.RawMessage(otherClientJoin),
	}
	return message
}
