package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

type Race struct {
	Finishers     []uint // Sailnumbers
	DidNotFinish  []uint // Sailnumbers
	AveragePoints []uint // Sailnumbers
}

type Result struct {
	SailNumber uint
	Rank       uint
	Tot        float32
	Scores     []Score
}

type Score struct {
	Score float32
	Drop  bool
}

func Calculate(races []Race) ([]Result, error) {
	if races == nil {
		return nil, errors.New("races may not be nil")
	}
	if len(races) <= 0 {
		return nil, errors.New("must have at least one race")
	}

	numRaces := len(races)
	numBoats := len(races[0].Finishers) + len(races[0].DidNotFinish) + len(races[0].AveragePoints)

	for _, race := range races {
		if len(race.Finishers)+len(race.DidNotFinish)+len(race.AveragePoints) != numBoats {
			return nil, errors.New("all races must have the same number of boats")
		}
		for _, finishedSailNumber := range race.Finishers {
			for _, dnfSailNumber := range race.DidNotFinish {
				if finishedSailNumber == dnfSailNumber {
					return nil, fmt.Errorf("a DNF sailnumber (%d) cannot also be a finished sailnumber (%d)", finishedSailNumber, dnfSailNumber)
				}
			}
			for _, averagePointsSailNumber := range race.AveragePoints {
				if finishedSailNumber == averagePointsSailNumber {
					return nil, fmt.Errorf("an average points sailnumber (%d) cannot also be a finished sailnumber (%d)", finishedSailNumber, averagePointsSailNumber)
				}
			}
		}
	}

	// Map sailnumbers to Boats and create the Scores array in each Boat.
	boats := make(map[uint]Boat, numBoats)
	for _, sailNumber := range races[0].Finishers {
		var boat Boat
		boat.Scores = make([]Score, numRaces)
		boats[sailNumber] = boat
	}
	for _, sailNumber := range races[0].DidNotFinish {
		var boat Boat
		boat.Scores = make([]Score, numRaces)
		boats[sailNumber] = boat
	}
	for _, sailNumber := range races[0].AveragePoints {
		var boat Boat
		boat.Scores = make([]Score, numRaces)
		boats[sailNumber] = boat
	}

	// Iterate over the races and assign scores to each Finished or DNF Boat.
	for i, race := range races {
		for pos, sailNumber := range race.Finishers {
			boat := boats[sailNumber]
			boat.Scores[i].Score = float32(pos + 1)
			boats[sailNumber] = boat
		}
		for _, sailNumber := range race.DidNotFinish {
			boat := boats[sailNumber]
			boat.Scores[i].Score = float32(numBoats + 1)
			boats[sailNumber] = boat
		}
	}

	// Calculate the average points for each boat (dropped scores ARE included)...
	for sailNumber, boat := range boats {
		numOfNonAveragePoinRaces := 0
		for _, score := range boat.Scores {
			if score.Score > 0 {
				boat.AveragePoints += score.Score
				numOfNonAveragePoinRaces++
			}
		}
		if numOfNonAveragePoinRaces > 0 {
			boat.AveragePoints /= float32(numOfNonAveragePoinRaces)
		}
		boats[sailNumber] = boat
	}
	// ... then insert this value in each race for each boat tagged as AveragePoints.
	for i, race := range races {
		for _, sailNumber := range race.AveragePoints {
			boat := boats[sailNumber]
			if boat.AveragePoints == 0 {
				return nil, fmt.Errorf("sailnumber %d has average points in race %d, but has no actual races from which the points can be calculated", sailNumber, i+1)
			}
			boat.Scores[i].Score = boat.AveragePoints
			boats[sailNumber] = boat
		}
	}

	// Calculate the number of drops. One drop per completed eight races plus one for from race 4.
	drops := numRaces / 8
	if numRaces >= 4 {
		drops++
	}

	// Mark the worst scores as drops. May mark 'drops' number of scores as such.
	for sailNumber, boat := range boats {
		for dropped := 0; dropped < drops; {
			// Find the worst non-discard score.
			theWorst := -1
			for i, score := range boat.Scores {
				// First the first non-discard score if we haven't already found one or the next non-discard that is worse
				// then the currently worst score.
				if !score.Drop && (theWorst == -1 || boat.Scores[theWorst].Score < score.Score) {
					theWorst = i
				}
			}
			boat.Scores[theWorst].Drop = true
			dropped++
		}
		boats[sailNumber] = boat
	}

	// Calculate the total score for each boat. Ignore dropped races
	for sailNumber, boat := range boats {
		for _, score := range boat.Scores {
			if !score.Drop {
				boat.Tot += score.Score
			}
		}
		boats[sailNumber] = boat
	}

	// Create a slice of sail numbers, then sort the sailnumbers slice such that the first place boat is first, etc.
	sortedSailNumbers := make([]uint, 0, numBoats)
	for sailNumber := range boats {
		sortedSailNumbers = append(sortedSailNumbers, sailNumber)
	}
	sort.SliceStable(sortedSailNumbers, func(i, j int) bool {
		boatI := boats[sortedSailNumbers[i]]
		boatJ := boats[sortedSailNumbers[j]]

		if boatI.Tot < boatJ.Tot {
			return true
		} else if boatI.Tot > boatJ.Tot {
			return false
		}

		// The total scores are the same. Who has the most wins, then (if still undecided), who has
		// the most seconds, then (if still undecided) etc... Discards are not included here (RRS A8.1).
		scoresNoDiscardsI := sortedLowToHighWithoutDiscards(boatI.Scores)
		scoresNoDiscardsJ := sortedLowToHighWithoutDiscards(boatJ.Scores)
		for a := 0; a < len(scoresNoDiscardsI) && a < len(scoresNoDiscardsJ); a++ {
			if scoresNoDiscardsI[a] < scoresNoDiscardsJ[a] {
				return true
			} else if scoresNoDiscardsI[a] > scoresNoDiscardsJ[a] {
				return false
			}
		}

		// The number of firsts, seconds, etc are the same. Who did best in the most recent race?
		// Discards ARE included here (RRS A8.2).
		for race := numRaces - 1; race >= 0; race-- {
			if boatI.Scores[race].Score < boatJ.Scores[i].Score {
				return true
			} else if boatI.Scores[i].Score > boatJ.Scores[i].Score {
				return false
			}

		}

		log.Printf("Unable to separate boats %d and %d", sortedSailNumbers[i], sortedSailNumbers[j])
		return false
	})

	// Populate the results in increasing rank order.
	results := make([]Result, numBoats)
	for rank, sailNumber := range sortedSailNumbers {
		result := Result{
			SailNumber: sailNumber,
			Rank:       uint(rank + 1),
			Tot:        boats[sailNumber].Tot,
			Scores:     make([]Score, numRaces),
		}
		copy(result.Scores, boats[sailNumber].Scores)
		results[rank] = result
	}

	return results, nil
}

func sortedLowToHighWithoutDiscards(in []Score) []float32 {
	out := make([]float32, 0)

	for _, score := range in {
		if !score.Drop {
			out = append(out, score.Score)
		}
	}

	sort.SliceStable(out, func(a, b int) bool {
		return out[a] < out[b]
	})

	return out
}

type Boat struct {
	Tot           float32
	Scores        []Score
	AveragePoints float32
}
