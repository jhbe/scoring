package main

import (
	"log"
	"os"
)

func main() {
	result, err := Calculate([]Race{
		{Finishers: []uint{39, 55, 38, 12, 13, 99, 69}},
		{Finishers: []uint{13, 99, 38, 12, 39, 55, 69}},
		{Finishers: []uint{13, 38, 99, 12, 39, 55, 69}},
		{Finishers: []uint{38, 39, 99, 13, 12, 69, 55}},
		{Finishers: []uint{38, 39, 13, 12, 99, 55, 69}},
		{Finishers: []uint{55, 38, 99, 13, 39, 12, 69}},
		{Finishers: []uint{39, 99, 69, 12}, DidNotFinish: []uint{13, 55}, AveragePoints: []uint{38}},
		{Finishers: []uint{39, 12, 69, 99}, DidNotFinish: []uint{38, 55}, AveragePoints: []uint{13}},
	})
	if err != nil {
		log.Fatalln(err)
	}
	var skippers = map[uint]string{
		38: "Tim Arland",
		39: "Kym Stringer",
		13: "Alex Hayter",
		99: "Rob Lee",
		12: "Gary Loughhead",
		55: "Kevin Bartlett",
		69: "Peter Phillis",
	}

	f, err := os.Create("results.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		f.Close()
	}()
	err = Print(result, skippers, f)
	if err != nil {
		log.Fatalln(err)
	}
}
