package main

import (
	"encoding/json"
	"log"
)

func (client *Client) handleBinaryMessage(message []byte) {
	log.Println("WARNING: Received binary message (unsupported)")
}

func (client *Client) handleTextMessage(message []byte) {
	// Parse JSON
	msgType, msgJSON := Preprocess(message)

	switch msgType {
	case ID_MESSAGE:
		var idMessage IdMessage
		err := json.Unmarshal(msgJSON, &idMessage)
		if err != nil {
			log.Println("Unmarshal: ", err)
		}

		client.handleIdMessage(idMessage)

	case GAME_PARAMS_REQ:
		var gameParamsReq GameParamsReq
		err := json.Unmarshal(msgJSON, &gameParamsReq)
		if err != nil {
			log.Println("Unmarshal: ", err)
		}

		client.handleGameParamsReq(gameParamsReq)

	case GAME_PARAMS_RESP:
		// This message type should never be sent by a client
		log.Println("Received GAME_PARAMS_RESP message from client. Closing connection")
		client.Close()
	case MOVE_MESSAGE:
		// This message type should never be sent by a client
		log.Println("Received MOVE_MESSAGE message from client. Closing connection")
		client.Close()
	case CLIENT_MOVE_MESSAGE:
		var clientMoveMessage ClientMoveMessage
		err := json.Unmarshal(msgJSON, &clientMoveMessage)
		if err != nil {
			log.Println("Unmarshal: ", err)
		}

		client.handleClientMoveMessage(clientMoveMessage)
	case JOIN_GAME_REQ:
		var joinGameReq JoinGameReq
		err := json.Unmarshal(msgJSON, &joinGameReq)
		if err != nil {
			log.Println("Unmarshal: ", err)
		}

		client.handleJoinGameReq(joinGameReq)

	default:
		if client != nil {
			log.Println("Received invalid message type from client ", client.id)
		} else {
			log.Println("Received invalid message type")
		}
		client.Close()
		return
	}
}

func (client *Client) handleIdMessage(idMessage IdMessage) {
	NewPlayer(client, idMessage.Id)

	// Respond with client ID
	idResponse := NewMessageIdMessage(client.id)
	client.write(idResponse)

	// If the client is in a game, we send them a JoinGameResp
	// or a GameParamsResp, depending on whether they have an opponent yet
	// TODO we must check if there is a
	if client.game != 0 {
		board, player1, player2, turnPlayerId := GetBoard(client.game)
		myTurn := turnPlayerId == client.id

		var playerNum uint
		var opponent string
		if client.id == player1 {
			playerNum = 0
			opponent = player2
		} else if client.id == player2 {
			playerNum = 1
			opponent = player1
		} else {
			log.Fatal("Got board for client that isn't in game")
		}

		// If there is no opponent yet, we send a GameParamsResp
		if opponent == "" {
			inGameMessage := NewMessageGameParamsResp(board, client.game)
			client.write(inGameMessage)
		} else {
			inGameMessage := NewMessageJoinGameResp(true, playerNum, board, myTurn)
			client.write(inGameMessage)
		}
	}
}

func (client *Client) handleGameParamsReq(gameParamsReq GameParamsReq) {
	if !client.isValid(gameParamsReq.Id) {
		return
	}

	gameId, board, err := NewGame(gameParamsReq)
	if err != nil {
		// TODO More granular error handling
		log.Println("NewGame: ", err)
	}
	client.game = gameId

	response := NewMessageGameParamsResp(board, gameId)
	client.write(response)
}

func (client *Client) handleClientMoveMessage(clientMoveMessage ClientMoveMessage) {
	if !client.isValid(clientMoveMessage.Id) {
		return
	}

	board, player1, player2, turn := GetBoard(client.game)
	if board == nil {
		// TODO handle this better
		client.Close()
		return
	}

	// Determine which player is which and whose turn it is
	var playerNum uint8
	var opponent string
	if player1 == client.id {
		playerNum = 1
		opponent = player2
	}
	if player2 == client.id {
		playerNum = 2
		opponent = player1
	}
	if client.id != turn {
		// Ignore the move since it's not their turn
		// TODO Maybe we should synchronize by sending a JoinGameResp here
		return
	}

	board.ChangePlayerColor(playerNum-1, Color(clientMoveMessage.Color))
	// TODO check for win

	SetBoard(client.game, board, opponent)

	// Respond to the mover, it is no longer their turn
	moverResponse := NewMessageMoveMessage(board, false)
	client.write(moverResponse)

	// Send move to their opponent, it is now their turn
	opponentMessage := NewMessageMoveMessage(board, true)
	TryWrite(opponent, opponentMessage)
}

func (client *Client) handleJoinGameReq(joinGameReq JoinGameReq) {
	if !client.isValid(joinGameReq.Id) {
		return
	}

	success, board, opponent, turn := TryJoin(client.id, joinGameReq.GameId)
	if success {
		client.game = joinGameReq.GameId
	}
	response := NewMessageJoinGameResp(success, BOT_PLAYER, board, turn)
	client.write(response)

	// Notify the opponent that another client has joined their game
	opponentMessage := NewMessageOtherClientJoin(!turn)
	TryWrite(opponent, opponentMessage)
}
