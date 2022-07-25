package processor

// Processes seeds and determines the next seed. Should occur before wipe.
func (wipedata *WipeData) ProcessSeeds() bool {

	return true
}

// Gets the next seed in the array.
func (wipedata *WipeData) GetNextSeed(seeds []int, curseed int) int {
	seed := -1

	for v, s := range seeds {
		if curseed == s {
			// If we're on the last seed, return 0 as the array item (starting item). Otherwise, return index.
			if (len(seeds) - 1) == v {
				s = 0
			} else {
				s = v
			}
		}
	}

	return seed 
}
