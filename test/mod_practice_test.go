package test

import (
	"fmt"
	"github.com/fightingBald/py-ds/counter"
	"testing"
)

func TestCounter(t *testing.T) {
	c := counter.Counter{}
	c.Add("apple")
	c.Add("apple")
	fmt.Println("apple count:", c.Count("apple")) // apple count: 2
}
