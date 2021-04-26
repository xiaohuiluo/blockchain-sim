/**
Tools package. It's contain some useful tools, just like vote and so on.
This file is created by magic at 2018-9-3
**/
package blockchain

import (
	"errors"

	log "github.com/go-fastlog/fastlog"
)

const (
	// DefaultTicket 节点默认的票数
	DefaultTicket = 10
)

var nodeVoteMap = make(map[string]int)

func NodeVoteMap() map[string]int {
	return nodeVoteMap
}

func InitVote(id string) {
	nodeVoteMap[id] = DefaultTicket
	log.Infof("init vote for node id = %s, ticket = %d", id, DefaultTicket)
}

func Vote(id string, ticket int) error {

	if value, ok := nodeVoteMap[id]; ok {
		nodeVoteMap[id] = value + ticket
		return nil
	}

	log.Errorf("id = %s is not exist", id)
	return errors.New("id = " + id + " is not exist")

}

func ResetVote() {
	for k, _ := range nodeVoteMap {
		nodeVoteMap[k] = DefaultTicket
	}
}

// // Vote for node. The votes of node is origin vote plus new vote.
// // votes = originVote + vote
// func Vote(context *cli.Context) error {
// 	name := context.String("name")
// 	vote := context.Int("v")

// 	if name == "" {
// 		log.Errorf("节点名称不能为空")
// 	}

// 	if vote < 1 {
// 		log.Errorf("最小投票数目为1")
// 	}

// 	f, err := ioutil.ReadFile(FileName)
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		return err
// 	}
// 	res := strings.Split(string(f), "\n")

// 	voteMap := make(map[string]string)
// 	for _, node := range res {
// 		nodeSplit := strings.Split(node, ":")
// 		if len(nodeSplit) > 1 {
// 			voteMap[nodeSplit[0]] = fmt.Sprintf("%s", nodeSplit[1])
// 		}
// 	}

// 	originVote, err := strconv.Atoi(voteMap[name])
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		return err
// 	}
// 	votes := originVote + vote
// 	voteMap[name] = fmt.Sprintf("%d", votes)

// 	log.Infof("节点%s新增票数%d", name, vote)
// 	str := ""
// 	for k, v := range voteMap {
// 		str += k + ":" + v + "\n"
// 	}

// 	file, err := os.OpenFile(FileName, os.O_RDWR, 0666)
// 	if err != nil {
// 		return err
// 	}

// 	file.WriteString(str)
// 	defer file.Close()

// 	return nil
// }
