package processor

// Processes seeds and determines the next seed. Should occur before wipe.
func (wipedata *WipeData) ProcessSeeds() bool {

	return true
}

// Gets the next seed in the array.
func (wipedata *WipeData) GetNextSeed(seeds []int, curseed int, picktype int) int {
	seed := -1

	if picktype == 1 {
		for v, s := range seeds {
			if curseed == s {
				// If we're on the last seed, return 0 as the array item (starting item). Otherwise, return index.
				if (len(seeds) - 1) == v {
					seed = 0
				} else {
					seed = v
				}
			}
		}
	} else {
		seed = rand.Intn((len(seeds) - 1) + 1) + 0
	}

	return seed 
}
