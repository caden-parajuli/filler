import { makeEmptyGrid } from "./modules/board.mjs";
import { initControls } from "./modules/controls.mjs";
import { connect, joinGame, requestGame, WS_SERVER_ADDRESS } from "./modules/socket.mjs";
import { hideHome, showWaitingRoom } from "./modules/room.mjs"

function host() {
    hideHome();
    showWaitingRoom();
    requestGame();
}
function join() {
    let input = /** @type {HTMLInputElement} */ (document.getElementById("join-game-id")).value;
    let game_id = parseInt(input);
    joinGame(game_id);
}

document.getElementById("host-button").addEventListener("click", host)
document.getElementById("join-button").addEventListener("click", join)

function main() {
    makeEmptyGrid();
    initControls();
    connect(WS_SERVER_ADDRESS);
}

document.addEventListener("DOMContentLoaded", main, false);
