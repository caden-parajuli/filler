export const rowLen = 30;
// Must be odd
export const numRows = 15;
export const numColors = 7;

// Useful for console debugging
export const CYAN = "color0"
export const GREEN = "color1";
export const RED = "color2";
export const YELLOW = "color3";
export const BLUE = "color4";
export const PINK = "color5";
export const ORANGE = "color6";

export const COLORS = [CYAN, GREEN, RED, YELLOW, BLUE, PINK, ORANGE]
const UNCOLORED = "uncolored"

/** @type {[number, number][]} */
export const player_positions = [[0, rowLen - 1], [numRows - 1, 0]]

export var player_num = -1;

/** @type {Board} */
export var board

/** @type {string} */
var player_color;

/**
 * @param {[number, number]} pos
 * @returns {[number, number][]} the neighbors of a diamond in the order [topLeft, topRight, bottomleft, bottomRight], whichever ones exist
 */
function getNeighbors(pos) {
    // Not exactly elegant, but it's good enough
    let i = pos[0];
    let j = pos[1];
    if (i & 1) {
        return [
            [i - 1, j],
            [i - 1, j + 1],
            [i + 1, j],
            [i + 1, j + 1],
        ];
    }

    /** @type {[number, number][]} */
    let result = [];
    if (i > 0 && j > 0) {
        result.push([i - 1, j - 1]);
    }
    if (i > 0 && j < rowLen - 1) {
        result.push([i - 1, j]);
    }
    if (i < numRows - 1 && j > 0) {
        result.push([i + 1, j - 1]);
    }
    if (i < numRows - 1 && j < rowLen - 1) {
        result.push([i + 1, j]);
    }
    return result;
}

/**
 * @param {HTMLElement} diamond
 * @param {string} color - the color to change to, e.g "color1"
 */
function changeSingleColor(diamond, color) {
    diamond.classList.replace(diamond.dataset.color, color);
    diamond.dataset.color = color;
}

/**
 * @param {[number, number]} pos
 * @returns {HTMLElement}
 */
function getDiamond(pos) {
    return document.querySelector(`[data-pos="${pos[0]},${pos[1]}"]`);
}

/**
 * Changes the player color.
 * Used to anticipate a move for responsiveness until the server responds.
 *
 * @param {[number, number]} start_pos - the position of the first square
 * @param {string} color - the color to change to, e.g "color1"
 */
export function changePlayerColor(start_pos, color) {
    let first_diamond = getDiamond(start_pos);
    let before_color = first_diamond.dataset.color;
    if (before_color == color) {
        return;
    }
    changeSingleColor(first_diamond, color);

    let to_visit = [start_pos];
    while (to_visit.length != 0) {
        let pos = to_visit.pop();

        // Visit the neighbors, changing their color and queueing if their color matches
        let neighbors = getNeighbors(pos);
        while (neighbors.length != 0) {
            let neighbor_pos = neighbors.pop();
            let neighbor = getDiamond(neighbor_pos);

            if (neighbor.dataset.color == before_color) {
                changeSingleColor(neighbor, color);
                to_visit.push(neighbor_pos);
            }
        }
    }
}

/**
 * @param {Board} newBoard
 */
export function setBoard(newBoard) {
    board = newBoard
    for (let i = 0; i < (board.num_rows - 1) / 2; ++i) {
        for (let j = 0; j < board.num_cols; ++j) {
            let diamond = getDiamond([2 * i, j]);
            changeSingleColor(diamond, COLORS[board.data[2 * i][j]]);
        }
        for (let j = 0; j < board.num_cols - 1; ++j) {
            let diamond = getDiamond([2 * i + 1, j]);
            changeSingleColor(diamond, COLORS[board.data[2 * i + 1][j]]);
        }
    }
    for (let j = 0; j < board.num_cols; ++j) {
        let diamond = getDiamond([board.num_rows - 1, j]);
        changeSingleColor(diamond, COLORS[board.data[board.num_rows - 1][j]]);
    }
    player_color = getDiamond(player_positions[player_num]).dataset.color;
}

export function makeEmptyGrid() {
    let grid = document.getElementById("grid");

    for (let i = 0; i < (numRows - 1) / 2; ++i) {
        // Non-offset rows
        let on_row = document.createElement("div");
        on_row.className = "on-row";
        if (i == 0) {
            on_row.className = "first-row";
        }

        for (let j = 0; j < rowLen; ++j) {
            let diamond = document.createElement("div");
            diamond.className = "diamond";

            let row = (2 * i).toString();
            let col = j.toString();
            setNewDiamond(diamond, row, col);

            on_row.appendChild(diamond);
        }
        grid.appendChild(on_row);

        // Offset rows
        let off_row = document.createElement("div");
        off_row.className = "off-row";

        for (let j = 0; j < rowLen - 1; ++j) {
            let diamond = document.createElement("div");
            diamond.className = "diamond";

            let row = (2 * i + 1).toString();
            let col = j.toString();
            setNewDiamond(diamond, row, col);

            off_row.appendChild(diamond);
        }
        grid.appendChild(off_row);
    }

    // Last row
    let on_row = document.createElement("div");
    on_row.className = "on-row";
    for (let j = 0; j < rowLen; ++j) {
        let diamond = document.createElement("div");
        diamond.className = "diamond";

        let row = (numRows - 1).toString();
        let col = j.toString();
        setNewDiamond(diamond, row, col);

        on_row.appendChild(diamond);
    }
    grid.appendChild(on_row);
}

/**
 * @param {HTMLElement} diamond
 * @param {string} row
 * @param {string} col
 */
function setNewDiamond(diamond, row, col) {
    diamond.dataset.row = row;
    diamond.dataset.col = col;
    diamond.dataset.pos = row + "," + col;

    diamond.dataset.color = UNCOLORED;
    diamond.classList.add(UNCOLORED);

}

/**
 * @param {number} new_player_num
 */
export function setPlayerNum(new_player_num) {
    player_num = new_player_num
}
