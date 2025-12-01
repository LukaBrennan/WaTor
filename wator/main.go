package main

import (
	"flag" //	Allows for command line option parsing, used to parse optional parameters
	"fmt"  //	For printing text to terminal
	"os"   //	Provides functions interacting with the operating system
	"strconv"	//	Used to convert string to int 
)

/**
	@file main.go
	@brief Entry point for the wartor project
	
	the file handles:
	parsing different flags such as the chronons, draw frequency and benchmark
	reading in the 7 different parameters required for the simulation to work
	validation and preparation for the simulation

	@brief this is the program entrypoint

	Here is what happens:
	Parses input from the user including, Optional flags {chronons, draw, bench}
	The 7 required positional arguments:
	1. NumShark  
	2. NumFish  
	3. FishBreed  
	4. SharkBreed  
	5. Starve  
	6. GridSize  
	7. Threads 
	
	After validation, these values will be used to configure and start the Wa-Tor simulation
*/

func main() {
	/**
	    Define command-line flags
    	@param chrononsFlag  Number of chronons to run (0 = infinite)
    	@param drawFlag      Draw every N chronons
    	@param benchFlag     Output benchmark CSV file (optional)
	*/
	chrononsFlag := flag.Int("chronons", 0, "Number of chronons to run (0 = run forever)")
	drawFlag := flag.Int("draw", 1, "Draw every N chronons")
	benchFlag := flag.String("bench", "", "Write benchmark CSV to this file")

	// Read in user inputted flags for the program
	flag.Parse()

// Read the 7 required positional arguments
args := flag.Args()
if len(args) < 7 {
    fmt.Println("Usage: wa-tor NumShark NumFish FishBreed SharkBreed Starve GridSize Threads")
    os.Exit(1)
}
	//@Error checking
// Convert arguments, ensuring that String values are checked and errored correctly

numShark, err := strconv.Atoi(args[0])
if err != nil {
    fmt.Println("Error: NumShark must be an integer.")
    os.Exit(1)
}

numFish, err := strconv.Atoi(args[1])
if err != nil {
    fmt.Println("Error: NumFish must be an integer.")
    os.Exit(1)
}

fishBreed, err := strconv.Atoi(args[2])
if err != nil {
    fmt.Println("Error: FishBreed must be an integer.")
    os.Exit(1)
}

sharkBreed, err := strconv.Atoi(args[3])
if err != nil {
    fmt.Println("Error: SharkBreed must be an integer.")
    os.Exit(1)
}

starve, err := strconv.Atoi(args[4])
if err != nil {
    fmt.Println("Error: Starve must be an integer.")
    os.Exit(1)
}

gridSize, err := strconv.Atoi(args[5])
if err != nil {
    fmt.Println("Error: GridSize must be an integer.")
    os.Exit(1)
}

threads, err := strconv.Atoi(args[6])
if err != nil {
    fmt.Println("Error: Threads must be an integer.")
    os.Exit(1)
}


	//@	Validation
	// ensuing that there cant be any negative values or any other incorrect configuration for the simulation

if numShark < 0 {
    fmt.Println("Error: NumShark must be 0 or greater.")
    os.Exit(1)
}

if numFish < 0 {
    fmt.Println("Error: NumFish must be 0 or greater.")
    os.Exit(1)
}

if fishBreed <= 0 {
    fmt.Println("Error: FishBreed must be greater than 0.")
    os.Exit(1)
}

if sharkBreed <= 0 {
    fmt.Println("Error: SharkBreed must be greater than 0.")
    os.Exit(1)
}

if starve <= 0 {
    fmt.Println("Error: Starve must be greater than 0.")
    os.Exit(1)
}

if gridSize <= 1 {
    fmt.Println("Error: GridSize must be greater than 1.")
    os.Exit(1)
}

if threads < 1 {
    fmt.Println("Error: Threads must be 1 or greater.")
    os.Exit(1)
}

cfg := Config{
    NumShark:   numShark,
    NumFish:    numFish,
    FishBreed:  fishBreed,
    SharkBreed: sharkBreed,
    Starve:     starve,
    GridSize:   gridSize,
    Threads:    threads,
    Chronons:   *chrononsFlag,
    DrawEvery:  *drawFlag,
    BenchFile:  *benchFlag,
}

fmt.Printf("Loaded configuration: %+v\n", cfg)
world := NewWorld(cfg)
world.Populate(cfg.NumFish, cfg.NumShark)
RunSimulation(cfg, world)

}
