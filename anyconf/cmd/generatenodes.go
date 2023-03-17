package cmd

import (
	"fmt"
	"github.com/anytypeio/any-sync-tools/gen"
	"github.com/anytypeio/any-sync/nodeconf"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	addressesFlag = "addresses"
	nodesYaml     = "nodes.yml"
)

var generateNodes = &cobra.Command{
	Use:   "generate-nodes",
	Short: "Generate nodes",
	Args:  cobra.RangeArgs(0, 10),
	Run: func(cmd *cobra.Command, args []string) {
		addresses, err := cmd.Flags().GetStringArray(addressesFlag)
		types, err := cmd.Flags().GetStringArray(typesFlag)

		var nodeTypes []nodeconf.NodeType
		for _, nodeType := range types {
			nodeType := nodeconf.NodeType(nodeType)

			if !slices.Contains(validOptions, nodeType) {
				fmt.Println(nodeType)
				panic("Wrong node 'type' parameter")
			}

			nodeTypes = append(nodeTypes, nodeType)
		}

		nodesList, accountsList, err := gen.GenerateNodesConfigs(nodeTypes, addresses)
		nodes := Nodes{nodesList}

		nodesBytes, err := yaml.Marshal(nodes)
		if err != nil {
			panic(fmt.Sprintf("could Marshal nodes: %v", err))
		}

		err = os.WriteFile(nodesYaml, nodesBytes, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("could not write nodes.yml to file: %v", err))
		}

		for index, account := range accountsList {
			pc := PrivateConf{Account: account}

			accountBytes, err := yaml.Marshal(pc)

			accountFilePath := fmt.Sprintf("account%d.yml", index)

			err = os.WriteFile(accountFilePath, accountBytes, os.ModePerm)
			if err != nil {
				panic(fmt.Sprintf("could not write accountBytes to file: %v", err))
			}
		}
	},
}

func init() {
	generateNodes.Flags().StringArray(typesFlag, []string{}, "fill this flag with one of three options [tree, file, coordinator]")
	generateNodes.MarkFlagRequired(typesFlag)

	generateNodes.Flags().StringArray(addressesFlag, []string{}, "")
}
