const rowLen = 30;
// Must be odd
const numRows = 15;
const numColors = 7;

// Useful for console debugging
const cyan = "color0";
const green = "color1";
const red = "color2";
const yellow = "color3";
const blue = "color4";
const pink = "color5";
const orange = "color6";

/** @type {[number, number][]} */
const player_positions = [[0, rowLen - 1], [numRows - 1, 0]]

// TODO for debugging
var player_num = 0;

/** @type {string} */
var player_color;

// In reality, the colors must come from the server
function randColor() {
  return "color" + Math.floor(Math.random() * numColors);
}

function makeGrid() {
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
      diamond.dataset.row = row;
      diamond.dataset.col = col;
      diamond.dataset.pos = row + "," + col;

      let color = randColor();
      diamond.dataset.color = color;
      diamond.classList.add(color);

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
      diamond.dataset.row = row;
      diamond.dataset.col = col;
      diamond.dataset.pos = row + "," + col;

      let color = randColor();
      diamond.dataset.color = color;
      diamond.classList.add(color);

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
    diamond.dataset.row = row;
    diamond.dataset.col = col;
    diamond.dataset.pos = row + "," + col;

    let color = randColor();
    diamond.dataset.color = color;
    diamond.classList.add(color);

    on_row.appendChild(diamond);
  }
  grid.appendChild(on_row);

  player_color = getDiamond(player_positions[player_num]).dataset.color;
}

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
 * @param {[number, number]} start_pos - the position of the first square
 * @param {string} color - the color to change to, e.g "color1"
 */
function changePlayerColor(start_pos, color) {
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

function initControls() {
  for (let i = 0; i < numColors; ++i) {
    let control = document.querySelector(`#diamond-control${i}`);
    control.addEventListener("click", (_) => {
      // TODO Do some validation first to make sure the control is pressable
      // @ts-ignore
      changePlayerColor(player_positions[player_num], control.dataset.color)
      // TODO message server
    });
  }
}


function main() {
  document.addEventListener("DOMContentLoaded", makeGrid, false);
  document.addEventListener("DOMContentLoaded", initControls, false);
  protobuf.load("../protos/board.proto" /* , Callback*/);
}

main();
