package main

/**
	@file cell.go
	@brief Defines the Cell structure used in the Wa-Tor simulation grid
	Each cell represents one position in the ocean and may contain:
		(1) fish
		(2) shark
		(3) empty
	Cells also track breeding timers and (for sharks) energy levels
*/

//	@brief Cell stores information about a single grid tile in the simulation
type Cell struct {
    Entity     Entity //	What occupies the cell (Empty, Fish, Shark)
    BreedTimer int    //	Counts how many chronons since last reproduction

    //	Only used by sharks
    Energy int //	Remaining energy before starvation
}
