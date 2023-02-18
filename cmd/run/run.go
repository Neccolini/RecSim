package run

import (
	"fmt"
	"log"

	"github.com/Neccolini/RecSimu/cmd/instruction"
	"github.com/Neccolini/RecSimu/cmd/node"
)

type SimulationConfig struct {
	nodeNum       int
	totalCycle    int
	adjacencyList map[int][]int
	nodes         map[int]*node.Node
}

func NewSimulationConfig(nodeNum int, cycle int, adjacencyList map[int][]int, nodesType map[int]string) *SimulationConfig {
	config := &SimulationConfig{}
	config.nodeNum = nodeNum
	config.totalCycle = cycle
	config.adjacencyList = adjacencyList
	config.nodes = make(map[int]*node.Node, nodeNum)
	for i := 0; i < nodeNum; i++ {
		nodeI, err := node.NewNode(i, nodesType[i], []instruction.Instruction{})
		if err != nil {
			log.Fatal(err)
		}
		config.nodes[i] = nodeI
	}
	return config
}

func (config *SimulationConfig) Simulate(outputFile string) error {
	// サイクルごとのシミュレートを実行
	for cycle := 1; cycle <= config.totalCycle; cycle++ {
		// todo トポロジーの変更

		// シミュレートを実行
		if err := config.SimulateCycle(cycle); err != nil {
			return err
		}
		// todo 各サイクル後の状態を記録
		fmt.Printf("cycle %d\n", cycle)
		for _, node := range config.nodes {
			fmt.Println(node.String())
		}
	}
	return nil
}

func (config *SimulationConfig) SimulateCycle(cycle int) error {
	// ノードごとに送信
	for _, node := range config.nodes {
		node.CycleSend()
	}
	// 衝突したらどうにかする（ランダムで一つ選んで他は待機＋再送信）
	// メッセージを配信
	for _, node := range config.nodes {
		if !node.SendingMessage.IsEmpty() {
			// メッセージをブロードキャストする
			for _, adjacentNodeId := range config.adjacencyList[node.Id()] {
				config.nodes[adjacentNodeId].Receive(node.SendingMessage)
			}
			config.nodes[node.Id()].SendingMessage.Clear()
		}
	}

	for _, node := range config.nodes {
		node.CycleReceive()
		node.SimulateCycle()
	}

	return nil
}
