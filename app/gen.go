package app

import (
	"math/rand"
	"time"
)

type Generator struct {
	rand *rand.Rand
}

func (generator *Generator) Identity() int {
	vales := []int{0, 1}
	length := len(vales)

	return vales[generator.rand.Intn(length)]
}

func CreateGen() Generator {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return Generator{
		rand: r,
	}
}
