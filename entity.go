package main

/**
	@file entity.go
	@brief Defines the different types of entities in the Wa-Tor simulation
	The world consists of three possible occupants:
		Empty (no creature)
		Fish  (moves and reproduces)
		Shark (moves, eats fish, starves, and reproduces)
*/

//	@brief Entity represents what occupies a cell in the world grid.
type Entity int

const (
    Empty Entity = iota  //	 No creature in this cell, iota automatically increments the values 
    Fish                 //	 Fish entity
    Shark                //	 Shark entity
)
