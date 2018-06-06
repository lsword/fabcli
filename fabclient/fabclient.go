package fabclient

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	"net/http"
	"bytes"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang/protobuf/proto"
	//pkgfab "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	//"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	//"github.com/hyperledger/fabric-sdk-go/pkg/msp"
	rwsetutil "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"

	fpc "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	//"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	//fpp "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	fpu "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"
	//selection "github.com/hyperledger/fabric-sdk-go/pkg/client/common/selection/dynamicselection"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/discovery/dynamicdiscovery"
  "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defsvc"

	"liujian/fabcli/common/log"
)

type FabricClient struct {
	sdk *fabsdk.FabricSDK
	configBackend     core.ConfigBackend
	configCryptoSuite core.CryptoSuiteConfig
	configEndpoint    fab.EndpointConfig
	configIdentity    msp.IdentityConfig
	channelClient     *channel.Client
	ledgerClient      *ledger.Client
	resmgmtClient     *resmgmt.Client
	eventClient       *event.Client
	blockDB           *sql.DB

	sysResmgmtClient *resmgmt.Client
	sysLedgerClient  *ledger.Client
}

var fclient FabricClient

type ChainCodeResponse struct {
		Status int32 `json:"status"`
		Payload string `json:"payload,omitempty"`
}

type OrgAnchorPeer struct {
	Org  string `json:"org"`
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type ChannelConfig struct {
	Id string `json:"id"`
	BlockNumber uint64 `json:"blocknumber"`
	//MSPs []*mspCfg.MSPConfig
	AnchorPeers []*OrgAnchorPeer `json:"anchorpeers"`
	Orderers []string `json:"orderers"`
	//Versions() *Versions
}

type Peer struct {
	Url string `json:"url"`
	Mspid string `json:"mspid"`
}

// DynamicDiscoveryProviderFactory is configured with dynamic (endorser) selection provider
type DynamicDiscoveryProviderFactory struct {
	defsvc.ProviderFactory
}

// CreateDiscoveryProvider returns a new dynamic discovery provider
func (f *DynamicDiscoveryProviderFactory) CreateDiscoveryProvider(config fab.EndpointConfig) (fab.DiscoveryProvider, error) {
	return dynamicdiscovery.New(config), nil
}

// CreateLocalDiscoveryProvider returns a new local dynamic discovery provider
func (f *DynamicDiscoveryProviderFactory) CreateLocalDiscoveryProvider(config fab.EndpointConfig) (fab.LocalDiscoveryProvider, error) {
	return dynamicdiscovery.New(config), nil
}

func ChannelQueryPeers() string {
	//chProvider := fclient.sdk.ChannelContext("mychannel", fabsdk.WithUser("admin"), fabsdk.WithOrg("org1"))
	chProvider := fclient.sdk.ChannelContext("mychannel")
	chCtx, err := chProvider()
	if err != nil {
		fmt.Printf("chProvider error:", err)
	}

	peers, err := chCtx.DiscoveryService().GetPeers()
	if err != nil {
		fmt.Printf("getpeers error:", err)
	}

	var peerarray []Peer

	for _, p := range peers {
		peerarray = append(peerarray, Peer{p.URL(), p.MSPID()})
	}
	res, _ := json.Marshal(peerarray)
	return string(res)
}

/*
func Query() {
	response, err := fclient.channelClient.Query(channel.Request{ChaincodeID: "mycc", Fcn: "query", Args: [][]byte{[]byte("b")}})
	//response, err := client.Invoke(channel.Request{ChaincodeID: "mycc", Fcn: "invoke", Args: [][]byte{[]byte("query"), []byte("a"), []byte("b")}})
	if err != nil {
		fmt.Sprintf("Failed to query funds after transaction: %s", err)
	}

	fmt.Printf(fmt.Sprintf("payload:%s\n", string(response.Payload)))
}
func Execute() {
	response, err := fclient.channelClient.Execute(channel.Request{ChaincodeID: "mycc", Fcn: "invoke", Args: [][]byte{[]byte("a"), []byte("b"), []byte("1")}})
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to move funds: %s", err))
	}
	fmt.Printf(fmt.Sprintf("payload:%s\n", string(response.Payload)))
}
*/


func ChannelQuery(chaincodeid string, args [][]byte) string {
	response, err := fclient.channelClient.Query(channel.Request{ChaincodeID: chaincodeid, Fcn: "query", Args: args})
	//response, err := client.Invoke(channel.Request{ChaincodeID: "mycc", Fcn: "invoke", Args: [][]byte{[]byte("query"), []byte("a"), []byte("b")}})
	if err != nil {
		fmt.Sprintf("Failed to query funds after transaction: %s", err)
	}

	//return string(response.Payload)
	rsp := &ChainCodeResponse{response.ChaincodeStatus, string(response.Payload)}
	b, _ := json.Marshal(rsp)
	return string(b)
}

func ChannelExecute(chaincodeid string, args [][]byte) string {
	response, err := fclient.channelClient.Execute(channel.Request{ChaincodeID: chaincodeid, Fcn: "invoke", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to move funds: %s", err))
	}

	rsp := &ChainCodeResponse{response.ChaincodeStatus, string(response.Payload)}
	b, _ := json.Marshal(rsp)
	return string(b)
}

func LedgerQueryBlockChain() (*fab.BlockchainInfoResponse, error) {
	info, err := fclient.ledgerClient.QueryInfo()
	if err != nil {
		fmt.Printf("QueryBlockByHash return error: %v", err)
		return nil, err
	}
	fmt.Printf("Endorser:%s, Status:%d\n", info.Endorser, info.Status)
	bci, err := json.Marshal(info.BCI)
	fmt.Println(string(bci))
	return info, nil
}

func LedgerQueryInfo() string {
	info, err := fclient.ledgerClient.QueryInfo()
	if err != nil {
		fmt.Printf("Ledger_QueryInfo return error: %v", err)
		return ""
	}
	fmt.Printf("Endorser:%s, Status:%d\n", info.Endorser, info.Status)
	bci, _ := json.Marshal(info.BCI)
	return string(bci)
}

func LedgerQueryConfig() string {
	cfg, err := fclient.ledgerClient.QueryConfig()
	if err != nil {
		fmt.Printf("Ledger_QueryConfig return error: %v", err)
		return ""
	}
	var channelConfig ChannelConfig
	channelConfig.Id = cfg.ID()
	channelConfig.BlockNumber = cfg.BlockNumber()
	//channelConfig.AnchorPeers
	anchorPeers := cfg.AnchorPeers()
	for _, v := range anchorPeers {
		var item OrgAnchorPeer
		item.Org = v.Org
		item.Host = v.Host
		item.Port = v.Port
		channelConfig.AnchorPeers = append(channelConfig.AnchorPeers, &item)
	}
	res, _ := json.Marshal(channelConfig)
	return string(res)
}

func LedgerQueryBlock(num uint64) string {
	block, err := fclient.ledgerClient.QueryBlock(num)
	if err != nil {
		fmt.Printf("Ledger_QueryBlock return error: %v", err)
		return ""
	}
	//header, _ := json.Marshal(block.Header)

	b, err := proto.Marshal(block)
	if err != nil {
		return ""
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Block", "application/octet-stream", bytes.NewReader(b))
	resBody, err := ioutil.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	if err != nil {
		return ""
	}
	return string(resBody)
}

func LedgerQueryBlockByHash(blockhash []byte) string {
	block, err := fclient.ledgerClient.QueryBlockByHash(blockhash)
	if err != nil {
		fmt.Printf("Ledger_QueryBlockByHash return error: %v", err)
		return ""
	}
	//header, _ := json.Marshal(block.Header)

	b, err := proto.Marshal(block)
	if err != nil {
		return ""
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Block", "application/octet-stream", bytes.NewReader(b))
	resBody, err := ioutil.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	if err != nil {
		return ""
	}
	return string(resBody)
}

func LedgerQueryBlockByTxID(txid string) string {
	block, err := fclient.ledgerClient.QueryBlockByTxID(fab.TransactionID(txid))
	if err != nil {
		fmt.Printf("Ledger_QueryBlockByTxID return error: %v", err)
		return ""
	}
	//header, _ := json.Marshal(block.Header)

	b, err := proto.Marshal(block)
	if err != nil {
		return ""
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/common.Block", "application/octet-stream", bytes.NewReader(b))
	resBody, err := ioutil.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	if err != nil {
		return ""
	}
	return string(resBody)
}

func LedgerQueryTransaction(txid string) string {
	tx, err := fclient.ledgerClient.QueryTransaction(fab.TransactionID(txid))
	if err != nil {
		fmt.Printf("Ledger_QueryTransaction return error: %v", err)
		return ""
	}
	//header, _ := json.Marshal(block.Header)

	b, err := proto.Marshal(tx)
	if err != nil {
		return ""
	}
	httpResp, err := http.Post("http://127.0.0.1:7059/protolator/decode/protos.ProcessedTransaction", "application/octet-stream", bytes.NewReader(b))
	resBody, err := ioutil.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	if err != nil {
		return ""
	}
	return string(resBody)
}

func ResmgmtQueryConfigFromOrderer(channelID string) string {
	//cfg, err := fclient.resmgmtClient.QueryConfigFromOrderer("mychannel")
	cfg, err := fclient.sysResmgmtClient.QueryConfigFromOrderer("testchainid")
	if err != nil {
		fmt.Printf("Ledger_QueryConfig return error: %v", err)
		return ""
	}
	var channelConfig ChannelConfig
	channelConfig.Id = cfg.ID()
	channelConfig.BlockNumber = cfg.BlockNumber()
	//channelConfig.AnchorPeers
	anchorPeers := cfg.AnchorPeers()
	for _, v := range anchorPeers {
		var item OrgAnchorPeer
		item.Org = v.Org
		item.Host = v.Host
		item.Port = v.Port
		channelConfig.AnchorPeers = append(channelConfig.AnchorPeers, &item)
	}
	res, _ := json.Marshal(channelConfig)
	return string(res)
}

func ResmgmtQueryChannels() string {
	channels, _ := fclient.resmgmtClient.QueryChannels(resmgmt.WithTargetEndpoints("peer0.org1.example.com:7051"))
	channelsJson, _ := json.Marshal(channels)
	fmt.Println(string(channelsJson))
	return string(channelsJson)
}

func ResmgmtQueryInstalledChaincodes() string {
	chaincodes, _ := fclient.resmgmtClient.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints("peer0.org1.example.com:7051"))
	chaincodesJson, _ := json.Marshal(chaincodes)
	fmt.Println(string(chaincodesJson))
	return string(chaincodesJson)
	/*
		configBackend, _ := config.FromFile("/Users/liujian/code/gopath/src/liujian/fabcli/config.yaml")()
		_, endpointConfig, _, _ := config.FromBackend(configBackend)()
		peersConfig, _ := endpointConfig.PeersConfig("org1")
		for index := range peersConfig {
			fmt.Printf("peer[%d].URL:%s\n", index, peersConfig[index].URL)
		}
		peerConfig, _ := endpointConfig.PeerConfig("org1", "peer0.org1.example.com")
		fmt.Printf("peer.URL:%s\n", peerConfig.URL)
	*/
}

func ResmgmtQueryInstantiatedChaincodes() string {
	chaincodes, _ := fclient.resmgmtClient.QueryInstantiatedChaincodes("mychannel", resmgmt.WithTargetEndpoints("peer0.org1.example.com:7051"))
	chaincodesJson, _ := json.Marshal(chaincodes)
	fmt.Println(string(chaincodesJson))
	return string(chaincodesJson)
}

func ResmgmtInstallCC(name, version, path string) string {
	ccpackage, err := gopackager.NewCCPackage("github.com/chaincode/chaincode_example02/go", "/Users/liujian/code/gopath")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return "newccpackage err"
	}
	fmt.Println(ccpackage.Type, "\n", len(ccpackage.Code))

	// Chaincode that is not installed already (it will be installed)
	req := resmgmt.InstallCCRequest{Name: name, Version: version, Path: path, Package: ccpackage}
	//responses, err := resmgmtClient.InstallCC(req, WithTargets(&peer2))
	responses, err := fclient.resmgmtClient.InstallCC(req, resmgmt.WithTargetEndpoints("peer0.org1.example.com:7051"))
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return "err"
	}
	if responses == nil || len(responses) != 1 {
		fmt.Printf("Should have one successful response\n")
		return "err"
	}
	return "success"
}

func ResmgmtInstantiateCC(name, version, path, channel string) string {
	// Valid request
	ccPolicy := cauthdsl.SignedByMspMember("Org1MSP")
	req := resmgmt.InstantiateCCRequest{Name: name, Version: version, Path: path, Policy: ccPolicy}

	req.Args = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
	// Test both targets and filter provided (error condition)
	//_, err := rc.InstantiateCC("mychannel", req, WithTargets(peers...), WithTargetFilter(&mspFilter{mspID: "Org1MSP"}))
	_, err := fclient.resmgmtClient.InstantiateCC(channel, req, resmgmt.WithTargetEndpoints("peer0.org1.example.com:7051"))
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return "err"
	}
	return "success"
}

func ResmgmtUpgradeCC(name, version, path, channel string) string {
	// Valid request
	ccPolicy := cauthdsl.SignedByMspMember("Org1MSP")
	req := resmgmt.UpgradeCCRequest{Name: name, Version: version, Path: path, Policy: ccPolicy}

	req.Args = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
	// Test both targets and filter provided (error condition)
	//_, err := rc.InstantiateCC("mychannel", req, WithTargets(peers...), WithTargetFilter(&mspFilter{mspID: "Org1MSP"}))
	_, err := fclient.resmgmtClient.UpgradeCC(channel, req, resmgmt.WithTargetEndpoints("peer0.org1.example.com:7051"))
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return "err"
	}
	return "success"
}

func QueryBlockChainHeight(channelName string) (uint64, error) {
	bir, err := LedgerQueryBlockChain()
	if err != nil {
		return 0, err
	}
	return bir.BCI.Height, nil
}

func SaveBlockToDB(block *fpc.Block) error {
	//blockHeaderJson, _ := json.Marshal(block.Header)
	//fmt.Println("block header===================")
	//fmt.Println(string(blockHeaderJson))
	//fmt.Println("block header===================")


	b, _ := proto.Marshal(block)

	//fmt.Printf("datahash:%s\n", hex.EncodeToString(block.Header.DataHash))
	stmt, _ := fclient.blockDB.Prepare("INSERT INTO blocks(blocknum, datahash, prehash, channelname, txcount, data) values(?,?,?,?,?,?)")
	//fmt.Println("block header===================")
	stmt.Exec(block.Header.Number, hex.EncodeToString(block.Header.DataHash), hex.EncodeToString(block.Header.PreviousHash), "mychannel", len(block.Data.Data), b)

	for i := range block.Data.Data {
		envelope, err := fpu.ExtractEnvelope(block, i)
		if err != nil {
			break
		}
		payload, err := fpu.ExtractPayload(envelope)
		if err != nil {
			break
		}
		channelHeader := &fpc.ChannelHeader{}
		//proto.Unmarshal(payload.Header.ChannelHeader, header)
		//fmt.Printf(payload.String())
		channelHeader, _ = fpu.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if channelHeader.Type != int32(fpc.HeaderType_ENDORSER_TRANSACTION) {
			continue
		}
		dt := time.Unix(channelHeader.Timestamp.GetSeconds(), int64(channelHeader.Timestamp.GetNanos()))
		fmt.Println(dt)

		//channelHeaderJson, _ := json.Marshal(channelHeader)
		//fmt.Println("payload header===================")
		//fmt.Println(string(channelHeaderJson))
		//fmt.Println("payload header===================")

		tx, _ := fpu.GetTransaction(payload.Data)
		//txjson, _ := json.Marshal(tx)
		//fmt.Println("payload data===================")
		//fmt.Println(string(txjson))
		//fmt.Println("payload data===================")

		proposal, _ := fpu.GetChaincodeActionPayload(tx.Actions[0].Payload)
		//proposaljson, _ := json.Marshal(proposal)
		//fmt.Println("actions[0].payload===================")
		//fmt.Println(string(proposaljson))
		//fmt.Println("actions[0].payload===================")

		propRespPayload, err := fpu.GetProposalResponsePayload(proposal.Action.ProposalResponsePayload)
		//proposalrspjson, _ := json.Marshal(propRespPayload)
		//fmt.Println(string(proposalrspjson))

		ccAction, err := fpu.GetChaincodeAction(propRespPayload.Extension)
		//ccActionjson, _ := json.Marshal(ccAction)
		//fmt.Println("ACTION:", string(ccActionjson))

		//ccEvent, err := fpu.GetChaincodeEvents(ccAction.Events)
		//ccEventjson, _ := json.Marshal(ccEvent)
		//fmt.Println(string(ccEventjson))

		txRWSets := &rwsetutil.TxRwSet{}
		txRWSets.FromProtoBytes(ccAction.Results)
		//txRWSetsJson, _ := json.Marshal(txRWSets)
		//fmt.Println("rwsets===================")
		//fmt.Println(string(txRWSetsJson))
		//fmt.Println("rwsets===================")
		//fmt.Println("---->", txRWSets.NsRwSets[1].NameSpace)

		fclient.blockDB.Exec("INSERT INTO transactions(channelname, blockid, txhash, createdate, chaincodename) VALUES (?,?,?,?,?)",
			channelHeader.ChannelId, block.Header.Number, channelHeader.TxId, dt, txRWSets.NsRwSets[1].NameSpace)

		//txReadWriteSet := &rwset.TxReadWriteSet{}
		//proto.Unmarshal(ccAction.Results, txReadWriteSet)
		//txReadWriteSetJson, _ := json.Marshal(txReadWriteSet)
		//fmt.Println(string(txReadWriteSetJson))
	}

	return nil
}

func QueryBlock(blockNum uint64) (*fpc.Block, error) {
	block, err := fclient.ledgerClient.QueryBlock(blockNum)
	if err != nil {
		fmt.Printf("QueryBlockByHash return error: %v", err)
		return nil, err
	}
	if block.Data == nil {
		fmt.Printf("QueryBlockByHash block data is nil")
		return nil, fmt.Errorf("Cannot find block by blocknum(%s)", blockNum)
	}

	return block, nil
}

func SyncBlocks() {
	t := time.NewTimer(time.Second * 5)
	for {
    select {
    case <-t.C:

			var curBlockNum uint64
			curBlockNum = 0
			maxBlockNum, err := QueryBlockChainHeight("")
			if err != nil {
				return
			}
			rows, err := fclient.blockDB.Query("SELECT max(blocknum) FROM blocks WHERE channelname='mychannel'")
			if err != nil {
				fmt.Printf("error:%v", err)
				return
			}
			defer rows.Close()
			for rows.Next() {
				rows.Scan(&curBlockNum)
			}
			for i := curBlockNum; i < maxBlockNum-1; i++ {
				block, err := QueryBlock(uint64(i))
				if err != nil {
					return
				}
				err = SaveBlockToDB(block)
				if err != nil {
					return
				}
			}

      t.Reset(time.Second * 5)
    }
  }
}

func GetBlockData(channelid string, displayStart int, displayLength int) string {
		type BlockData struct {
			BlockNum int `json:"blocknum"`
			TxCount int `json:"txcount"`
			DataHash string `json:"datahash"`
			PreHash string `json:"prehash"`
		}
		type BlockDataResult struct {
			TotalRecords int `json:"iTotalRecords"`
			TotalDisplayRecords int `json:"iTotalDisplayRecords"`
			Data []BlockData `json:"aaData"`
		}
		var blockDataResult BlockDataResult

		var sqlstr string
		sqlstr = fmt.Sprintf("SELECT count(*) FROM blocks WHERE channelname='%s'", channelid)
		//rows, err := fclient.blockDB.Query("SELECT count(*) FROM blocks WHERE channelname='mychannel'")
		fmt.Println(sqlstr)
		rows, err := fclient.blockDB.Query(sqlstr)
		if err != nil {
			fmt.Printf("error:%v", err)
			return ""
		}
		for rows.Next() {
			rows.Scan(&blockDataResult.TotalDisplayRecords)
		}
		rows.Close()
		//rows, err = fclient.blockDB.Query("SELECT blocknum, txcount, datahash, prehash FROM blocks WHERE channelname='mychannel'")
		sqlstr = fmt.Sprintf("SELECT blocknum, txcount, datahash, prehash FROM blocks WHERE channelname='%s' limit %d offset %d", channelid, displayLength, displayStart)
		fmt.Println(sqlstr)
		rows, err = fclient.blockDB.Query(sqlstr)
		if err != nil {
			fmt.Printf("error:%v", err)
			return ""
		}
		defer rows.Close()

		blockDataResult.TotalRecords = 0
		for rows.Next() {
			var blockData BlockData
			rows.Scan(&blockData.BlockNum, &blockData.TxCount, &blockData.DataHash, &blockData.PreHash)
			blockDataResult.Data = append(blockDataResult.Data, blockData)
			blockDataResult.TotalRecords += 1
		}

		resJson, _ := json.Marshal(blockDataResult)
		fmt.Println(string(resJson))
		return string(resJson)
}

func GetTxData(channelid string, displayStart int, displayLength int) string {
		type TxData struct {
			TxHash string `json:"txhash"`
			BlockNum int `json:"blocknum"`
			TxTime string `json:"txtime"`
			Chaincode string `json:"chaincode"`
		}
		type TxDataResult struct {
			TotalRecords int `json:"iTotalRecords"`
			TotalDisplayRecords int `json:"iTotalDisplayRecords"`
			Data []TxData `json:"aaData"`
		}
		var txDataResult TxDataResult

		var sqlstr string
		sqlstr = fmt.Sprintf("SELECT count(*) FROM transactions WHERE channelname='%s'", channelid)
		//rows, err := fclient.blockDB.Query("SELECT count(*) FROM blocks WHERE channelname='mychannel'")
		fmt.Println(sqlstr)
		rows, err := fclient.blockDB.Query(sqlstr)
		if err != nil {
			fmt.Printf("error:%v", err)
			return ""
		}
		for rows.Next() {
			rows.Scan(&txDataResult.TotalDisplayRecords)
		}
		rows.Close()
		//rows, err = fclient.blockDB.Query("SELECT blocknum, txcount, datahash, prehash FROM blocks WHERE channelname='mychannel'")
		sqlstr = fmt.Sprintf("SELECT txhash, blockid, createdate, chaincodename FROM transactions WHERE channelname='%s' limit %d offset %d", channelid, displayLength, displayStart)
		fmt.Println(sqlstr)
		rows, err = fclient.blockDB.Query(sqlstr)
		if err != nil {
			fmt.Printf("error:%v", err)
			return ""
		}
		defer rows.Close()

		txDataResult.TotalRecords = 0
		for rows.Next() {
			var txData TxData
			var txTime time.Time
			rows.Scan(&txData.TxHash, &txData.BlockNum, &txTime, &txData.Chaincode)
			txData.TxTime = txTime.Format("2006-01-02 15:04:05")
			txDataResult.Data = append(txDataResult.Data, txData)
			txDataResult.TotalRecords += 1
		}

		resJson, _ := json.Marshal(txDataResult)
		fmt.Println(string(resJson))
		return string(resJson)
}

func GetBlockCount(channelid string) string {
	var blockcount int
	var sqlstr string
	sqlstr = fmt.Sprintf("SELECT count(*) FROM blocks WHERE channelname='%s'", channelid)
	rows, err := fclient.blockDB.Query(sqlstr)
	if err != nil {
		fmt.Printf("error:%v", err)
		return "0"
	}
	for rows.Next() {
		rows.Scan(&blockcount)
	}
	rows.Close()
	return fmt.Sprintf("%d", blockcount)
}

func GetTxCount(channelid string) string {
	var txcount int
	var sqlstr string
	sqlstr = fmt.Sprintf("SELECT count(*) FROM transactions WHERE channelname='%s'", channelid)
	rows, err := fclient.blockDB.Query(sqlstr)
	if err != nil {
		fmt.Printf("error:%v", err)
		return "0"
	}
	for rows.Next() {
		rows.Scan(&txcount)
	}
	rows.Close()

	return fmt.Sprintf("%d", txcount)
}

func dealEvent() {
	regBlock, blockEventChannel, err := fclient.eventClient.RegisterBlockEvent()
	if err != nil {
		fmt.Printf("error registering for block events: %s\n", err)
	}
	defer fclient.eventClient.Unregister(regBlock)

	select {
	case blockEvent, ok := <-blockEventChannel:
		if !ok {
			fmt.Printf("unexpected closed channel\n")
		} else {
			fmt.Printf("blockEvent:%s\n", blockEvent.SourceURL)
		}
		//case <-time.After(5 * time.Second):
		//	fmt.Printf("timed out waiting for block event\n")
	}
}
func Init() error {
	blockDB, err := sql.Open("sqlite3", "/Users/liujian/code/gopath/src/liujian/fabcli/fabcli.db")
	if err != nil {
		fmt.Sprintf("open block db error\n")
		return err
	}
	fclient.blockDB = blockDB

	configProvider := config.FromFile("/Users/liujian/code/gopath/src/liujian/fabcli/config.yaml")
	//configBackend, err := config.FromFile("/Users/liujian/code/gopath/src/liujian/fabcli/config.yaml")()
	//configBackend, err := configProvider()
	//if err != nil {
	//	fmt.Sprintf("Failed to init config: %s", err)
	//	return err
	//}

	//configCryptoSuite, configEndpoint, configIdentity, err := config.FromBackend(configBackend)()
	//if err != nil {
	//	fmt.Sprintf("Failed to init config: %s", err)
	//	return err
	//}
	//configCryptoSuite := cryptosuite.ConfigFromBackend(configBackend)
	//configEndpoint := fab.ConfigFromBackend(configBackend)
	//configIdentity := msp.ConfigFromBackend(configBackend)

	//mychannelUser := selection.ChannelUser{ChannelID: "mychannel", Username: "admin", OrgName: "org1"}

	//sdk, err := fabsdk.New(config.FromFile("/Users/liujian/code/gopath/src/liujian/fabcli/config.yaml"))
	sdk, err := fabsdk.New(configProvider, fabsdk.WithLoggerPkg(&log.CliLoggingProvider{}))
		//fabsdk.WithServicePkg(&DynamicDiscoveryProviderFactory{})
		//fabsdk.WithServicePkg(&DynamicSelectionProviderFactory{ChannelUsers: []selection.ChannelUser{mychannelUser}}))
	if err != nil {
		fmt.Sprintf("Failed to init sdk: %s", err)
		return err
	}

	clientChannelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("admin"), fabsdk.WithOrg("org1"))

	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Sprintf("Failed to create new channel client: %s", err)
		return err
	}

	ledgerClient, err := ledger.New(clientChannelContext)
	if err != nil {
		fmt.Printf("Failed to create new ledger client: %s", err)
		return err
	}

	eventClient, err := event.New(clientChannelContext)
	if err != nil {
		fmt.Printf("Failed to create new event client: %s", err)
		return err
	}

	clientContext := sdk.Context(fabsdk.WithUser("admin"), fabsdk.WithOrg("org1"))
	resmgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s", err)
	}

	sysClientContext := sdk.Context(fabsdk.WithUser("admin"), fabsdk.WithOrg("ordererorg"))
	fmt.Println("====================>\n")

	sysResmgmtClient, err := resmgmt.New(sysClientContext)
	fmt.Println("====================>\n")
	if err != nil {
		fmt.Printf("error:%v\n", err)
		fmt.Sprintf("Failed to create new channel client: %s", err)
		return err
	}
	fclient.sysResmgmtClient = sysResmgmtClient

	/*
	sysClientChannelContext := sdk.ChannelContext("testchainid", fabsdk.WithOrg("ordererorg"))
	sysLedgerClient, err := ledger.New(sysClientChannelContext)
	if err != nil {
		fmt.Printf("Failed to create new ledger client: %s", err)
		return err
	}
	fclient.sysLedgerClient = sysLedgerClient
	*/

	//fclient.configBackend = configBackend
	//fclient.configCryptoSuite = configCryptoSuite
	//fclient.configEndpoint = configEndpoint
	//fclient.configIdentity = configIdentity
	fclient.sdk = sdk
	fclient.channelClient = channelClient
	fclient.ledgerClient = ledgerClient
	fclient.eventClient = eventClient
	fclient.resmgmtClient = resmgmtClient

	sdkContext := sdk.Context()
	ctx, _ := sdkContext()
	networkConfig, _ := ctx.EndpointConfig().NetworkConfig()
	//configBackends, _ := sdk.Config()
	//endpointConfig, _ := pkgfab.ConfigFromBackend(configBackends)
	//fmt.Println("----->")
	//networkConfig, _ := endpointConfig.NetworkConfig()
	fmt.Println(networkConfig.Name, networkConfig.Description, networkConfig.Version)
	fmt.Println("----->")
	for k, v := range networkConfig.Organizations {
		fmt.Println(k)
		fmt.Println(v.MSPID)
		for u, _ := range v.Users {
			fmt.Println("\t", u)
		}
	}
	fmt.Println("----->")
	for k, _ := range networkConfig.Channels {
		fmt.Println(k)
	}

	return nil
	//go dealEvent()

	//test config
	/*
	fmt.Println("----->")
	networkConfig, _ := configEndpoint.NetworkConfig()
	fmt.Println(networkConfig.Name, networkConfig.Description, networkConfig.Version)
	fmt.Println("----->")

	fmt.Println("networkConfig.Channels")
	for k, v := range networkConfig.Channels {
		fmt.Println("\tChannel:", k)
		for _, orderer := range v.Orderers {
			fmt.Println("\t\tOrderer:", orderer)
		}
		for peerName, _ := range v.Peers {
			fmt.Println("\t\tPeer:", peerName)
		}
		fmt.Println("\tChannelPolicies:")
		fmt.Println("\t\tMinResponses:", v.Policies.QueryChannelConfig.MinResponses)
		fmt.Println("\t\tMaxTargets:", v.Policies.QueryChannelConfig.MaxTargets)
	}
	fmt.Println("networkConfig.Organizations")
	for k, v := range networkConfig.Organizations {
		fmt.Println("\tOrganization:", k)
		fmt.Println("\t\tMSPID:", v.MSPID)
		fmt.Println("\t\tCryptoPath:", v.CryptoPath)
	}
	fmt.Println("networkConfig.Orderers")
	for k, v := range networkConfig.Orderers {
		fmt.Println("\tOrderer:", k)
		fmt.Println("\t\tURL:", v.URL)
	}
	fmt.Println("networkConfig.Peers")
	for k, v := range networkConfig.Peers {
		fmt.Println("\tPeer:", k)
		fmt.Println("\t\tURL:", v.URL)
		fmt.Println("\t\tEventURL:", v.EventURL)
	}
	*/

	return nil
}
