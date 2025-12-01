package main

import (
    "math/rand"
)

/**
	@file world.go
	@brief Defines the World structure and grid operations for the Wa-Tor simulation
	The World contains:
		A 2D toroidal grid of Cells
		Simulation parameters from Config
		Helper functions for movement and neighbor retrieval
*/

//	@brief World represents the Wa-Tor ocean grid and all simulation state
type World struct {
    Size  int       //	Grid dimension
    Cells [][]Cell  //	2D array storing every cell in the world

    // Parameters copied from Config for convenience
    FishBreed  int
    SharkBreed int
    Starve     int
}

/**
	@brief Creates an empty world of given size using configuration parameters
	Creating a square gird of size x size
*/
func NewWorld(cfg Config) *World {
    //	Allocate a 2D slice of Cells
    cells := make([][]Cell, cfg.GridSize)
    for i := range cells {
        cells[i] = make([]Cell, cfg.GridSize)
    }

    return &World{
        Size:       cfg.GridSize,
        Cells:      cells,
        FishBreed:  cfg.FishBreed,
        SharkBreed: cfg.SharkBreed,
        Starve:     cfg.Starve,
    }
}

/**
	@brief Wraps a grid index so the world moves in a cycle, no out of bounds, instead returning the entity back to the first row or column depending on where they moved
*/
func (w *World) wrap(i int) int {
    if i < 0 {
        return i + w.Size
    }
    if i >= w.Size {
        return i - w.Size
    }
    return i
}

/**
	@brief Returns the indices of the 4 neighboring cells
*/
func (w *World) Neighbors(row, col int) [][2]int {
    return [][2]int{
        {w.wrap(row-1), col}, //	North
        {w.wrap(row+1), col}, //	South
        {row, w.wrap(col-1)}, //	West
        {row, w.wrap(col+1)}, //	East
    }
}

/**
	@brief Randomly places sharks and fish into empty cells at the start of the simulation
*/
func (w *World) Populate(numFish, numShark int) {
    total := w.Size * w.Size

    //	Generate a list of all cell positions
    positions := make([][2]int, 0, total)
    for row := 0; row < w.Size; row++ {
        for col := 0; col < w.Size; col++ {
            positions = append(positions, [2]int{row, col})
        }
    }

    //	Shuffle positions
    rand.Shuffle(len(positions), func(i, j int) {
        positions[i], positions[j] = positions[j], positions[i]
    })

    index := 0

    //	Place sharks
    for shark := 0; shark < numShark && index < len(positions); shark++ {
        pos := positions[index]
        index++
        w.Cells[pos[0]][pos[1]].Entity = Shark
        w.Cells[pos[0]][pos[1]].Energy = w.Starve
        w.Cells[pos[0]][pos[1]].BreedTimer = 0
    }

    //	Place fish
    for fish := 0; fish < numFish && index < len(positions); fish++ {
        pos := positions[index]
        index++
        //	Only place fish in empty cells
        if w.Cells[pos[0]][pos[1]].Entity == Empty {
            w.Cells[pos[0]][pos[1]].Entity = Fish
            w.Cells[pos[0]][pos[1]].BreedTimer = 0
        }
    }
}




