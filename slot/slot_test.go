package slot

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestMachineManual(t *testing.T) {
	m := NewMachine()
	result, err := m.Bet(1000, 100)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result.debug())
	fmt.Println(string(data))
}

type TestCase struct {
	name     string
	chips    int
	wager    int
	testSpin func() [][]int
	result   *Result
}

var testCases = []TestCase{
	{
		name:  "test1",
		chips: 1000,
		wager: 100,
		testSpin: func() [][]int {
			return [][]int{
				{25, 30, 23, 11, 5},
				{26, 31, 24, 12, 6},
				{27, 0, 25, 13, 7},
			}
		},
		result: &Result{
			Win:   50,
			Chips: 950,
			Bet:   100,
			Spins: []Spin{
				{
					Type: "main",
					Win:  50,
					PayLines: [][]int{
						{2, 1, 2, 1, 2},
					},
					Stops: [][]int{
						{25, 30, 23, 11, 5},
						{26, 31, 24, 12, 6},
						{27, 0, 25, 13, 7},
					},
				},
			},
		},
	},
	{
		name:  "test2",
		chips: 1000,
		wager: 100,
		testSpin: func() [][]int {
			return [][]int{
				{25, 12, 1, 29, 25},
				{26, 13, 2, 30, 26},
				{27, 14, 3, 31, 27},
			}
		},
		result: &Result{
			Win:   75,
			Chips: 975,
			Bet:   100,
			Spins: []Spin{
				{
					Type: "main",
					Win:  75,
					PayLines: [][]int{
						{1, 1, 0, 1, 1},
					},
					Stops: [][]int{
						{25, 12, 1, 29, 25},
						{26, 13, 2, 30, 26},
						{27, 14, 3, 31, 27},
					},
				},
			},
		},
	},
}

func TestMachineWithSpin(t *testing.T) {
	m := NewMachine()
	for _, tc := range testCases {
		m.testSpin = tc.testSpin
		result, err := m.Bet(tc.chips, tc.wager)
		if err != nil {
			t.Fatalf("%s %v", tc.name, err)
		}
		if !reflect.DeepEqual(result, tc.result) {
			t.Errorf("%s got result: \n%#v\n, expected: \n%#v\n", tc.name, result, tc.result)
		}
	}
}

var prize int
var lines [][]int

// BenchmarkBaseWinnings-8   	  200000	      8609 ns/op	      36 B/op	       0 allocs/op
func BenchmarkBaseWinnings(b *testing.B) {
	m := NewMachine()
	var w int
	var l [][]int
	for i := 0; i < b.N; i++ {
		// generate random stops
		b.StopTimer()
		stops := m.spin()
		b.StartTimer()
		// call baseWinnings on them
		w, l = m.baseWinnings(stops)
	}
	prize = w
	lines = l
}

// BenchmarkBaseWinningsTree-8   	  300000	      5424 ns/op	    1898 B/op	      32 allocs/op
func BenchmarkBaseWinningsTree(b *testing.B) {
	m := NewMachine()
	var w int
	var l [][]int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// generate random stops
		b.StopTimer()
		stops := m.spin()
		b.StartTimer()
		// call baseWinnings on them
		w, l = m.baseWinningsTree(stops)
	}
	prize = w
	lines = l
}

// debug returns information for this result so you can understand if it is correct
func (r *Result) debug() string {
	var buf strings.Builder
	for i, spin := range r.Spins {
		buf.WriteString(fmt.Sprintf("%s spin: %d\n", spin.Type, i))
		stops := spin.Stops
		win := spin.Win
		// display the reels for the stops
		buf.WriteString("\treels\n")
		for i := range stops {
			buf.WriteString("\t")
			for j := range stops[i] {
				buf.WriteString(fmt.Sprintf("%s ", reels[stops[i][j]][j]))
			}
			buf.WriteString("\n")
		}
		buf.WriteString("\tpaylines\n")
		for j, pl := range payLines {
			// the symbols in each payLine
			symbols := []int{stops[pl[0]][0], stops[pl[1]][1], stops[pl[2]][2], stops[pl[3]][3], stops[pl[4]][4]}
			buf.WriteString(fmt.Sprintf("\tpay line: %d", j))
			for j, symbol := range symbols {
				buf.WriteString(fmt.Sprintf(" %s", reels[symbol][j]))
			}
			if win > 0 {
				payTableLine, lineWin := matchLineSymbols(symbols)
				if lineWin > 0 {
					buf.WriteString(fmt.Sprintf(", win %d for pay table line: %v", lineWin, payTable[payTableLine]))
				}
			}
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("\tscatter wins: %d\n", scatterWinnings(stops)))
	}
	return buf.String()
}
