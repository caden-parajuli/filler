import { getNeighbors, player_positions, player_num, getSize } from "./board.mjs"
import { ws } from "./socket.mjs";

/**
 * Calculates both players scores
 * @param {Board} gameBoard
 */
function calculateScores(gameBoard) {
    let scores = player_positions.map((start_pos) => {
        let numControlled = 0;
        let color = gameBoard.data[start_pos[0]][start_pos[1]];

        let visited = [];
        let to_visit = [start_pos];
        while (to_visit.length != 0) {
            numControlled += 1;
            let pos = to_visit.pop();
            visited[pos[0] * gameBoard.num_cols + pos[1]] = true

            // Visit the neighbors, incrementing total and queueing if their color matches
            let neighbors = getNeighbors(pos);
            while (neighbors.length != 0) {
                let neighbor_pos = neighbors.pop();

                if (gameBoard.data[neighbor_pos[0]][neighbor_pos[1]] == color && !visited[neighbor_pos[0] * gameBoard.num_cols + neighbor_pos[1]]) {
                    to_visit.push(neighbor_pos);
                }
            }
        }

        return (numControlled / getSize(gameBoard)) * 100;
    })

    return scores;
}

/**
 * @param {Board} [board]
 */
export function updateScoreboard(board) {
    let scores = calculateScores(board);
    console.debug(scores);
    checkWin(scores);

    let leftScoreText = document.getElementById("left-score-num");
    let rightScoreText = document.getElementById("right-score-num");

    leftScoreText.innerHTML = scores[1].toFixed() + "%";
    rightScoreText.innerHTML = scores[0].toFixed() + "%";
}

/**
 * @param {number[]} scores
 */
function checkWin(scores) {
    if (scores[0] > 50) {
        if (player_num == 0) {
            win();
        } else {
            lose();
        }
    } else if (scores[1] > 50) {
        if (player_num == 1) {
            win();
        } else {
            lose();
        }
    }
}

/**
 * @param {string} message
 */
function endGame(message) {
    alert(message);
    // TODO send quit message
    ws.close();
    location.reload();
}

function win() {
    endGame("You won!");
}

function lose() {
    endGame("You lost. Don't worry, you'll get 'em next time!");
}
