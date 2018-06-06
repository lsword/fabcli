package main

import (
	//"bytes"
	//"encoding/json"
	"flag"
	"fmt"
	//"log"
	//"io/ioutil"
	"os"
	//"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	// cliconfig "liujian/fabcli/config"
	"liujian/fabcli/common/log"
	//"liujian/fabcli/fabclient"
	httpserver "liujian/fabcli/httpserver"
	//"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/modlog"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	_ "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	_ "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/msp"
	_ "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/orderer"
	_ "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	//putils "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"
	//"github.com/op/go-logging"
	//"github.com/golang/protobuf/proto"
)

func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fabcli",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//cmd.HelpFunc()(cmd, args)
		},
	}

	//flags := cmd.PersistentFlags()
	//cliconfig.InitConfigFile(flags)

	return cmd
}

var signalHandlerMap map[os.Signal]func()

func SigHupHandler() {
	glog.Flush()
	fmt.Println("sig hup")
	os.Exit(0)
}

func SigIntHandler() {
	glog.Flush()
	fmt.Println("sig int")
	os.Exit(0)
}

func SetSignalHandler(handler func(), sig os.Signal) {
	if signalHandlerMap == nil {
		signalHandlerMap = make(map[os.Signal]func())
	}
	signalHandlerMap[sig] = handler
}

func dealSignal() {
	SetSignalHandler(SigHupHandler, syscall.SIGHUP)
	SetSignalHandler(SigIntHandler, syscall.SIGINT)

	signalChan := make(chan os.Signal, 1)
	var signals []os.Signal
	for k := range signalHandlerMap {
		signals = append(signals, k)
	}
	signal.Notify(signalChan, signals...)
	go func() {
		for sig := range signalChan {
			if signalHandlerMap[sig] != nil {
				signalHandlerMap[sig]()
			}
		}
	}()
}

var fabric_bin_path = "/Users/liujian/fabric/bin"
var fabric_sysconfig_block_path = "/Users/liujian/code/gopath/src/liujian/fabcli/sys_config.block"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go dealSignal()

	flag.Parse()
	flag.Set("log_dir", "/Users/liujian/code/gopath/src/liujian/fabcli/log")
	defer glog.Flush()

	modlog.InitLogger(&log.CliLoggingProvider{})
	modlog.SetLevel("Main", api.DEBUG)
	logger := modlog.LoggerProvider().GetLogger("Main")
	logger.Info("fabcli start")

	if cmd().Execute() != nil {
		fmt.Println("cmd.Execute error")
		os.Exit(1)
	}

	httpserver.Run()

	//"channel", "fetch", "config", "/Users/liujian/code/gopath/src/liujian/fabcli/sys_config_block.pb", "-o", "orderer.example.com:7050", "-c", "testchainid", "--tls", "--cafile", "$ORDERER_CA"

	/*
		cmd := exec.Command(fabric_bin_path+"/peer",
			"channel",
			"fetch",
			"config",
			fabric_sysconfig_block_path,
			"-o",
			"orderer.example.com:7050",
			"-c",
			"testchainid",
			"--tls",
			"--cafile",
			"/Users/liujian/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem")
		cmd.Env = append(os.Environ(),
			"FABRIC_CFG_PATH=/Users/liujian/fabric/config",
			"ORDERER_CA=/Users/liujian/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem",
			"CORE_PEER_LOCALMSPID=OrdererMSP",
			"CORE_PEER_TLS_ROOTCERT_FILE=$ORDERER_CA",
			"CORE_PEER_MSPCONFIGPATH=/Users/liujian/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp",
		)
		var out bytes.Buffer
		cmd.Stderr = &out
		cmd.Stdout = &out
		if err = cmd.Run(); err != nil {
			fmt.Printf("exec peer error:%v\n", err)
			os.Exit(1)
		}

		cmd = exec.Command(fabric_bin_path+"/configtxlator",
			"proto_decode",
			"--type=common.Block",
			"--input="+fabric_sysconfig_block_path,
			"--output="+fabric_sysconfig_block_path+".json")
		if err = cmd.Run(); err != nil {
			fmt.Printf("exec peer error:%v\n", err)
			os.Exit(1)
		}

		blockJson, err := ioutil.ReadFile(fabric_sysconfig_block_path + ".json")
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(blockJson))
	*/

	runtime.Goexit()
}
