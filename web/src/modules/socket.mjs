import { numColors, numRows, rowLen, setPlayerNum, setBoard } from "./board.mjs"
import { setTurn } from "./controls.mjs"
import { showWaitingRoom } from "./room.mjs"
import { hideHome, hideWaitingRoom, showGame, setGameId } from "./room.mjs"

/**
 * @type {WebSocket}
 */
export var ws;
export var id = "";

export const WS_SERVER_ADDRESS = "ws://localhost/ws/";

const ID_MESSAGE_STR = "id_message";
const GAME_PARAMS_REQ_STR = "game_params_req";
const GAME_PARAMS_RESP_STR = "game_params_resp";
const MOVE_MESSAGE_STR = "move_message";
const CLIENT_MOVE_MESSAGE_STR = "client_move_message";
const JOIN_GAME_REQ_STR = "join_game_req";
const JOIN_GAME_RESP_STR = "join_game_resp";
const OTHER_CLIENT_JOIN_STR = "other_client_join";

/**
 * @param {string} address - websocket server address to connect to
 */
export function connect(address) {
    ws = new WebSocket(address, "JSON-v1");
    getId();
    ws.onmessage = messageHandler;

    // Send id upon connection
    ws.onopen = (_) => {
        ws.send(JSON.stringify({
            message_type: ID_MESSAGE_STR,
            message: { id: id }
        }));
    };
}

export function requestGame() {
    console.log("Requesting game");
    /** @type {GameParamsReq} */
    let game_params_req = {
        id: id,
        is_diamonds: true,
        num_rows: numRows,
        num_cols: rowLen,
        num_colors: numColors,
    };

    ws.send(JSON.stringify({
        message_type: GAME_PARAMS_REQ_STR,
        message: game_params_req
    }));
}

/**
 * @param {number} game_id
 */
export function joinGame(game_id) {
    console.log(`Joining game ${game_id}`);
    /** @type {JoinGameReq} */
    let join_game_req = {
        id,
        game_id
    };

    ws.send(JSON.stringify({
        message_type: JOIN_GAME_REQ_STR,
        message: join_game_req
    }));
}

/**
 * @param {number} color
 */
export function sendMove(color) {
    console.log(`Sending move ${color}`);
    /** @type ClientMoveMessage */
    let client_move_message = {
        id,
        color
    };

    ws.send(JSON.stringify({
        message_type: CLIENT_MOVE_MESSAGE_STR,
        message: client_move_message
    }));
}

/**
 * @param {MessageEvent} event - WebSocket message event
 */
function messageHandler(event) {
    console.log(event.data);
    /** @type {Message} */
    let message = JSON.parse(event.data);

    // Call appropriate handler
    switch (message.message_type) {
        case ID_MESSAGE_STR:
            updateId(message.message);
            break;
        case GAME_PARAMS_RESP_STR:
            handleGameParamsResp(message.message);
            break;
        case MOVE_MESSAGE_STR:
            handleMoveMessage(message.message);
            break;
        case JOIN_GAME_RESP_STR:
            handleJoinGameResp(message.message);
            break;
        case OTHER_CLIENT_JOIN_STR:
            handleOtherClientJoin(message.message);
            break;
        default:
            console.log(`Received unsupported message from server: ${event.data}`)
    }
}

/**
 * @param {IdMessage} idMessage
 */
function updateId(idMessage) {
    id = idMessage.id;
    localStorage.setItem("id", id);
}

function getId() {
    id = localStorage.getItem("id")
    if (id == null) {
        id = "";
    }
}

/**
 * @param {GameParamsResp} game_params_resp
 */
function handleGameParamsResp(game_params_resp) {
    setPlayerNum(0)
    setBoard(game_params_resp.board)
    setGameId(game_params_resp.game_id)
    hideHome()
    showWaitingRoom()
}

/**
 * @param {MoveMessage} move_message
 */
function handleMoveMessage(move_message) {
    setBoard(move_message.board)
    setTurn(move_message.my_turn)
}

/**
 * @param {JoinGameResp} join_game_resp
 */
function handleJoinGameResp(join_game_resp) {
    setPlayerNum(join_game_resp.player_num)
    setBoard(join_game_resp.board)
    setTurn(join_game_resp.my_turn)
    hideHome()
    showGame()
}

/**
 * @param {OtherClientJoin} other_client_join
 */
function handleOtherClientJoin(other_client_join) {
    setTurn(other_client_join.my_turn)
    hideWaitingRoom()
    showGame()
}
