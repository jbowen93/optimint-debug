package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/celestiaorg/go-cnc"
	"github.com/celestiaorg/optimint/da"
	"github.com/celestiaorg/optimint/types"
	pb "github.com/celestiaorg/optimint/types/pb/optimint"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/gogo/protobuf/proto"
)

const bridgeNode = "http://35.208.160.145:26658"
const fullNode = "http://146.148.83.114:26658"

type DataAvailabilityLayerClient struct {
	client *cnc.Client

	namespaceID [8]byte
	// config      Config
	logger log.Logger
}

func main() {
	startHeight, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	endHeight, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	bridgeClient, err := cnc.NewClient(bridgeNode)
	if err != nil {
		panic("failed to setup bridgeNode cnc client")
	}
	bridgeDAClient := DataAvailabilityLayerClient{
		bridgeClient,
		[8]byte{0, 1, 2, 3, 4, 5, 6, 8},
		log.NewTMLogger(os.Stdout),
	}

	fullClient, err := cnc.NewClient(fullNode)
	if err != nil {
		panic("failed to setup fullNode cnc client")
	}
	fullDAClient := DataAvailabilityLayerClient{
		fullClient,
		[8]byte{0, 1, 2, 3, 4, 5, 6, 8},
		log.NewTMLogger(os.Stdout),
	}
	for i := startHeight; i <= endHeight; i++ {
		bridgeBlocks := bridgeDAClient.RetrieveBlocks(uint64(i))
		for _, block := range bridgeBlocks.Blocks {
			fmt.Printf("daHeight: %v, bridgeHeight: %v, ", i, block.Header.Height)
		}
		fullBlocks := fullDAClient.RetrieveBlocks(uint64(i))
		for _, block := range fullBlocks.Blocks {
			fmt.Printf("fullHeight: %v\n", block.Header.Height)
		}
	}
}

func (c *DataAvailabilityLayerClient) RetrieveBlocks(dataLayerHeight uint64) da.ResultRetrieveBlocks {
	data, err := c.client.NamespacedData(context.TODO(), c.namespaceID, dataLayerHeight)
	if err != nil {
		return da.ResultRetrieveBlocks{
			DAResult: da.DAResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	blocks := make([]*types.Block, len(data))
	for i, msg := range data {
		var block pb.Block
		err = proto.Unmarshal(msg, &block)
		if err != nil {
			c.logger.Error("failed to unmarshal block", "daHeight", dataLayerHeight, "position", i, "error", err)
			continue
		}
		blocks[i] = new(types.Block)
		err := blocks[i].FromProto(&block)
		if err != nil {
			return da.ResultRetrieveBlocks{
				DAResult: da.DAResult{
					Code:    da.StatusError,
					Message: err.Error(),
				},
			}
		}
	}

	return da.ResultRetrieveBlocks{
		DAResult: da.DAResult{
			Code:     da.StatusSuccess,
			DAHeight: dataLayerHeight,
		},
		Blocks: blocks,
	}
}
