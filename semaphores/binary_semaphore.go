package semaphores

type empty struct{}

// Create a binary semaphore.  The semaphore can be initialized
// to 'up' (full) or 'down'.
func MakeBinarySemaphore(full bool) Semaphore {
	const max = 1
	var initial = 0
	if full {
		initial = 1
	}
	return MakeCountingSemaphore(initial, max)
}
