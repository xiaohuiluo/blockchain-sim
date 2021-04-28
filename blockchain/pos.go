package blockchain

import "errors"

var nodeTokenAgeMap map[string]int64

func InitTokenAgeMap() {
	nodeTokenAgeMap = make(map[string]int64)
}

func NodeTokenAgeMap() map[string]int64 {
	return nodeTokenAgeMap
}

func SetTokenAge(id string, TokenAge int64) {
	nodeTokenAgeMap[id] = TokenAge
}

func GetTokenAge(id string) int64 {
	return nodeTokenAgeMap[id]
}

// PickWinnerPos 根据代币和持有时间选择生成区块的节点
func PickWinnerWithPos() (bp string, err error) {
	if len(nodeTokenAgeMap) < 1 {
		err = errors.New("Error: failed to pick winner, node tokenAge map is empty")
		return
	}

	// 选择BlockProducer
	var maxTokenAge int64
	maxTokenAge = 0
	for nodeId, tokenAge := range nodeTokenAgeMap {
		if tokenAge > maxTokenAge {
			maxTokenAge = tokenAge
			bp = nodeId
		}
	}

	err = nil
	return
}
