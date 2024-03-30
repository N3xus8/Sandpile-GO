// https://www.youtube.com/watch?v=SfWWaZ1AoQE
//
// If all tiles have less that 4 grains of sand the sandpile is stable
// If at least one vertex is unstable
//  the whole configuration  is said to be unstable.
// In this case, choose any unstable vertex at random.
// Topple this vertex by reducing its grain number by four and by increasing the grain numbers of each of its (at maximum four) direct neighbors by one, i.e. set

// Use the raylib package for graphics
// Implemented on 29/03/2024.
// go version go1.21.8
// gcc version 8.1.0 (x86_64-posix-seh-rev0, Built by MinGW-W64 project)
package main

import (
	"fmt"
	"math/rand"

	raylb "github.com/gen2brain/raylib-go/raylib"
)

// The raylib window
const WINDOWWIDTH int32 = 800
const WINDOWHEIGHT int32 = 800

type Sandpile struct {
	width  int32
	height int32
	grains []int32
}

type Rectangle struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// Color are used to display the number of grains of sand
var Palette = [5]raylb.Color{
	raylb.Black,
	raylb.Purple,
	raylb.DarkGreen,
	raylb.Orange,
	raylb.Magenta,
}

// Use a Map with Key Value pair. In this case the value is unused
// Probably there's a better data structure for this.
// Basically the map registers the tiles with more than 4 grains of sand.
// The keys correspond to the index of the tiles on the sandpile.
var excessMap = make(map[int32]bool)

// function that return all the keys for a given map
// Uses Generics
// the ~ indicates derived map
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// with/without outline around the tiles
var outline bool = false

// function that draw the sandpile
func draw_sandpile(location Rectangle, sandpile *Sandpile) {

	//
	var cellWidth float32 = float32((location.Width) / sandpile.width)
	var cellHeight float32 = float32((location.Height) / sandpile.height)

	cell := Rectangle{X: int32(0), Y: int32(0), Width: int32(0), Height: int32(0)}

	for row := int32(0); row < sandpile.height; row++ {
		for col := int32(0); col < sandpile.width; col++ {
			cell = Rectangle{
				X:      location.X + col*int32(cellWidth),
				Y:      location.Y + row*int32(cellHeight),
				Width:  int32(cellWidth),
				Height: int32(cellHeight),
			}
			grains := sandpile.grains[row*sandpile.width+col]
			if grains >= 4 {
				grains = 4
			}
			var color raylb.Color = Palette[grains]

			raylb.DrawRectangle(cell.X, cell.Y, cell.Width, cell.Height, color)
			if outline {
				raylb.DrawRectangleLines(cell.X, cell.Y, cell.Width, cell.Height, raylb.RayWhite)
			}

			// convert custom rectangle struct to raylib rectangle
			// to be able to use the checkCollisionPointRec function.
			// Some raylib function use int32 for coordinates but the raylib rectangle uses float.
			var cellf = raylb.Rectangle{
				X:      float32(cell.X),
				Y:      float32(cell.Y),
				Width:  float32(cell.Width),
				Height: float32(cell.Height),
			}
			var mousePosition raylb.Vector2 = raylb.GetMousePosition()
			if raylb.CheckCollisionPointRec(mousePosition, cellf) {
				if raylb.IsMouseButtonPressed(raylb.MouseButtonLeft) {
					sandpile.grains[row*sandpile.width+col] += 1
					if sandpile.grains[row*sandpile.width+col] >= 4 {
						excessMap[row*sandpile.width+col] = true
					}
				}
			}
		}
	}
}

func updatePile(sandpile *Sandpile) {
	if len(excessMap) == 0 {
		return
	}

	keys := Keys(excessMap)             // Use the Keys function to retrieve all keys in the map. Basically all tiles >= 4 grains of sand.
	randomIndex := rand.Intn(len(keys)) // Pick random index(keys of the map) in the slice
	updateIndex := keys[randomIndex]    // Pick random keys

	sandpile.grains[updateIndex] -= 4     // remove 4 grains
	if sandpile.grains[updateIndex] < 4 { // if it becomes < 4 remove from the map.
		delete(excessMap, updateIndex)
	}
	// Retrieve the row and column from the index
	row := updateIndex / sandpile.width
	col := updateIndex % sandpile.width

	// distribute excess grains of sand to neighbour tiles E-W-N-S
	// loop through neighbours
	// [-1,0] [-1,1] [0, 1] [1, 1] [1, 0] [1, -1] [0, -1] [-1, -1]
	for dx := int32(-1); dx <= 1; dx++ {
		for dy := int32(-1); dy <= 1; dy++ {
			if dx == 0 && dy == 0 { // center do nothing
				continue
			}
			if dx != 0 && dy != 0 { // diagonal do nothing
				continue
			}

			newRow := row + dy
			newCol := col + dx

			if newRow < 0 || newRow >= sandpile.height { // outside indices
				continue
			}
			if newCol < 0 || newCol >= sandpile.width { // outside indices
				continue
			}

			// add 1 grain to neighbours
			index := newRow*sandpile.width + newCol
			sandpile.grains[index] += 1
			if sandpile.grains[index] >= 4 {
				excessMap[int32(index)] = true
			}

		}
	}

	// for key := range excessMap {
	// 	println("this square has an excess of grains of sand: %v", key)
	// }
}

func main() {

	raylb.InitWindow(WINDOWWIDTH, WINDOWHEIGHT, "Sandpiles")
	defer raylb.CloseWindow()

	raylb.SetTargetFPS(120)

	//window := [4]int32{0, 0, WINDOWWIDTH, WINDOWHEIGHT}

	var window = Rectangle{
		X:      int32(0),
		Y:      int32(0),
		Width:  WINDOWWIDTH,
		Height: WINDOWHEIGHT,
	}

	const sandpileWidth int32 = 50
	const sandpileHeight int32 = 50

	// sandpile struct
	var sandpile = new(Sandpile)
	sandpile.width = sandpileWidth
	sandpile.height = sandpileHeight
	grains := make([]int32, sandpileWidth*sandpileHeight) // initiliaze the slice
	sandpile.grains = grains

	// Display refresh
	var timeSinceLastUpdate float32 = 0
	var updateInterval float32 = 0.0

	// Set the initial pile of sand at the center
	initialIndex := sandpileWidth*(sandpileWidth/2) + (sandpileHeight / 2)
	sandpile.grains[initialIndex] = 870
	excessMap[initialIndex] = true

	fmt.Printf("Initial sandpile \xE2\x8F\xB3: %v grains on tile No %v", sandpile.grains[initialIndex], initialIndex)

	for !raylb.WindowShouldClose() { // while loop stop on closing window
		raylb.BeginDrawing()

		raylb.ClearBackground(raylb.RayWhite)
		raylb.DrawFPS(5, 5)
		draw_sandpile(window, sandpile)

		timeSinceLastUpdate += raylb.GetFrameTime()
		if timeSinceLastUpdate >= updateInterval {
			timeSinceLastUpdate = 0.0
			updatePile(sandpile)
		}

		raylb.EndDrawing()
	}
}
