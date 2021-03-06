This is an implementation of atkins-diet slot machine:
https://wizardofodds.com/play/slots/atkins-diet/atkins-diet.pdf
You can find clearer rules and discussion here:
http://giorasimchoni.com/2017/05/06/2017-05-06-don-t-drink-and-gamble/

The purpose of this code is to understand the game and the implementation. There are very few tests. There are improvements that can be made performance wise.

Just run the server and try requests, e.g.:
```curl -v -X POST --data '{"uid": "xyz", "chips": 1000, "bet":100}' "http://127.0.0.1:9999/api/machines/atkins-diet/spins"```
```grpcurl -plaintext -import-path third_party/proto/ -import-path .  -proto  proto/bet/v1/bet.proto -d '{"uid": "xyz", "chips": 1000, "bet":100}' localhost:8080 bet.v1.SlotMachineService/CreateBet```
Example response:
```{
    "Spins": [
        {
            "Type": "main",
            "Stops": [
                [11,23,14,16,30],
                [12,24,15,17,31],
                [13,25,16,18,0]
            ],
            "Win": 100,
            "PayLines": [
                [1,2,2,2,1],
            ]
        }
    ],
    "Win": 100,
    "Chips": 1000,
    "Bet": 100,
    "JWT": {
        "UID": "xyz",
        "Chips": 1000,
        "Bet": 100
    }
}
```
All indexes start with 0. "Stops" represents the result grid, "PayLines" represents the winning pay lines and "Type" can be main or free. There is only a main spin, and all others are free, returned in the spinning order. To find more about the result you can run TestMachineManual or TestMachineWithSpin and print the result of the Result.debug method.
Example:
```
go test -v github.com/gadumitrachioaiei/slotserver/slot -run TestMachineManual

=== RUN   TestMachineManual
main spin: 0
        reels
        Sym6 Sym3 Sym5 Sym3 Sym7
        Sym3 Sym1 Sym8 Sym8 Sym4
        Sym8 Sym4 Sym9 Sym9 Sym2
        paylines
        pay line: 0 Sym3 Sym1 Sym8 Sym8 Sym4
        pay line: 1 Sym6 Sym3 Sym5 Sym3 Sym7
        pay line: 2 Sym8 Sym4 Sym9 Sym9 Sym2
        pay line: 3 Sym6 Sym1 Sym9 Sym8 Sym7
        pay line: 4 Sym8 Sym1 Sym5 Sym8 Sym2
        pay line: 5 Sym3 Sym3 Sym5 Sym3 Sym4, win 2 for pay table line: {"Sym3", 2, 2},
        pay line: 6 Sym3 Sym4 Sym9 Sym9 Sym4
        pay line: 7 Sym6 Sym3 Sym8 Sym9 Sym2
        pay line: 8 Sym8 Sym4 Sym8 Sym3 Sym7
        pay line: 9 Sym3 Sym3 Sym8 Sym9 Sym4, win 2 for pay table line: {"Sym3", 2, 2},
        pay line: 10 Sym3 Sym4 Sym8 Sym3 Sym4
        pay line: 11 Sym6 Sym1 Sym8 Sym8 Sym7
        pay line: 12 Sym8 Sym1 Sym8 Sym8 Sym2
        pay line: 13 Sym6 Sym1 Sym5 Sym8 Sym7
        pay line: 14 Sym8 Sym1 Sym9 Sym8 Sym2
        pay line: 15 Sym3 Sym1 Sym5 Sym8 Sym4
        pay line: 16 Sym3 Sym1 Sym9 Sym8 Sym4
        pay line: 17 Sym6 Sym3 Sym9 Sym3 Sym7
        pay line: 18 Sym8 Sym4 Sym5 Sym9 Sym2
        pay line: 19 Sym6 Sym4 Sym9 Sym9 Sym7
        scatter wins: 0

{"Spins":[{"Type":"main","Stops":[[9,1,16,23,9],[10,2,17,24,10],[11,3,18,25,11]],"Win":20,"PayLines":[[1,0,0,0,1],[1,0,1,2,1]]}],"Win":20,"Chips":920,"Bet":100}
```
Short description of the rules:
1. A pay line wins if it contains a pay table entry ( scatter pay table is handled separated ), matched starting with first symbol.
Wild symbol can be replaced with any symbol.
If a pay line contains several pay table entries, than the one with biggest prize wins.
If there are several pay table entries the shortest one wins.
2. For the scatter symbols to give you a win, they can be found anywhere in the grid.
Besides the prize found in the scatter pay table, at least 3 scatter symbols will also give you 10 free spins with triple the normal win.
A free spin can also give you more free spins, unlimited.
3. All pay lines that contain a pay table entry win.
3. All prizes are multiplied by wager / number of pay lines.

Header: :status: 200
Header: grpc-status: 13
Header: grpc-message: placing the bet: Not enough chips
    
    
    