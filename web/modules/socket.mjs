
/**
 * @type {WebSocket}
 */
export var ws
export var id = ""

export const WS_SERVER_ADDRESS = "ws://localhost/ws/"

const ID_MESSAGE_STR = "id_message"
const GAME_PARAMS_REQ_STR = "game_params_req"
const GAME_PARAMS_RESP_STR = "game_params_resp"
const MOVE_MESSAGE_STR = "move_message"
const CLIENT_MOVE_MESSAGE_STR = "client_move_message"

/**
 * @param {string} address - websocket server address to connect to
 */
export function connect(address) {
    ws = new WebSocket(address, "JSON-v1")
    getId()
    ws.onmessage = messageHandler

    // Send id upon connection
    ws.onopen = (_) => {
        ws.send(JSON.stringify({
            message_type: "id_message",
            message: { id: id }
        }));
    };
}

/**
 * @param {any} event - WebSocket message event
 */
function messageHandler(event) {
    console.log(event.data)
    let message = JSON.parse(event.data)

    // Call appropriate handler
    switch (message.message_type) {
        case ID_MESSAGE_STR:
            updateId(message.message.id)
            break
        case GAME_PARAMS_RESP_STR:
            handleGameParamsResp(message.message)
            break
        case MOVE_MESSAGE_STR:
            handleMoveMessage(message.message)
            break
        default:
            console.log(`Received unsupported message from server: $event.data`)

    }
}

/**
 * @param {string} newId
 */
function updateId(newId) {
    id = newId
    localStorage.setItem("id", newId)
}

function getId() {
    id = localStorage.getItem("id")
    if (id == null) {
        id = ""
    }
}

function handleGameParamsResp(game_params_resp) {

}

function handleMoveMessage(move_message) {

}
