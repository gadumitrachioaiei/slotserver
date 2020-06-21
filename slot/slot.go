package slot

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// reels represens the reels of the slot machine
var reels = [][]string{
	{"Scatter", "Sym9", "Sym2", "Sym2", "Sym8"},
	{"Sym9", "Sym3", "Sym6", "Sym7", "Scatter"},
	{"Sym2", "Sym1", "Sym5", "Wild", "Sym1"},
	{"Sym4", "Sym4", "Scatter", "Scatter", "Sym2"},
	{"Sym8", "Sym7", "Sym7", "Sym6", "Sym7"},
	{"Sym5", "Sym9", "Sym9", "Sym8", "Sym4"},
	{"Sym7", "Sym2", "Sym6", "Sym7", "Sym6"},
	{"Sym9", "Sym6", "Sym2", "Sym4", "Sym8"},
	{"Sym4", "Sym8", "Sym4", "Sym1", "Sym3"},
	{"Sym6", "Sym1", "Sym8", "Sym5", "Sym7"},
	{"Sym3", "Sym4", "Sym1", "Sym8", "Sym4"},
	{"Sym8", "Sym9", "Sym3", "Sym9", "Sym2"},
	{"Sym5", "Sym2", "Sym6", "Sym4", "Sym6"},
	{"Sym9", "Wild", "Sym9", "Sym7", "Sym1"},
	{"Sym1", "Sym6", "Sym7", "Sym6", "Sym9"},
	{"Sym3", "Sym5", "Sym4", "Sym2", "Sym5"},
	{"Sym6", "Sym7", "Sym5", "Sym9", "Sym4"},
	{"Sym7", "Sym8", "Sym8", "Sym8", "Sym2"},
	{"Sym5", "Sym4", "Sym9", "Sym3", "Wild"},
	{"Wild", "Sym3", "Sym3", "Sym4", "Sym6"},
	{"Sym8", "Scatter", "Sym2", "Sym7", "Sym3"},
	{"Sym9", "Sym9", "Sym4", "Sym5", "Sym9"},
	{"Sym2", "Sym6", "Sym8", "Sym6", "Sym5"},
	{"Sym7", "Sym7", "Sym7", "Sym3", "Sym2"},
	{"Sym5", "Sym8", "Sym5", "Sym8", "Sym8"},
	{"Scatter", "Sym5", "Wild", "Sym9", "Sym6"},
	{"Sym6", "Sym3", "Sym3", "Sym5", "Sym1"},
	{"Sym8", "Sym9", "Sym8", "Sym2", "Sym9"},
	{"Sym4", "Sym1", "Sym6", "Sym4", "Sym4"},
	{"Sym3", "Sym2", "Sym7", "Sym1", "Sym5"},
	{"Sym1", "Sym7", "Sym9", "Sym9", "Sym7"},
	{"Sym6", "Sym8", "Sym1", "Sym8", "Sym3"},
}

// payLines represent the pay lines of the slot machine
var payLines = [][]int{
	{1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0},
	{2, 2, 2, 2, 2},
	{0, 1, 2, 1, 0},
	{2, 1, 0, 1, 2},
	{1, 0, 0, 0, 1},
	{1, 2, 2, 2, 1},
	{0, 0, 1, 2, 2},
	{2, 2, 1, 0, 0},
	{1, 0, 1, 2, 1},
	{1, 2, 1, 0, 1},
	{0, 1, 1, 1, 0},
	{2, 1, 1, 1, 2},
	{0, 1, 0, 1, 0},
	{2, 1, 2, 1, 2},
	{1, 1, 0, 1, 1},
	{1, 1, 2, 1, 1},
	{0, 0, 2, 0, 0},
	{2, 2, 0, 2, 2},
	{0, 2, 2, 2, 0},
}

// WinRule encodes a winning combination of symbols that must be found in a pay line
type WinRule struct {
	symbol       string
	consecutives int
	prize        int
}

func (p WinRule) String() string {
	return fmt.Sprintf("{\"%s\", %d, %d},", p.symbol, p.consecutives, p.prize)
}

// payTable is the table pay for all symbols except Scatter
// e.g.: 5 consecutives Wild symbols will win you 5000 chips
var payTable = []WinRule{
	{"Wild", 5, 5000},
	{"Sym1", 5, 1000},
	{"Wild", 4, 500},
	{"Sym2", 5, 500},
	{"Sym3", 5, 300},
	{"Sym1", 4, 200},
	{"Sym5", 5, 200},
	{"Sym4", 5, 200},
	{"Sym2", 4, 150},
	{"Sym3", 4, 100},
	{"Sym6", 5, 100},
	{"Sym7", 5, 100},
	{"Sym5", 4, 75},
	{"Sym4", 4, 75},
	{"Wild", 3, 50},
	{"Sym7", 4, 50},
	{"Sym6", 4, 50},
	{"Sym9", 5, 50},
	{"Sym8", 5, 50},
	{"Sym1", 3, 40},
	{"Sym2", 3, 30},
	{"Sym3", 3, 25},
	{"Sym8", 4, 25},
	{"Sym9", 4, 25},
	{"Sym5", 3, 20},
	{"Sym4", 3, 20},
	{"Sym7", 3, 15},
	{"Sym6", 3, 15},
	{"Sym8", 3, 10},
	{"Sym9", 3, 10},
	{"Wild", 2, 5},
	{"Sym1", 2, 3},
	{"Sym2", 2, 2},
	{"Sym3", 2, 2},
}

// scatterPayTable represents prizes for the Scatter symbols
// map between number of Scatter symbols and prize
var scatterPayTable = map[int]int{
	5: 100,
	4: 25,
	3: 5,
}

// SpinTypeString represents the type of spin
type SpinTypeInt int32

const (
	SpinTypeIntMain SpinTypeInt = iota + 1
	SpinTypeIntFree
)

// SpinTypeString represents the type of spin, as a string
type SpinType string

const (
	SpinTypeMain SpinType = "MAIN"
	SpinTypeFree SpinType = "FREE"
)

func spinTypeToInt(typ SpinType) SpinTypeInt {
	switch typ {
	case SpinTypeMain:
		return SpinTypeIntMain
	case SpinTypeFree:
		return SpinTypeIntFree
	}
	return 0
}

// Spin represents a spin of reels
type Spin struct {
	TypeInt  SpinTypeInt `json:"-"` // this field should match the Type field, and it is here just for the needs of the rpc service
	Type     SpinType
	Stops    [][]int // stops for this spin
	Win      int     // how much this spin won
	PayLines [][]int // which lines won
}

// Result expresses the result of a bet and it can contain multiple spins: the main one and maybe free spins
type Result struct {
	Spins []Spin
	Win   int // sum of wins from all spins
	Chips int // chips after all spins
	Bet   int // wager of a bet
}

// Machine is an Atkin Diet slot machine
type Machine struct {
	c                   int            // number of pay lines
	freeSpins           int            // number of free spins in case Scatter happens
	freeSpinsMultiplier int            // win multiplier in case of free spins winnings
	testSpin            func() [][]int // set only from tests
	payLinesTree        *N             // paylines as a tree
}

// NewMachine returns a machine with default settings
func NewMachine() *Machine {
	m := Machine{c: len(payLines), freeSpins: 10, freeSpinsMultiplier: 3}
	m.payLinesTree = tree(payLines)
	return &m
}

// Bet bets the wager and returns the result
func (m *Machine) Bet(chips, wager int) (*Result, error) {
	if wager > chips || chips <= 0 || wager <= 0 {
		return nil, errors.New("incorrect chips and wager amount")
	}
	r := m.bet(wager, false)
	r.Chips = chips - wager + r.Win
	r.Bet = wager
	return r, nil
}

// bet bets the wager and returns the result
func (m *Machine) bet(wager int, freeSpinMode bool) *Result {
	// spin the reels
	var stops [][]int
	if m.testSpin == nil {
		stops = m.spin()
	} else {
		stops = m.testSpin()
	}
	wins, lines := m.baseWinnings(stops)
	scatterWins := scatterWinnings(stops)
	// current spin result
	var result Result
	var typ SpinType = SpinTypeMain
	if freeSpinMode {
		typ = SpinTypeFree
	}
	wins = (wins + scatterWins) * (wager / m.c)
	if freeSpinMode {
		wins *= m.freeSpinsMultiplier
	}
	result = Result{
		Spins: []Spin{{TypeInt: spinTypeToInt(typ), Type: typ, Win: wins, Stops: stops, PayLines: lines}},
		Win:   wins,
	}
	if scatterWins == 0 {
		return &result
	}
	// if we have scatter wins we provide free spins and add their results to the base one
	for i := 0; i < m.freeSpins; i++ {
		r := m.bet(wager, true)
		result.Spins = append(result.Spins, r.Spins...)
		result.Win += r.Win
	}
	return &result
}

// baseWinnings returns the total wins and the winning pay lines for the stops
// it does not take into consideration scatter symbols
func (m *Machine) baseWinnings(stops [][]int) (int, [][]int) {
	wins := 0
	var lines [][]int
	for _, pl := range payLines {
		// what are the symbols in each payLine
		symbols := []int{stops[pl[0]][0], stops[pl[1]][1], stops[pl[2]][2], stops[pl[3]][3], stops[pl[4]][4]}
		_, win := matchLineSymbols(symbols)
		wins += win
		if win > 0 {
			lines = append(lines, pl)
		}
	}
	return wins, lines
}

// scatterWinnings returns winnings from Scatter symbols only for the stops
func scatterWinnings(stops [][]int) int {
	scatters := 0
	for i := range stops {
		for j := range stops[i] {
			if reels[i][j] == "Scatter" {
				scatters++
			}
		}
	}
	return scatterPayTable[scatters]
}

// matchLineSymbols matches symbols against the payTable to check for any winnings
// returns the winning payTable line index and the win
// symbols are the symbols in a spin's stops, i.e.: indexes in the reels columns
func matchLineSymbols(symbols []int) (int, int) {
	win := 0
	var i int
	var rule WinRule
	for i, rule = range payTable {
		match := true
		for j, symbol := range symbols[0:rule.consecutives] {
			if !(reels[symbol][j] == rule.symbol || reels[symbol][j] == "Wild") {
				match = false
				break
			}
		}
		if match {
			win += rule.prize
			break
		}
	}
	return i, win
}

// spin spins the reels and returns the result as a matrix of symbols
// ( values are indexes into reels rows )
func (m *Machine) spin() [][]int {
	// generate 5 random numbers between 0 and 31
	// ( because the slot has 5 reels each with 32 possible values )
	r := rand.New(rand.NewSource(time.Now().Unix()))

	r1, r2, r3, r4, r5 := r.Intn(32), r.Intn(32), r.Intn(32), r.Intn(32), r.Intn(32)
	resultn := make([][]int, 3)
	resultn[1] = []int{r1, r2, r3, r4, r5}
	resultn[0] = make([]int, 5)
	resultn[2] = make([]int, 5)
	for i, result := range resultn[1] {
		firstStop := result - 1
		if result == 0 {
			firstStop = 31
		}
		thirdStop := result + 1
		if result == 31 {
			thirdStop = 0
		}
		resultn[0][i] = firstStop
		resultn[2][i] = thirdStop
	}
	return resultn
}
