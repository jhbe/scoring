package main

import (
	"reflect"
	"testing"
)

func TestCalculate(t *testing.T) {
	cases := []struct {
		in       []Race
		expected []Result
		err      bool
	}{
		{
			// Nil argument.
			in:       nil,
			expected: nil,
			err:      true,
		},
		{
			// No races.
			in:       []Race{},
			expected: nil,
			err:      true,
		},
		{
			// Races with different number of boats.
			in: []Race{
				{Finishers: []uint{36, 34, 30}},
				{Finishers: []uint{30, 34}},
				{Finishers: []uint{30, 36, 34}},
				{Finishers: []uint{30, 34, 36}},
			},
			expected: nil,
			err:      true,
		},
		{
			// Single boat, single race.
			in:       []Race{{Finishers: []uint{10}}},
			expected: []Result{{SailNumber: 10, Rank: 1, Tot: 1, Scores: []Score{{1, false}}}},
		},
		{
			// Boat finished and was tagged DNF.
			in:       []Race{{Finishers: []uint{10}, DidNotFinish: []uint{10}}},
			expected: nil,
			err:      true,
		},
		{
			// Boat awarded average points in a single race event.
			in:       []Race{{AveragePoints: []uint{10}}},
			expected: nil,
			err:      true,
		},
		{
			// Single boat, single race, DNF.
			in:       []Race{{Finishers: []uint{}, DidNotFinish: []uint{10}}},
			expected: []Result{{SailNumber: 10, Rank: 1, Tot: 2, Scores: []Score{{2, false}}}},
		},
		{
			// Single boat, two races.
			in: []Race{
				{Finishers: []uint{10}},
				{Finishers: []uint{10}},
			},
			expected: []Result{{SailNumber: 10, Rank: 1, Tot: 2, Scores: []Score{{1, false}, {1, false}}}},
		},
		{
			// Regular race with one DNF. 11 and 18 have equal scores, but 11 has a win.
			in: []Race{
				{Finishers: []uint{10, 11, 23, 18}},
				{Finishers: []uint{10, 23, 18}, DidNotFinish: []uint{11}},
			},
			expected: []Result{
				{SailNumber: 10, Rank: 1, Tot: 2, Scores: []Score{{1, false}, {1, false}}},
				{SailNumber: 23, Rank: 2, Tot: 5, Scores: []Score{{3, false}, {2, false}}},
				{SailNumber: 11, Rank: 3, Tot: 7, Scores: []Score{{2, false}, {5, false}}},
				{SailNumber: 18, Rank: 4, Tot: 7, Scores: []Score{{4, false}, {3, false}}},
			},
		},
		{
			// Two boats, one race.
			in: []Race{
				{Finishers: []uint{36, 34}},
			},
			expected: []Result{
				{SailNumber: 36, Rank: 1, Tot: 1, Scores: []Score{{1, false}}},
				{SailNumber: 34, Rank: 2, Tot: 2, Scores: []Score{{2, false}}},
			},
		},
		{
			// Two boats, two races, equal scores. 34 did best in the last race.
			in: []Race{
				{Finishers: []uint{36, 34}},
				{Finishers: []uint{34, 36}},
			},
			expected: []Result{
				{SailNumber: 34, Rank: 1, Tot: 3, Scores: []Score{{2, false}, {1, false}}},
				{SailNumber: 36, Rank: 2, Tot: 3, Scores: []Score{{1, false}, {2, false}}},
			},
		},
		{
			// Two boats, three races, one Average Points.
			in: []Race{
				{Finishers: []uint{36, 34}},
				{Finishers: []uint{34, 36}},
				{Finishers: []uint{34}, AveragePoints: []uint{36}},
			},
			expected: []Result{
				{SailNumber: 34, Rank: 1, Tot: 4, Scores: []Score{{2, false}, {1, false}, {1, false}}},
				{SailNumber: 36, Rank: 2, Tot: 4.5, Scores: []Score{{1, false}, {2, false}, {1.5, false}}},
			},
		},
		{
			// 34 and 36 have equal scores, but 36 has a win. One drop.
			in: []Race{
				{Finishers: []uint{36, 34, 30}},
				{Finishers: []uint{30, 34, 36}},
				{Finishers: []uint{30, 36, 34}},
				{Finishers: []uint{30, 34, 36}},
			},
			expected: []Result{
				{SailNumber: 30, Rank: 1, Tot: 3, Scores: []Score{{3, true}, {1, false}, {1, false}, {1, false}}},
				{SailNumber: 36, Rank: 2, Tot: 6, Scores: []Score{{1, false}, {3, true}, {2, false}, {3, false}}},
				{SailNumber: 34, Rank: 3, Tot: 6, Scores: []Score{{2, false}, {2, false}, {3, true}, {2, false}}},
			},
		},
		{
			in: []Race{
				{Finishers: []uint{39, 55, 38, 12, 13, 99, 69}},
				{Finishers: []uint{13, 99, 38, 12, 39, 55, 69}},
				{Finishers: []uint{13, 38, 99, 12, 39, 55, 69}},
				{Finishers: []uint{38, 39, 99, 13, 12, 69, 55}},
				{Finishers: []uint{38, 39, 13, 12, 99, 55, 69}},
				{Finishers: []uint{55, 38, 99, 13, 39, 12, 69}},
				{Finishers: []uint{39, 99, 69, 12}, DidNotFinish: []uint{13, 55}, AveragePoints: []uint{38}},
				{Finishers: []uint{39, 12, 69, 99}, DidNotFinish: []uint{38, 55}, AveragePoints: []uint{13}},
			},
			expected: []Result{
				{SailNumber: 38, Rank: 1, Tot: 11.857143, Scores: []Score{{3, true}, {3, false}, {2, false}, {1, false}, {1, false}, {2, false}, {2.857143, false}, {8, true}}},
				{SailNumber: 39, Rank: 2, Tot: 12, Scores: []Score{{1, false}, {5, true}, {5, true}, {2, false}, {2, false}, {5, false}, {1, false}, {1, false}}},
				{SailNumber: 13, Rank: 3, Tot: 16.7142856, Scores: []Score{{5, true}, {1, false}, {1, false}, {4, false}, {3, false}, {4, false}, {8, true}, {3.7142856, false}}},
				{SailNumber: 99, Rank: 4, Tot: 17, Scores: []Score{{6, true}, {2, false}, {3, false}, {3, false}, {5, true}, {3, false}, {2, false}, {4, false}}},
				{SailNumber: 12, Rank: 5, Tot: 22, Scores: []Score{{4, false}, {4, false}, {4, false}, {5, true}, {4, false}, {6, true}, {4, false}, {2, false}}},
				{SailNumber: 55, Rank: 6, Tot: 28, Scores: []Score{{2, false}, {6, false}, {6, false}, {7, false}, {6, false}, {1, false}, {8, true}, {8, true}}},
				{SailNumber: 69, Rank: 7, Tot: 33, Scores: []Score{{7, true}, {7, true}, {7, false}, {6, false}, {7, false}, {7, false}, {3, false}, {3, false}}},
			},
		},
	}

	for i, c := range cases {
		t.Logf("Case %d...", i)
		out, err := Calculate(c.in)
		if !c.err && err != nil {
			t.Errorf("Case %d: Got an error \"%s\", but didn't expect one", i, err)
		} else if c.err && err == nil {
			t.Errorf("Case %d: Got no error, but expect one", i)
		} else if err == nil && !reflect.DeepEqual(out, c.expected) {
			t.Errorf("Case %d: Expected %v, got %v", i, c.expected, out)
		}
	}
}
