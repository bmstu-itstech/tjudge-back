package util

// Iter iterates over all k-sized combinations of n-sized set
type Iter struct {
	N, K        int
	Combination []int
}

func (c *Iter) Next() bool {
	if c.N < c.K {
		// set is too small
		return false
	}
	if len(c.Combination) != c.K {
		// initial call
		c.Combination = make([]int, c.K)
		for i := 0; i < c.K; i++ {
			c.Combination[i] = i
		}
		return true
	}

	if c.Combination[0] == c.N-c.K {
		// final combination already reached
		return false
	}

	var i int
	for i = c.K - 1; i >= 0; i-- {
		c.Combination[i]++
		if c.Combination[i] < c.N-c.K+1+i {
			break
		}
	}
	for i = i + 1; i < c.K; i++ {
		c.Combination[i] = c.Combination[i-1] + 1
	}
	return true
}
