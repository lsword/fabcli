package httpserver

import (
	"fmt"
	//"io"
	//"io/ioutil"
	//"net"
	"strings"
	"strconv"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"liujian/fabcli/fabclient"

)

func InitRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/network.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "network.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/blocks.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "blocks.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/transactions.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transactions.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/chaincodes.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chaincodes.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/ajax/blockdata/:channelid", getBlockData)
	router.GET("/ajax/txdata/:channelid", getTxData)
	router.GET("/ajax/blockcount/:channelid", getBlockCount)
	router.GET("/ajax/txcount/:channelid", getTxCount)
	router.GET("/hello", getHello)
	router.GET("/channel/querypeers", getChannelPeers)
	router.GET("/channel/query/:chaincodename/:args", getChannelQuery)
	router.GET("/channel/execute/:chaincodename/:args", getChannelExecute)
	router.GET("/ledger/queryinfo", getLegerQueryInfo)
	router.GET("/ledger/queryconfig", getLegerQueryConfig)
	router.GET("/ledger/queryblock/:num", getLegerQueryBlock)
	router.GET("/ledger/queryblockbyhash/:blockhash", getLegerQueryBlockByHash)
	router.GET("/ledger/queryblockbytxid/:txid", getLegerQueryBlockByTxID)
	router.GET("/ledger/querytransaction/:txid", getLegerQueryTransaction)
	router.GET("/resmgmt/queryconfigfromorderer", getResmgmtQueryConfigFromOrderer)
	router.GET("/resmgmt/querychannels", getResmgmtQueryChannels)
	router.GET("/resmgmt/queryinstalledchaincodes", getResmgmtQueryInstalledChaincodes)
	router.GET("/resmgmt/queryinstantiatedchaincodes", getResmgmtQueryInstantiatedChaincodes)
	router.POST("/resmgmt/installcc", postResmgmtInstallCC)
	router.POST("/resmgmt/instantiatecc", postResmgmtInstantiateCC)
	router.POST("/resmgmt/upgradecc", postResmgmtUpgradeCC)
	return router
}

func getBlockData(c *gin.Context) {
	channelid := c.Param("channelid")
	displayStart, _ := strconv.Atoi(c.Query("iDisplayStart"))
	displayLength, _ := strconv.Atoi(c.Query("iDisplayLength"))
	result := fabclient.GetBlockData(channelid, displayStart, displayLength)
	c.String(http.StatusOK, result)
}

func getTxData(c *gin.Context) {
	channelid := c.Param("channelid")
	displayStart, _ := strconv.Atoi(c.Query("iDisplayStart"))
	displayLength, _ := strconv.Atoi(c.Query("iDisplayLength"))
	result := fabclient.GetTxData(channelid, displayStart, displayLength)
	c.String(http.StatusOK, result)
}

func getBlockCount(c *gin.Context) {
	channelid := c.Param("channelid")
	result := fabclient.GetBlockCount(channelid)
	c.String(http.StatusOK, result)
}

func getTxCount(c *gin.Context) {
	channelid := c.Param("channelid")
	result := fabclient.GetTxCount(channelid)
	c.String(http.StatusOK, result)
}

func getHello(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func getChannelPeers(c *gin.Context) {
	result := fabclient.ChannelQueryPeers()
	c.String(http.StatusOK, result)
}

func getChannelQuery(c *gin.Context) {
	chainCodeName := c.Param("chaincodename")
	args := c.Param("args")
	argsarray := strings.Split(args, ",")
	var queryArgs [][]byte
	for index := range argsarray {
		queryArgs = append(queryArgs, []byte(argsarray[index]))
	}

	result := fabclient.ChannelQuery(chainCodeName, queryArgs)
	c.String(http.StatusOK, result)
}

func getChannelExecute(c *gin.Context) {
	chainCodeName := c.Param("chaincodename")
	args := c.Param("args")
	argsarray := strings.Split(args, ",")
	var executeArgs [][]byte
	for index := range argsarray {
		executeArgs = append(executeArgs, []byte(argsarray[index]))
	}

	result := fabclient.ChannelExecute(chainCodeName, executeArgs)
	c.String(http.StatusOK, result)
}

func getLegerQueryInfo(c *gin.Context) {
	result := fabclient.LedgerQueryInfo()
	c.String(http.StatusOK, result)
}

func getLegerQueryConfig(c *gin.Context) {
	result := fabclient.LedgerQueryConfig()
	c.String(http.StatusOK, result)
}

func getLegerQueryBlock(c *gin.Context) {
	num, _ := strconv.Atoi(c.Param("num"))

	result := fabclient.LedgerQueryBlock(uint64(num))
	c.String(http.StatusOK, result)
}

func getLegerQueryBlockByHash(c *gin.Context) {
	blockhash := []byte(c.Param("blockhash"))

	result := fabclient.LedgerQueryBlockByHash(blockhash)
	c.String(http.StatusOK, result)
}

func getLegerQueryBlockByTxID(c *gin.Context) {
	txid := c.Param("txid")

	result := fabclient.LedgerQueryBlockByTxID(txid)
	c.String(http.StatusOK, result)
}

func getLegerQueryTransaction(c *gin.Context) {
	txid := c.Param("txid")

	result := fabclient.LedgerQueryTransaction(txid)
	c.String(http.StatusOK, result)
}

func getResmgmtQueryConfigFromOrderer(c *gin.Context) {
	result := fabclient.ResmgmtQueryConfigFromOrderer("")
	c.String(http.StatusOK, result)
}

func getResmgmtQueryChannels(c *gin.Context) {
	result := fabclient.ResmgmtQueryChannels()
	c.String(http.StatusOK, result)
}

func getResmgmtQueryInstalledChaincodes(c *gin.Context) {
	result := fabclient.ResmgmtQueryInstalledChaincodes()
	c.String(http.StatusOK, result)
}

func getResmgmtQueryInstantiatedChaincodes(c *gin.Context) {
	result := fabclient.ResmgmtQueryInstantiatedChaincodes()
	c.String(http.StatusOK, result)
}

func postResmgmtInstallCC(c *gin.Context) {
	name := c.PostForm("name")
	version := c.PostForm("version")
	path := c.PostForm("path")

	result := fabclient.ResmgmtInstallCC(name, version, path)
	c.String(http.StatusOK, result)
}

func postResmgmtInstantiateCC(c *gin.Context) {
	name := c.PostForm("name")
	version := c.PostForm("version")
	path := c.PostForm("path")
	channel := c.PostForm("channel")

	result := fabclient.ResmgmtInstantiateCC(name, version, path, channel)
	c.String(http.StatusOK, result)
}

func postResmgmtUpgradeCC(c *gin.Context) {
	name := c.PostForm("name")
	version := c.PostForm("version")
	path := c.PostForm("path")
	channel := c.PostForm("channel")

	result := fabclient.ResmgmtUpgradeCC(name, version, path, channel)
	c.String(http.StatusOK, result)
}

func Run() {
	err := fabclient.Init()
	if err != nil {
		fmt.Printf("fabclient init error:%v\n", err)
		return
	}

	//fabclient.Execute()
	//fabclient.Query()
	///fabclient.QueryBlock(41)
	///fabclient.QueryInstalledChaincodes()
	///fabclient.QueryBlockChain()
	///height, _ := fabclient.QueryBlockChainHeight("")
	//fabclient.SyncBlocks()

	go fabclient.SyncBlocks()

	router := InitRouter()

	go router.Run(fmt.Sprintf(":%d", 8888))
}
