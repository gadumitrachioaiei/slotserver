package docs

// swagger:route POST /api/machines/atkins-diet/spins casino createBet
// You place a bet and spin an atkins-diet slot machine.
//
// responses:
//   200: createBet
//   default: errorResponse

// Bet created
// swagger:response createBet
type CreateBetResponseWrapper struct {
	// in:body
	Body Response
}

// swagger:parameters createBet
type CreateBetParamsWrapper struct {
	// User and bet data.
	//
	// in: body
	Body Request
}

type Response struct {
	*Result
	JWT Request
}

type Result struct {
	// spins of the slot machine
	Spins []Spin
	// sum of wins from all spins
	Win int
	// chips after all spins
	Chips int
	// wager of the bet
	Bet int
}

// swagger:enum SpinType
type SpinType string

const (
	SpinTypeMain SpinType = "MAIN"
	SpinTypeFree SpinType = "FREE"
)

// Spin represents a spin of reels
type Spin struct {
	Type SpinType
	// stops for this spin
	Stops [][]int
	// how much this spin won
	Win int
	// which lines won
	PayLines [][]int
}

type Request struct {
	// user id
	// required: true
	UID string
	// chips balance
	// required: true
	Chips int
	// amount of chips bet
	// required: true
	Bet int
}

// Error response
// swagger:response errorResponse
type BadRequestWrapper struct {
	// in:body
	Body struct {
		// error message
		// example: incorrect chips and wager amount
		Message string
		// status code
		// example: 400
		Code int
	}
}
