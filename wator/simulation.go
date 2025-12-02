package main

import (
    "fmt"
    "math/rand"
    "os"
    "sync"
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
    start := time.Now()
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

    chronon := 0

    for {
        chronon++

        // advance one chronon (potentially using multiple threads)
        w = StepWorld(w, cfg, rnd)

        // draw occasionally (only with small grids / Threads=1 ideally)
        if cfg.DrawEvery > 0 && chronon%cfg.DrawEvery == 0 {
            drawWorld(w, chronon)
        }

        // stop if either species is extinct
        if countEntities(w, Fish) == 0 || countEntities(w, Shark) == 0 {
            break
        }

        // optional chronon limit
        if cfg.Chronons > 0 && chronon >= cfg.Chronons {
            break
        }
    }

    elapsed := time.Since(start)
    fmt.Printf("Threads: %d  Time: %v\n", cfg.Threads, elapsed)

    // If a benchmark file was provided, append a CSV line
    writeBenchmarkLine(cfg, elapsed)
}

//  @brief Writes one line of benchmark CSV if BenchFile is set
func writeBenchmarkLine(cfg Config, elapsed time.Duration) {
    if cfg.BenchFile == "" {
        return
    }

    f, err := os.OpenFile(cfg.BenchFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        fmt.Printf("Could not open benchmark file %s: %v\n", cfg.BenchFile, err)
        return
    }
    defer f.Close()

    // If file is empty, write a header row
    info, err := f.Stat()
    if err == nil && info.Size() == 0 {
        fmt.Fprintln(f, "Threads,GridSize,NumFish,NumShark,FishBreed,SharkBreed,Starve,Chronons,TimeMillis")
    }

    millis := elapsed.Milliseconds()

    // One CSV row per run
    fmt.Fprintf(
        f,
        "%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
        cfg.Threads,
        cfg.GridSize,
        cfg.NumFish,
        cfg.NumShark,
        cfg.FishBreed,
        cfg.SharkBreed,
        cfg.Starve,
        cfg.Chronons,
        millis,
    )
}

//  @brief Prints the current world grid to the terminal
func drawWorld(w *World, chronon int) {
    fmt.Printf("Chronon: %d\n", chronon)

    for row := 0; row < w.Size; row++ {
        for col := 0; col < w.Size; col++ {
            cell := w.Cells[row][col]
            switch cell.Entity {
            case Empty:
                fmt.Print("~")
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

//  @brief Advances the world by one chronon (multi-threaded using goroutines)
func StepWorld(w *World, cfg Config, rnd *rand.Rand) *World {
    next := newEmptyWorldLike(w)

    threads := cfg.Threads
    if threads < 1 {
        threads = 1
    }
    if threads > w.Size {
        // no point having more threads than rows
        threads = w.Size
    }

    rowsPerThread := w.Size / threads
    remainder := w.Size % threads

    var wg sync.WaitGroup
    var mu sync.Mutex // protects writes to "next"

    startRow := 0
    for t := 0; t < threads; t++ {
        extra := 0
        if t < remainder {
            extra = 1
        }
        endRow := startRow + rowsPerThread + extra

        wg.Add(1)

        go func(start, end int) {
            defer wg.Done()

            // per-goroutine RNG
            localRnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(start)))

            for row := start; row < end; row++ {
                for col := 0; col < w.Size; col++ {

                    cell := w.Cells[row][col]
                    if cell.Entity == Empty {
                        continue
                    }

                    switch cell.Entity {
                    case Fish:
                        stepFish(w, next, row, col, cfg, localRnd, &mu)
                    case Shark:
                        stepShark(w, next, row, col, cfg, localRnd, &mu)
                    }
                }
            }

        }(startRow, endRow)

        startRow = endRow
    }

    wg.Wait()

    return next
}

//  @brief Handles movement and reproduction for a single fish at (row, column)
func stepFish(current *World, next *World, row, col int, cfg Config, rnd *rand.Rand, mu *sync.Mutex) {
    cell := current.Cells[row][col]
    neighbors := current.Neighbors(row, col)

    emptySpots := make([][2]int, 0)

    // Look for empty neighbors in CURRENT world (not next)
    for _, n := range neighbors {
        nr, nc := n[0], n[1]
        if current.Cells[nr][nc].Entity == Empty {
            emptySpots = append(emptySpots, [2]int{nr, nc})
        }
    }

    // No movement
    if len(emptySpots) == 0 {
        mu.Lock()
        next.Cells[row][col] = Cell{
            Entity:     Fish,
            BreedTimer: cell.BreedTimer + 1,
        }
        mu.Unlock()
        return
    }

    // Pick random move
    destination := emptySpots[rnd.Intn(len(emptySpots))]
    nr, nc := destination[0], destination[1]

    // Reproduction happens only ON MOVE
    if cell.BreedTimer+1 >= cfg.FishBreed {
        mu.Lock()
        // Leave baby at original position
        next.Cells[row][col] = Cell{
            Entity:     Fish,
            BreedTimer: 0,
        }
        // Parent moves
        next.Cells[nr][nc] = Cell{
            Entity:     Fish,
            BreedTimer: 0,
        }
        mu.Unlock()
        return
    }

    // Normal movement
    mu.Lock()
    next.Cells[nr][nc] = Cell{
        Entity:     Fish,
        BreedTimer: cell.BreedTimer + 1,
    }
    mu.Unlock()
}

//  @brief Handles movement, eating, reproduction and starvation for a shark at (row, column).
func stepShark(current *World, next *World, row, col int, cfg Config, rnd *rand.Rand, mu *sync.Mutex) {
    cell := current.Cells[row][col]

    // Shark loses 1 energy each turn
    newEnergy := cell.Energy - 1
    if newEnergy <= 0 {
        return // shark dies
    }

    neighbors := current.Neighbors(row, col)

    // 1. LOOK FOR FISH TO EAT
    fishTargets := make([][2]int, 0)
    for _, n := range neighbors {
        nr, nc := n[0], n[1]
        if current.Cells[nr][nc].Entity == Fish {
            fishTargets = append(fishTargets, [2]int{nr, nc})
        }
    }

    if len(fishTargets) > 0 {
        destination := fishTargets[rnd.Intn(len(fishTargets))]
        nr, nc := destination[0], destination[1]

        // Eating gives FULL energy
        gainedEnergy := cfg.Starve

        mu.Lock()
        defer mu.Unlock()

        // Reproduction?
        if cell.BreedTimer+1 >= cfg.SharkBreed {
            // Leave baby behind with HALF energy
            next.Cells[row][col] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     gainedEnergy / 2,
            }
            // Parent moves to fish
            next.Cells[nr][nc] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     gainedEnergy,
            }
            return
        }

        // Normal move & eat
        next.Cells[nr][nc] = Cell{
            Entity:     Shark,
            BreedTimer: cell.BreedTimer + 1,
            Energy:     gainedEnergy,
        }
        return
    }

    // 2. NO FISH — MOVE LIKE FISH
    emptyTargets := make([][2]int, 0)
    for _, n := range neighbors {
        nr, nc := n[0], n[1]
        if current.Cells[nr][nc].Entity == Empty {
            emptyTargets = append(emptyTargets, [2]int{nr, nc})
        }
    }

    if len(emptyTargets) > 0 {
        destination := emptyTargets[rnd.Intn(len(emptyTargets))]
        nr, nc := destination[0], destination[1]

        mu.Lock()
        defer mu.Unlock()

        // Reproduce?
        if cell.BreedTimer+1 >= cfg.SharkBreed {
            next.Cells[row][col] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     newEnergy / 2,
            }
            next.Cells[nr][nc] = Cell{
                Entity:     Shark,
                BreedTimer: 0,
                Energy:     newEnergy,
            }
            return
        }

        next.Cells[nr][nc] = Cell{
            Entity:     Shark,
            BreedTimer: cell.BreedTimer + 1,
            Energy:     newEnergy,
        }
        return
    }

    // 3. Can't move
    mu.Lock()
    next.Cells[row][col] = Cell{
        Entity:     Shark,
        BreedTimer: cell.BreedTimer + 1,
        Energy:     newEnergy,
    }
    mu.Unlock()
}
