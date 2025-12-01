package main

import (
    "fmt"
    "math/rand"
    "time"
)

/**
    @file simulation.go
    @brief Contains the main simulation loop for the Wa-Tor world
    This module is responsible for:
        Iterating over chronons/time steps
        Updating fish and shark positions and their state
        Applying reproduction and starvation rules
        Drawing the world
*/

//  @brief Creates a new World with the same size and parameters as an existing one, but with all cells empty
func newEmptyWorldLike(w *World) *World {
    cells := make([][]Cell, w.Size)
    for row := range cells {
        cells[row] = make([]Cell, w.Size)
    }
    return &World{
        Size:       w.Size,
        Cells:      cells,
        FishBreed:  w.FishBreed,
        SharkBreed: w.SharkBreed,
        Starve:     w.Starve,
    }
}

//  @brief Runs the Wa-Tor simulation using the given configuration
//   @param "cfg" The simulation configuration
//  @param "w" The initial world
func RunSimulation(cfg Config, w *World) {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

    chronon := 0

    for {
        chronon++

        //  Step the world by one chronon
        w = StepWorld(w, cfg, rnd)

        //  Draw every cfg.DrawEvery chronons
        if cfg.DrawEvery > 0 && chronon%cfg.DrawEvery == 0 {
            drawWorld(w, chronon)
        }

        // extinction stopping conditions
        if countEntities(w, Fish) == 0 {
            fmt.Println("All fish extinct. Simulation ending.")
            break
        }
        if countEntities(w, Shark) == 0 {
            fmt.Println("All sharks extinct. Simulation ending.")
            break
        }

        if cfg.Chronons > 0 && chronon >= cfg.Chronons {
            fmt.Println("Reached chronon limit.")
            break
        }
    }
}

//  @brief Prints the current world grid to the terminal
func drawWorld(w *World, chronon int) {
    fmt.Printf("Chronon: %d\n", chronon)

    for row := 0; row < w.Size; row++ {
        for col := 0; col < w.Size; col++ {
            cell := w.Cells[row][col]
            switch cell.Entity {
            case Empty:
                fmt.Print(".")
            case Fish:
                fmt.Print("F")
            case Shark:
                fmt.Print("S")
            }
        }
        fmt.Println()
    }

    fmt.Printf("Fish: %d  Sharks: %d\n", countEntities(w, Fish), countEntities(w, Shark))
    fmt.Println()
}

//  @brief Counts how many cells currently contain the given entity type
func countEntities(w *World, e Entity) int {
    count := 0
    for row := 0; row < w.Size; row++ {
        for col := 0; col < w.Size; col++ {
            if w.Cells[row][col].Entity == e {
                count++
            }
        }
    }
    return count
}

//  @brief Advances the world by one chronon
func StepWorld(w *World, cfg Config, rnd *rand.Rand) *World {
    next := newEmptyWorldLike(w)

    for row := 0; row < w.Size; row++ {
        for col := 0; col < w.Size; col++ {

            // prevent double-processing
            if next.Cells[row][col].Entity != Empty {
                continue
            }

            cell := w.Cells[row][col]

            switch cell.Entity {
            case Fish:
                stepFish(w, next, row, col, cfg, rnd)

            case Shark:
                stepShark(w, next, row, col, cfg, rnd)

            case Empty:
            }
        }
    }

    return next
}

//  @brief Handles movement and reproduction for a single fish at (row, column)
func stepFish(current *World, next *World, row, col int, cfg Config, rnd *rand.Rand) {
    cell := current.Cells[row][col]
    neighbors := current.Neighbors(row, col)
    emptySpots := make([][2]int, 0)

    for _, n := range neighbors {
        nr, nc := n[0], n[1]

        if current.Cells[nr][nc].Entity == Empty && next.Cells[nr][nc].Entity == Empty {
            emptySpots = append(emptySpots, [2]int{nr, nc})
        }
    }

    if len(emptySpots) == 0 {
        next.Cells[row][col] = Cell{
            Entity:     Fish,
            BreedTimer: cell.BreedTimer + 1,
        }
        return
    }

    destination := emptySpots[rnd.Intn(len(emptySpots))]
    nr, nc := destination[0], destination[1]

    if cell.BreedTimer+1 >= cfg.FishBreed {
        next.Cells[row][col] = Cell{
            Entity:     Fish,
            BreedTimer: 0,
        }
        next.Cells[nr][nc] = Cell{
            Entity:     Fish,
            BreedTimer: 0,
        }
        return
    }

    next.Cells[nr][nc] = Cell{
        Entity:     Fish,
        BreedTimer: cell.BreedTimer + 1,
    }
}

//  @brief Handles movement, eating, reproduction and starvation for a shark at (row, column).
func stepShark(current *World, next *World, row, col int, cfg Config, rnd *rand.Rand) {
    cell := current.Cells[row][col]

    newEnergy := cell.Energy - 1
    if newEnergy <= 0 {
        return
    }

    neighbors := current.Neighbors(row, col)

    fishTarget := make([][2]int, 0)
    for _, n := range neighbors {
        nr, nc := n[0], n[1]
        if current.Cells[nr][nc].Entity == Fish && next.Cells[nr][nc].Entity == Empty {
            fishTarget = append(fishTarget, [2]int{nr, nc})
        }
    }

    if len(fishTarget) > 0 {
        destination := fishTarget[rnd.Intn(len(fishTarget))]
        nr, nc := destination[0], destination[1]

        if cell.BreedTimer+1 >= cfg.SharkBreed {
            next.Cells[row][col] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     cfg.Starve,
            }
            next.Cells[nr][nc] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     cfg.Starve,
            }
        } else {
            next.Cells[nr][nc] = Cell{
                Entity:     Shark,
                BreedTimer: cell.BreedTimer + 1,
                Energy:     cfg.Starve,
            }
        }
        return
    }

    emptyTarget := make([][2]int, 0)
    for _, n := range neighbors {
        nr, nc := n[0], n[1]
        if current.Cells[nr][nc].Entity == Empty && next.Cells[nr][nc].Entity == Empty {
            emptyTarget = append(emptyTarget, [2]int{nr, nc})
        }
    }

    if len(emptyTarget) > 0 {
        destination := emptyTarget[rnd.Intn(len(emptyTarget))]
        nr, nc := destination[0], destination[1]

        if cell.BreedTimer+1 >= cfg.SharkBreed {
            next.Cells[row][col] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     cfg.Starve,
            }
            next.Cells[nr][nc] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     newEnergy,
            }
        } else {
            next.Cells[nr][nc] = Cell{
                Entity:     Shark,
                BreedTimer: cell.BreedTimer + 1,
                Energy:     newEnergy,
            }
        }
        return
    }

    next.Cells[row][col] = Cell{
        Entity:     Shark,
        BreedTimer: cell.BreedTimer + 1,
        Energy:     newEnergy,
    }
}
