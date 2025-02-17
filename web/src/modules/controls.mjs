import { changePlayerColor, numColors, player_positions, player_num, player_color, COLORS, op_color as opp_color } from "./board.mjs";
import { sendMove } from "./socket.mjs";

/**
 * @type {boolean}
 */
export var is_my_turn

export function initControls() {
    for (let colorNum = 0; colorNum < numColors; ++colorNum) {
        let control = document.getElementById(`diamond-control${colorNum}`);
        control.addEventListener("click", (_) => {
            if (control.getAttribute("disabled") != "disabled") {
                changePlayerColor(player_positions[player_num], control.dataset.color);
                sendMove(colorNum)
            }
        });
    }
}


/**
 *  Deactivates all controls
 */
export function deactivateControls() {
    for (let i = 0; i < numColors; ++i) {
        let control = document.getElementById(`diamond-control${i}`);
        control.setAttribute('disabled', 'disabled');
    }
}

/**
 *  Reactivates the controls, except for the 2 players' colors
 */
export function reactivateControls() {
    for (let i = 0; i < numColors; ++i) {
        let control = document.getElementById(`diamond-control${i}`);
        if (COLORS[i] == player_color || COLORS[i] == opp_color) {
            control.setAttribute('disabled', 'disabled');
        } else {
            control.removeAttribute('disabled');
        }
    }
}

/**
 * @param {number} color
 */
export function disableControl(color) {
    let control = document.getElementById(`diamond-control${color}`);
    control.removeAttribute('disabled');
}

/**
 * @param {number} color
 */
export function enableControl(color) {
    let control = document.getElementById(`diamond-control${color}`);
    control.setAttribute('disabled', 'disabled');
}


/**
 * @param {boolean} my_turn - Whether it is this client's turn
 */
export function setTurn(my_turn) {
    if (my_turn) {
        reactivateControls()
    } else {
        deactivateControls()
    }
}
