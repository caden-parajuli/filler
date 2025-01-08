export var gameId = 0;

export function hideHome() {
    document.getElementById("home").style.display = "none";
}

export function showWaitingRoom() {
    document.getElementById("waiting-room").style.display = "initial";
}

export function hideWaitingRoom() {
    document.getElementById("waiting-room").style.display = "none";
}

export function showGame() {
    document.getElementById("game").style.display = "initial";
}

/**
 * @param {number} newGameId 
 */
export function setGameId(newGameId) {
    gameId = newGameId;
    let gameIdDisplay = document.getElementById("game-id-display");
    if (gameId != 0) {
        // Show the game id
        gameIdDisplay.style.display = "initial";
    } else {
        gameIdDisplay.style.display = "none";
    }
    gameIdDisplay.innerHTML = `Room Code: ${gameId}`
}

