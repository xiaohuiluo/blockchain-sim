package cmd

import (
	"math"
	"math/rand"
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
	},
})

func init() {
	simCmd := &grumble.Command{
		Name:    "sim",
		Help:    "simulate blockchain consensus algorithm with p2p",
		Aliases: []string{"run"},
		Args: func(a *grumble.Args) {
			a.String("consensus", "blockchain consensus algorithm")
			a.Int("nodes", "nodes num of p2p")
			a.Int("rounds", "rounds num of simulate")
		},
		Run: func(c *grumble.Context) error {

			// set log level
			logLevel := c.Flags.String("log_level")
			setLogLevel(logLevel)

			// simulate
			consensus := c.Args.String("consensus")
			nodes := c.Args.Int("nodes")
			rounds := c.Args.Int("rounds")

			c.App.Println("simulate with consensus: ", consensus, " nodes: ", nodes, " rounds: ", rounds)
			simulate(c, consensus, nodes, rounds)

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

	// readCmd := &grumble.Command{
	// 	Name:    "read",
	// 	Help:    "read blockchain from a node",
	// 	Aliases: []string{"run"},
	// 	Flags: func(f *grumble.Flags) {
	// 		f.Duration("t", "timeout", time.Second, "timeout duration")
	// 	},
	// 	Args: func(a *grumble.Args) {
	// 		a.String("id", "id of blockchain node")
	// 	},
	// 	Run: func(c *grumble.Context) error {
	// 		logLevel := c.Flags.String("log_level")
	// 		setLogLevel(logLevel)

	// 		// read blockchain from a node
	// 		id := c.Args.String("id")

	// 		data := blockchain.ReadData(id)
	// 		c.App.Println("node id:", id)
	// 		c.App.Println("data:", data)
	// 		return nil
	// 	},
	// }

	// writeCmd := &grumble.Command{
	// 	Name:    "write",
	// 	Help:    "write data to blockchain by a node",
	// 	Aliases: []string{"run"},
	// 	Flags: func(f *grumble.Flags) {
	// 		f.Duration("t", "timeout", time.Second, "timeout duration")
	// 	},
	// 	Args: func(a *grumble.Args) {
	// 		a.Int("data", "data you want to write to blockchain")
	// 	},
	// 	Run: func(c *grumble.Context) error {
	// 		logLevel := c.Flags.String("log_level")
	// 		setLogLevel(logLevel)

	// 		// write blockchain from a node
	// 		data := c.Args.Int("data")
	// 		// pick win write new block node
	// 		winNodeId := blockchain.PickWinnerWithDpos()
	// 		c.App.Println("*****write new block win node=", winNodeId)
	// 		log.Infof("******node=%s win and create new block******", winNodeId)
	// 		rtl, err := blockchain.WriteData(winNodeId, data)
	// 		if rtl {
	// 			c.App.Println("success write data = ", data, " to blockchain")
	// 		} else {
	// 			c.App.Println("failed write data = ", data, " to blockchain, error is", err)
	// 		}

	// 		return nil
	// 	},
	// }

	Cli.AddCommand(simCmd)
	Cli.AddCommand(showCmd)
	// Cli.AddCommand(readCmd)
	// Cli.AddCommand(writeCmd)
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

func simulate(c *grumble.Context, consensus string, nodes int, rounds int) {

	if nodes < 1 {
		log.Errorf("the min of nodes is 1")
	}

	if rounds < 1 {
		log.Errorf("the min of rounds is 1")
	}

	// create nodes
	createBlockChainNodes(c, nodes, rounds)

	// run simulate
	for round := 1; round <= rounds; round++ {
		log.Infof("round: %d", round)
		switch consensus {
		case "dpos":
			runWithDpos(c, nodes, round)
		case "pos":
			runWithPos(c, nodes, round)
		default:
			// default use dpos consensus algorithm
			runWithDpos(c, nodes, round)
		}

		time.Sleep(time.Duration(2) * time.Second)

	}

}

func createBlockChainNodes(c *grumble.Context, nodes int, round int) {
	log.Info("clean previous simulate and init current simulate resource")
	blockchain.InitNodeResource()
	log.Infof("create %d nodes", nodes)
	port := 3000
	addr := ""
	fullAddr := ""
	var err error
	seed := int64(0)

	for node := 1; node <= nodes; node++ {
		log.Infof("create node: %d", node)

		if node == 1 {

			num := 0

			for num < 3 {
				addr, fullAddr, err = blockchain.CreateNode(port, "", seed)

				if err != nil {
					log.Errorf("failed to create node: %d, error is %s", node, err)
					time.Sleep(time.Duration(2) * time.Second)
					num++
					port++
				} else {
					break
				}
			}

			if num == 3 {
				log.Errorf("failed to create first block chain node: %d", node)
				c.App.Println("failed to create node: ", addr)
				return
			}

		} else {

			num := 0

			for num < 3 {
				addr, fullAddr, err = blockchain.CreateNode(port, fullAddr, seed)
				if err != nil {
					log.Errorf("failed to create node: %d, error is %s", node, err)
					time.Sleep(time.Duration(2) * time.Second)
					num++
					port++
				} else {
					break
				}
			}
			if num == 3 {
				log.Errorf("failed to create node: %d", node)
				c.App.Println("failed to create node: ", addr)
				return
			}
		}

		c.App.Println("success create node id = ", addr)

		port++
	}
}

// run with dpos consensus algorithm
func runWithDpos(c *grumble.Context, nodes int, round int) {
	log.Info("---------------------------------------------------------")
	log.Info("Start block chain simulate with dpos consensus algorithm")
	log.Info("init tocken normal distribution")

	c.App.Println("***** start round=", round, " random read and write to block chain *****")
	nodeIds := blockchain.GetNodes()
	tokenData := generateNormalDistribution(100, 20, nodes)
	var sum int64
	for _, value := range tokenData {
		sum += value
	}

	log.Info("compute vote weight")
	weight := make([]float64, nodes)
	for i, value := range tokenData {
		weight[i] = float64(value) / float64(sum)
		log.Infof("node id=%s, token=%d, vote weight=%f", nodeIds[i], value, weight[i])
	}

	log.Info("init vote normal distribution and compute read vote data")
	voteNormalData := generateNormalDistribution(1000, 300, nodes)
	voteData := make([]int, nodes)
	for i, tickets := range voteNormalData {
		voteData[i] = int(math.Floor(float64(tickets) * weight[i]))
	}

	// init real node vote data
	// simple to only vote to themself
	blockchain.InitVoteMap()
	for i, nodeId := range nodeIds {
		blockchain.Vote(nodeId, voteData[i])
		log.Infof("node id=%s, vote=%d", nodeId, voteData[i])
		c.App.Println("node id=", nodeId, " vote tickets=", voteData[i])
	}

	c.App.Println("start random read and write to block chain")
	it := 0
	for it < 10 {
		// read block chain
		log.Info("*****************************************************")
		randNodeId := nodeIds[randInt(0, len(nodeIds)-1)]
		data := blockchain.ReadData(randNodeId)

		c.App.Println("read success from node=", randNodeId)
		log.Infof("read node=%s, data=%s", randNodeId, data)

		time.Sleep(time.Duration(1) * time.Second)
		// new block chain
		randData := randInt(10, 100)
		// pick win write new block node
		winNodeId, err := blockchain.PickWinnerWithDpos()
		if err != nil {
			c.App.Println("failed to pick winner node")
			break
		}

		c.App.Println("*****write new block win node=", winNodeId)
		log.Infof("******node=%s win and create new block******", winNodeId)
		rlt, err := blockchain.WriteData(winNodeId, randData)
		if rlt {
			log.Infof("write new block by node=%s, data=%d", winNodeId, randData)
			c.App.Println("write success from node=", winNodeId, ", data=", randData)

			// revote
			if it%3 == 0 {
				c.App.Println("...revote start...")
				log.Infof("...revote start...")
				c.App.Println("write success from node=", winNodeId, ", data=", randData)

				voteNormalData := generateNormalDistribution(1000, 300, nodes)
				voteData := make([]int, nodes)
				for i, tickets := range voteNormalData {
					voteData[i] = int(math.Floor(float64(tickets) * weight[i]))
				}

				blockchain.InitVoteMap()
				for i, nodeId := range nodeIds {
					blockchain.Vote(nodeId, voteData[i])
					log.Infof("node id=%s, vote=%d", nodeId, voteData[i])
					c.App.Println("node id=", nodeId, " vote tickets=", voteData[i])
				}
				c.App.Println("...revote end...")
				log.Infof("...revote end...")
			}

		} else {
			log.Infof("write new block error by node=%s, data=%s, error=%s", winNodeId, randData, err.Error())
			c.App.Println("write failed from node=", winNodeId, ", data=", randData, ", error=", err.Error())
		}

		it++
		log.Info("*****************************************************")
	}

	c.App.Println("end random read and write to block chain")
	c.App.Println("***** end round=", round, " random read and write to block chain *****")
	log.Info("End block chain simulate with dpos consensus algorithm")
	log.Info("---------------------------------------------------------")
}

func runWithPos(c *grumble.Context, nodes int, round int) {
	log.Info("---------------------------------------------------------")
	log.Info("Start block chain simulate with pos consensus algorithm")
	log.Info("init tocken normal distribution")

	c.App.Println("***** start round=", round, " random read and write to block chain *****")
	nodeIds := blockchain.GetNodes()
	tokenData := generateNormalDistribution(1000, 300, nodes)
	var sum int64
	for _, value := range tokenData {
		sum += value
	}

	log.Info("init age normal distribution and compute tokenAge data")
	ageData := generateNormalDistribution(500, 100, nodes)
	tokenAgeData := make([]int64, nodes)
	for i, age := range ageData {
		tokenAgeData[i] = tokenData[i] * age
	}

	// init real node tokenAge data
	blockchain.InitTokenAgeMap()
	for i, nodeId := range nodeIds {
		blockchain.SetTokenAge(nodeId, tokenAgeData[i])
		log.Infof("node id=%s, tokenAge=%d", nodeId, tokenAgeData[i])
		c.App.Println("node id=", nodeId, " tokenAge=", tokenAgeData[i])
	}

	c.App.Println("start random read and write to block chain")
	it := 0
	for it < 10 {
		// read block chain
		log.Info("*****************************************************")
		randNodeId := nodeIds[randInt(0, len(nodeIds)-1)]
		data := blockchain.ReadData(randNodeId)

		c.App.Println("read success from node=", randNodeId)
		log.Infof("read node=%s, data=%s", randNodeId, data)

		time.Sleep(time.Duration(1) * time.Second)
		// new block chain
		randData := randInt(10, 100)
		// pick win write new block node
		winNodeId, err := blockchain.PickWinnerWithPos()
		if err != nil {
			c.App.Println("failed to pick winner node")
			break
		}
		c.App.Println("*****write new block win node=", winNodeId)
		log.Infof("******node=%s win and create new block******", winNodeId)
		rlt, err := blockchain.WriteData(winNodeId, randData)
		if rlt {
			log.Infof("write new block by node=%s, data=%d", winNodeId, randData)
			c.App.Println("write success from node=", winNodeId, ", data=", randData)

			// set winner node tokenAge to 0, and other node tokenAge +2000
			c.App.Println("...start recompute tokenAge...")
			log.Info("...start recompute tokenAge...")
			blockchain.SetTokenAge(winNodeId, 0)
			for _, nodeId := range nodeIds {
				if nodeId != winNodeId {
					oldTokenAge := blockchain.GetTokenAge(nodeId)
					blockchain.SetTokenAge(nodeId, oldTokenAge+2000)
				}

				log.Infof("node id=%s, tokenAge=%d", nodeId, blockchain.GetTokenAge(nodeId))
				c.App.Println("node id=", nodeId, " tokenAge=", blockchain.GetTokenAge(nodeId))
			}
			c.App.Println("...end recompute tokenAge...")
			log.Info("...end recompute tokenAge...")
		} else {
			log.Infof("write new block error by node=%s, data=%s, error=%s", winNodeId, randData, err.Error())
			c.App.Println("write failed from node=", winNodeId, ", data=", randData, ", error=", err.Error())
		}

		it++
		log.Info("*****************************************************")
	}

	c.App.Println("end random read and write to block chain")
	c.App.Println("***** end round=", round, " random read and write to block chain *****")
	log.Info("End block chain simulate with dpos consensus algorithm")
	log.Info("---------------------------------------------------------")

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
			c.App.Println("------------------------------------------------------------------")
		}
	}
}

// generate normal distribution data
func generateNormalDistribution(mu float64, sigma float64, size int) []int64 {
	dist := distuv.Normal{
		Mu:    mu,    // Mean of the normal distribution
		Sigma: sigma, // Standard deviation of the normal distribution
	}

	data := make([]int64, size)
	for i := 0; i < size; i++ {

		data[i] = int64(dist.Rand())

	}

	return data
}

func randInt(min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min+1) + min
}
