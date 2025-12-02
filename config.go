package main

/**
	@file config.go
	@brief Configuration structure for the Wa-Tor simulation
 */

//	@brief Holds all user-configurable parameters for the simulation
type Config struct {
    NumShark   int
    NumFish    int
    FishBreed  int
    SharkBreed int
    Starve     int
    GridSize   int
    Threads    int

    Chronons   int
    DrawEvery  int
    BenchFile  string
}
