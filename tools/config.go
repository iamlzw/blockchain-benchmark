package tools

const (
	ccID="benchmark4"
	ccPath="D:\\workspace\\go\\src\\github.com\\lifegoeson\\benchmark-chaincode"
	configPath="D:\\workspace\\go\\src\\github.com\\lifegoeson\\blockchain-benchmark\\config\\config_e2e.yaml"
	peer1 = "peer0.org1.example.com"
	peer2 = "peer0.org2.example.com"
	orderer = "orderer.example.com"
	channelID      = "mychannel"
	org1Name        = "Org1"
	org2Name        = "Org2"
	orgAdmin       = "Admin"
)
var initArgs = [][]byte{[]byte("a"), []byte("100"), []byte("b"), []byte("200")}
