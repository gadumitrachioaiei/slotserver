package slot

// baseWinningsTree returns the won prize and the winning paylines
// it uses pre-order tree traversal to match the paylines with the rules
func (m *Machine) baseWinningsTree(stops [][]int) (int, [][]int) {
	wt := WinningsTree{}
	for _, n := range m.payLinesTree.Ch {
		pathNodes := make([]int, 0, 5)
		// at beginning all rules are potential matches
		wt.win(n, pathNodes, stops, payTable)
	}
	return wt.prize, wt.winnings
}

// WinningsTree holds the winning paylines and the prize
// the data is set from the tree traversal matching rules
type WinningsTree struct {
	winnings [][]int
	prize    int
}

// win sets what paylines have won
// win determines all paylines that pass through this node which won
// rules are rules that have been satistfied before this node
// pathNodes is the tree branch traversed so far ( node values ), not including current node
// stops represent the reels according to current spin
func (wt *WinningsTree) win(n *N, pathNodes []int, stops [][]int, rules []WinRule) {
	// which of the rules are still satisfied by this node
	var nrules []WinRule
	depth := len(pathNodes)
	pathNodes = append(pathNodes, n.V)
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		if rule.consecutives <= depth {
			// this rule was already satisfied
			nrules = append(nrules, rule)
			continue
		}
		switch reels[stops[n.V][depth]][depth] {
		case rule.symbol, "Wild":
			// this rule is still valid
			nrules = append(nrules, rule)
		}
	}
	// no rules is matched so far, nothing else to check
	if len(nrules) == 0 {
		return
	}
	// if we have children, check if they match rules
	if len(n.Ch) > 0 {
		for _, child := range n.Ch {
			wt.win(child, pathNodes, stops, nrules)
		}
		return
	}
	// if we are a leaf node, we append the winning payline
	prize := 0
	for i := range nrules {
		if nrules[i].prize > prize {
			prize = nrules[i].prize
		}
	}
	wt.winnings = append(wt.winnings, pathNodes)
	wt.prize += prize
}

// N is a node in the paylines tree
type N struct {
	V  int
	Ch []*N
}

// func (n *N) String() string {
// 	data, err := json.MarshalIndent(n, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return string(data)
// }

// tree transforms the paylines into a tree, each payline becoming a branch of the tree
// tree has an artificial root, where all paylines are attached
func tree(lines [][]int) *N {
	root := &N{V: -1}
	for _, l := range lines {
		node := root
		for _, e := range l {
			found := false
			for _, child := range node.Ch {
				if child.V == e {
					found = true
					node = child
					break
				}
			}
			if !found {
				nn := &N{V: e}
				node.Ch = append(node.Ch, nn)
				node = nn
			}
		}
	}
	return root
}
