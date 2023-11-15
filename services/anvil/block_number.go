package anvil

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

func blockNumber(client *rpc.Client) (hexutil.Uint64, error) {
	var blockNumber hexutil.Uint64
	if err := client.Call(&blockNumber, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil
}
