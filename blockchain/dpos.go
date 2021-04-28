package blockchain

import (
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
func PickWinner() (bp string) {
	// 选择BlockProducer
	voteMap := NodeVoteMap()
	ticketList := make([]int, len(voteMap))
	for _, ticket := range voteMap {
		ticketList = append(ticketList, ticket)
	}

	sort.Slice(ticketList, func(i, j int) bool {
		return ticketList[i] > ticketList[j]
	})

	// 前一半作为producer
	ticketList = ticketList[0 : len(ticketList)/2]

	for k, v := range voteMap {
		if v > ticketList[len(ticketList)-1] {
			bp = k
		}
	}

	return
}
