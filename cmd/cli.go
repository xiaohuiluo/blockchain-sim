package cmd

import (
	"time"

	"github.com/desertbit/grumble"
	log "github.com/go-fastlog/fastlog"
	"github.com/xiaohuiluo/blockchain-sim/blockchain"
	"gonum.org/v1/gonum/stat/distuv"
)

var Cli = grumble.New(&grumble.Config{
	Name:        "simulate",
	Description: "simulate app for blockchain",
	Flags: func(f *grumble.Flags) {
		f.String("l", "log_level", "info", "log level of app")
		f.String("c", "consensus", "dpos", "blockchain consensus algorithm")
	},
})

func init() {
	simCmd := &grumble.Command{
		Name:    "sim",
		Help:    "simulate blockchain consensus algorithm with p2p",
		Aliases: []string{"run"},
		Args: func(a *grumble.Args) {
			a.Int("nodes", "nodes num of p2p")
			a.Int("rounds", "rounds num of simulate")
		},
		Run: func(c *grumble.Context) error {

			// set log level
			logLevel := c.Flags.String("log_level")
			setLogLevel(logLevel)

			// set consensus algorithm
			consensus := c.Flags.String("consensus")
			blockchain.SetConsensus(consensus)

			// simulate
			nodes := c.Args.Int("nodes")
			rounds := c.Args.Int("rounds")

			c.App.Println("simulate with consensus: ", consensus, " nodes: ", nodes, " rounds: ", rounds)
			sim(c, nodes, rounds)

			return nil
		},
	}

	showCmd := &grumble.Command{
		Name:    "show",
		Help:    "show blockchain node info",
		Aliases: []string{"run"},
		Flags: func(f *grumble.Flags) {
			f.Bool("d", "detail", false, "if show detail")
		},
		Args: func(a *grumble.Args) {
		},
		Run: func(c *grumble.Context) error {

			detail := c.Flags.Bool("detail")
			show(c, detail)

			return nil
		},
	}

	readCmd := &grumble.Command{
		Name:    "read",
		Help:    "read blockchain from a node",
		Aliases: []string{"run"},
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Args: func(a *grumble.Args) {
			a.String("id", "id of blockchain node")
		},
		Run: func(c *grumble.Context) error {
			logLevel := c.Flags.String("log_level")
			setLogLevel(logLevel)

			// read blockchain from a node
			id := c.Args.String("id")

			data := blockchain.ReadData(id)
			c.App.Println("node id:", id)
			c.App.Println("data:", data)
			return nil
		},
	}

	writeCmd := &grumble.Command{
		Name:    "write",
		Help:    "write data to blockchain by a node",
		Aliases: []string{"run"},
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Args: func(a *grumble.Args) {
			a.String("id", "id of blockchain node")
			a.Int("data", "data you want to write to blockchain")
		},
		Run: func(c *grumble.Context) error {
			logLevel := c.Flags.String("log_level")
			setLogLevel(logLevel)

			// write blockchain from a node
			id := c.Args.String("id")
			data := c.Args.Int("data")

			rtl, err := blockchain.WriteData(id, data)
			if rtl {
				c.App.Println("success write data = ", data, " to node id:", id)
			} else {
				c.App.Println("failed write data = ", data, " to node id:", id, " error is", err)
			}

			return nil
		},
	}

	voteCmd := &grumble.Command{
		Name:    "vote",
		Help:    "vote cmd for dpos consensus algorithrm",
		Aliases: []string{"run"},
		Flags: func(f *grumble.Flags) {
		},
		Args: func(a *grumble.Args) {
			a.String("id", "id of blockchain node")
			a.Int("ticket", "ticket value you want to vote")
		},
		Run: func(c *grumble.Context) error {
			consensus := c.Flags.String("consensus")
			if consensus != "dpos" {
				c.App.Println("vote cmd only support simulate with consensus is ", consensus)
				return nil
			}

			return nil
		},
	}

	Cli.AddCommand(simCmd)
	Cli.AddCommand(showCmd)
	Cli.AddCommand(readCmd)
	Cli.AddCommand(writeCmd)
	Cli.AddCommand(voteCmd)
}

func setLogLevel(logLevel string) {
	if logLevel == "error" {
		log.SetFlags(log.Flags() | log.Lerror)
	} else if logLevel == "warn" {
		log.SetFlags(log.Flags() | log.Lwarn)
	} else if logLevel == "info" {
		log.SetFlags(log.Flags() | log.Linfo)
	} else if logLevel == "debug" {
		log.SetFlags(log.Flags() | log.Ldebug)
	} else {
		log.SetFlags(log.LstdFlags)
	}
}

func sim(c *grumble.Context, nodes int, rounds int) {

	if nodes < 1 {
		log.Errorf("最小节点数为1")
	}

	if rounds < 1 {
		log.Errorf("最小模拟周期数为1")
	}

	log.Info("初始化高斯分布")
	dist := distuv.Normal{
		Mu:    100, // Mean of the normal distribution
		Sigma: 20,  // Standard deviation of the normal distribution
	}

	var sum_c int64 = 0
	data := make([][]int64, nodes)
	for i := 0; i < nodes; i++ {
		data[i] = make([]int64, rounds)
		log.Debugf("初始化 node = %d 高斯分布", i+1)

		sum_c = 0
		for j := 0; j < rounds; j++ {
			// Generate some random numbers from standard normal distribution
			data[i][j] = int64(dist.Rand())
			sum_c += data[i][j]
			log.Debugf("元素值：%d", data[i][j])
		}

		log.Debugf("node= %d, all rounds 平均值%d", i+1, sum_c/int64(len(data[i])))
	}

	sum_round := make([]int64, rounds)

	for i := 0; i < rounds; i++ {

		sum_round[i] = 0
		for j := 0; j < nodes; j++ {
			sum_round[i] += data[j][i]
		}

		log.Debugf("round=%d, sum=%d", i+1, sum_round[i])
	}

	weight := make([][]float64, nodes)

	for i := 0; i < nodes; i++ {
		for j := 0; j < rounds; j++ {
			weight[i] = make([]float64, rounds)
			weight[i][j] = float64(data[i][j]) / float64(sum_round[j])
			log.Debugf("node=%d, round=%d, weight=%f", i+1, j+1, weight[i][j])
		}
	}

	log.Infof("create %d nodes", nodes)
	port := 3000
	addr := ""
	fullAddr := ""
	seed := int64(0)

	for node := 1; node <= nodes; node++ {
		log.Infof("create node: %d", node)

		if node == 1 {
			addr, fullAddr, _ = blockchain.CreateNode(port, "", seed)
		} else {
			addr, fullAddr, _ = blockchain.CreateNode(port, fullAddr, seed)
		}

		c.App.Println("success create node id = ", addr)

		port++
	}

	// for round := 1; round <= rounds; round++ {
	// 	log.Infof("round: %d", round)

	// 	time.Sleep(time.Duration(2) * time.Second)

	// }

}

func show(c *grumble.Context, detail bool) {
	nodes := blockchain.GetNodes()
	if len(nodes) < 1 {
		c.App.Println("Empty node")
		return
	}

	c.App.Println("----------------------------------------------")
	for _, node := range nodes {
		c.App.Println(node)
	}
	c.App.Println("----------------------------------------------")

	if detail {
		for _, node := range nodes {
			c.App.Println("node id: ", node)
			c.App.Println("data: ", blockchain.ReadData(node))
		}
	}
}
