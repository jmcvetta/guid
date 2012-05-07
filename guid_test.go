// Copyright 2012 Jason McVetta.  This is Free Software, released under
// an MIT-style license.  See README.md for more info.

package guid

import (
	"testing"
)

func TestUniqueness(t *testing.T) {
}

func BenchmarkIdGeneration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		guid, err := NextId()
		if err != nil {
			b.Fatal(err)
		}
		println(guid)
	}
}

func BenchmarkParallel10(b *testing.B) {
	parallelIdGeneration(10, b)
}

func parallelIdGeneration(c int, b *testing.B) {
	// Setup the workers
	reqs := make(chan bool)
	guids := make(chan int64)
	for i := 0; i < c; i++ {
		go func() {
			for {
				<-reqs
				g, err := NextId()
				if err != nil {
					b.Fatal(err)
				}
				println(g)
				guids <- g
			}
		}()
	}
	// Request some GUIDs
	for n := 0; n < b.N; n++ {
		reqs <- true
	}
	// Wait for GUIDs to be generated
	for n := 0; n < b.N; n++ {
		<-guids
	}
}
