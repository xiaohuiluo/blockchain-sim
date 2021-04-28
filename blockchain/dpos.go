package blockchain

import (
	"errors"
	"sort"
)

var nodeVoteMap map[string]int

func InitVoteMap() {
	nodeVoteMap = make(map[string]int)
}

func NodeVoteMap() map[string]int {
	return nodeVoteMap
}

func Vote(id string, ticket int) {
	nodeVoteMap[id] = ticket
}

// PickWinner 根据投票数量选择生成区块的节点
func PickWinnerWithDpos() (bp string, err error) {

	if len(nodeVoteMap) < 1 {
		err = errors.New("Error: failed to pick winner, node vote map is empty")
		return
	}

	// 选择BlockProducer
	ticketList := make([]int, len(nodeVoteMap))
	for _, ticket := range nodeVoteMap {
		ticketList = append(ticketList, ticket)
	}

	sort.Slice(ticketList, func(i, j int) bool {
		return ticketList[i] > ticketList[j]
	})

	// 前一半作为producer
	ticketList = ticketList[0 : len(ticketList)/2]

	for k, v := range nodeVoteMap {
		if v > ticketList[len(ticketList)-1] {
			bp = k
			err = nil
		}
	}

	return
}
