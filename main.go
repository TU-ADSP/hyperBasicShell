package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"os/exec"
	"strings"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("connect",
		readline.PcItem("Org1"),
		readline.PcItem("Org2"),
	),
	readline.PcItem("CreateAsset"),
	readline.PcItem("ReadAsset"),
	readline.PcItem("UpdateAsset"),
	readline.PcItem("DeleteAsset"),
	readline.PcItem("AssetExists"),
	readline.PcItem("TransferAsset"),
	readline.PcItem("GetAllAssets"),
	readline.PcItem("exit"),
	readline.PcItem("quit"),
)

var (
	envOrg1 = []string{
		"PATH=/usr/share/blf-hyperledger/fabric-samples/test-network/../bin:$PATH",
		"FABRIC_CFG_PATH=/usr/share/blf-hyperledger/fabric-samples/test-network/../config/",
		"CORE_PEER_TLS_ENABLED=true",
		"CORE_PEER_LOCALMSPID=Org1MSP",
		"CORE_PEER_TLS_ROOTCERT_FILE=/usr/share/blf-hyperledger/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
		"CORE_PEER_MSPCONFIGPATH=/usr/share/blf-hyperledger/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp",
		"CORE_PEER_ADDRESS=localhost:7051",
	}
	envOrg2 = []string{
		"PATH=/usr/share/blf-hyperledger/fabric-samples/test-network/../bin:$PATH",
		"FABRIC_CFG_PATH=/usr/share/blf-hyperledger/fabric-samples/test-network/../config/",
		"CORE_PEER_TLS_ENABLED=true",
		"CORE_PEER_LOCALMSPID=Org2MSP",
		"CORE_PEER_TLS_ROOTCERT_FILE=/usr/share/blf-hyperledger/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt",
		"CORE_PEER_MSPCONFIGPATH=/usr/share/blf-hyperledger/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp",
		"CORE_PEER_ADDRESS=localhost:9051",
	}
	selectedEnv = []string{}
	connected   = false
)

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          " > ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
		//		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		switch {
		case line == "exit" || line == "quit":
			os.Exit(0)
		case line == "connect Org1":
			selectedEnv = envOrg1
			connected = true
		case line == "connect Org2":
			selectedEnv = envOrg2
			connected = true
		case strings.HasPrefix(line, "CreateAsset") || strings.HasPrefix(line, "ReadAsset") || strings.HasPrefix(line, "UpdateAsset") || strings.HasPrefix(line, "DeleteAsset") || strings.HasPrefix(line, "AssetExists") || strings.HasPrefix(line, "TransferAsset") || strings.HasPrefix(line, "GetAllAssets") || strings.HasPrefix(line, "exit") || strings.HasPrefix(line, "quit"):
			if !connected {
				fmt.Println("Please connect to an Org first!")
				continue
			}

			cmdPieces := strings.Split(line, " ")

			if strings.HasPrefix(line, "CreateAsset") || strings.HasPrefix(line, "TransferAsset") {
				if len(cmdPieces) > 1 {
					argsText := ""
					for i := 1; i < len(cmdPieces)-1; i++ {
						argsText += "\"" + cmdPieces[i] + "\","
					}
					argsText += "\"" + cmdPieces[len(cmdPieces)-1] + "\""
					fabricSamplesFolder := "/usr/share/blf-hyperledger/fabric-samples/test-network"
					cmd := exec.Command("/bin/bash", "-c", "peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "+fabricSamplesFolder+"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "+fabricSamplesFolder+"/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles "+fabricSamplesFolder+"/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{\"function\":\""+cmdPieces[0]+"\",\"Args\":["+argsText+"]}' | /usr/bin/jq .")
					cmd.Env = selectedEnv
					cmd.Wait()
					stdout, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("%s\n", stdout)
				} else {
					fmt.Println("You need to have additional arguments for these commands!")
				}
			} else {
				argsText := ""
				for i := 0; i < len(cmdPieces)-1; i++ {
					argsText += "\"" + cmdPieces[i] + "\","
				}
				argsText += "\"" + cmdPieces[len(cmdPieces)-1] + "\""
				cmd := exec.Command("/bin/bash", "-c", "peer chaincode query -C mychannel -n basic -c '{\"Args\":["+argsText+"]}' | /usr/bin/jq .")
				cmd.Env = selectedEnv
				cmd.Wait()
				stdout, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%s\n", stdout)
			}
		default:
			fmt.Println("Possible commands: 'connect Org1', 'connect Org2', 'CreateAsset', 'ReadAsset', 'UpdateAsset', 'DeleteAsset', 'AssetExists', 'TransferAssets', 'GetAllAssets', 'exit' or 'quit'")
		}
	}
}
