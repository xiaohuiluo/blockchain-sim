// This is the p2p network, handler the conn and communicate with nodes each other.

package blockchain

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	log "github.com/go-fastlog/fastlog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	host "github.com/libp2p/go-libp2p-host"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

// const (
// 	// DefaultVote 节点默认的票数
// 	DefaultVote = 10
// 	// FileName 节点信息保存配置文件
// 	FileName = "config.ini"
// )

var mutex = &sync.Mutex{}

//Validator 定义节点信息
// type Validator struct {
// 	name string
// 	vote int
// }

var consensus string

// node rw map
var nodeRWMap = make(map[string]*bufio.ReadWriter)

// MakeBasicHost 构建P2P网络
func MakeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, string, string, error) {
	var r io.Reader

	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	currentAddr := ""
	currentFullAddr := ""

	// 生产一对公私钥
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, currentAddr, currentFullAddr, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	if !secio {
		opts = append(opts, libp2p.NoSecurity)
	}
	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, currentAddr, currentFullAddr, err
	}

	currentAddr = basicHost.ID().Pretty()
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", currentAddr))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses;
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	currentFullAddr = fullAddr.String()

	log.Infof("create node: %s\n", currentAddr)

	return basicHost, currentAddr, currentFullAddr, nil
}

// HandleStream  handler stream info
func HandleStream(s network.Stream) {
	localPeer := s.Conn().LocalPeer().Pretty()
	remotePeer := s.Conn().RemotePeer().Pretty()

	log.Infof("%s 得到一个新的连接: %s", localPeer, remotePeer)
	// 将连接加入到
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	if nodeRWMap[remotePeer] == nil {
		nodeRWMap[remotePeer] = rw
	}

	// 启动同步数据协程
	go syncData(rw)
}

func SetConsensus(cons string) {
	consensus = cons
}

func GetNodes() []string {
	nodes := make([]string, 0, len(nodeRWMap))
	for k := range nodeRWMap {
		nodes = append(nodes, k)
	}

	return nodes
}

func ReadData(id string) string {
	rw := nodeRWMap[id]
	if rw == nil {
		log.Errorf("node id = %s is not exist", id)
		return "Error: node id = " + id + " is not exist"
	}

	var bt bytes.Buffer

	str, err := rw.ReadString('\n')
	if err != nil {
		log.Errorf(err.Error())
	}

	if str != "\n" {
		chain := make([]Block, 0)

		if err := json.Unmarshal([]byte(str), &chain); err != nil {
			log.Errorf(err.Error())
		}

		mutex.Lock()

		log.Debugf("chain len = %d, BlockChain len = %d", len(chain), len(BlockChain))
		if len(chain) > len(BlockChain) {
			BlockChain = chain
		}

		bytes, err := json.MarshalIndent(BlockChain, "", " ")
		if err != nil {
			log.Errorf(err.Error())
		}

		bt.Write(bytes)

		mutex.Unlock()
	}

	return bt.String()

}

func WriteData(id string, data int) (bool, error) {
	rw := nodeRWMap[id]
	if rw == nil {
		log.Errorf("node id = %s is not exist", id)
		return false, errors.New("Error: node id = " + id + " is not exist")
	}

	// pick选择block生产者
	address := PickWinner()
	log.Infof("******节点 %s 获得了记账权利******", address)
	lastBlock := BlockChain[len(BlockChain)-1]
	newBlock, err := GenerateBlock(lastBlock, data, address)
	if err != nil {
		log.Errorf(err.Error())
	}

	if IsBlockValid(newBlock, lastBlock) {
		mutex.Lock()
		BlockChain = append(BlockChain, newBlock)
		mutex.Unlock()
	}

	spew.Dump(BlockChain)

	bytes, err := json.Marshal(BlockChain)
	if err != nil {
		log.Errorf(err.Error())
	}
	mutex.Lock()
	rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	rw.Flush()
	mutex.Unlock()

	return true, nil
}

func syncData(rw *bufio.ReadWriter) {
	for {
		time.Sleep(2 * time.Second)
		mutex.Lock()
		bytes, err := json.Marshal(BlockChain)
		if err != nil {
			log.Errorf(err.Error())
		}
		mutex.Unlock()

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}
}

func CreateNode(port int, target string, seed int64) (string, string, error) {
	log.Infof("Create Node: port=%d, target=%s, seed=%d", port, target, seed)

	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, CaculateBlockHash(genesisBlock), "", ""}
	BlockChain = append(BlockChain, genesisBlock)

	currentAddr := ""
	currentFullAddr := ""

	if port == 0 {
		log.Fatal("请提供一个端口号")
	}
	// 构造一个host 监听地址
	ha, currentAddr, currentFullAddr, err := MakeBasicHost(port, false, seed)

	if err != nil {
		return currentAddr, currentFullAddr, err
	}

	if target == "" {
		log.Debug("等待节点连接...")
		ha.SetStreamHandler("/p2p/1.0.0", HandleStream)
		return currentAddr, currentFullAddr, err
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", HandleStream)
		ipfsaddr, err := ma.NewMultiaddr(target)
		if err != nil {
			return currentAddr, currentFullAddr, err
		}
		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			return currentAddr, currentFullAddr, err
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			return currentAddr, currentFullAddr, err
		}

		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// 现在我们有一个peerID和一个targetaddr，所以我们添加它到peerstore中。 让libP2P知道如何连接到它。
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		// 构建一个新的stream从hostB到hostA
		// 使用了相同的/p2p/1.0.0 协议
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			return currentAddr, currentFullAddr, err
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		nodeRWMap[s.Conn().RemotePeer().Pretty()] = rw

		go syncData(rw)

		return currentAddr, currentFullAddr, err
	}
}
