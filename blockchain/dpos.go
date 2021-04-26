package blockchain

import (
	"sort"
)

// BPCount 区块生产者的数量
const BPCount = 5

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

	if len(ticketList) > BPCount {
		ticketList = ticketList[0:BPCount] // 选择前面的5个节点作为Block producer
	}

	for k, v := range voteMap {
		if v > ticketList[len(ticketList)-1] {
			bp = k
		}
	}

	return
}
