package dfusesquadron

import (
	"context"
	"fmt"
	"io"
	"os"

	dfuse "github.com/dfuse-io/client-go"
	"github.com/dfuse-io/dgrpc"
	pbbstream "github.com/dfuse-io/pbgo/dfuse/bstream/v1"
	"github.com/golang/protobuf/ptypes"
	pbcodec "github.com/streamingfast/streamingfast-client/pb/dfuse/ethereum/codec/v1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func ReadBlock(blockstart int64, blockend uint64) {
	apiKey := "server_7d59808b3e874fd899e7f785e546642a"
	endpoint := "heco.streamingfast.io:443"
	cursor := ""
	filter := "true"
	var dialOptions []grpc.DialOption

	dfuse, err := dfuse.NewClient("api.streamingfast.io", apiKey)
	noError(err, "unable to create streamingfast client")

	conn, err := dgrpc.NewExternalClient(endpoint, dialOptions...)
	noError(err, "unable to create external gRPC client")

	streamClient := pbbstream.NewBlockStreamV2Client(conn)

stream:
	for {
		tokenInfo, err := dfuse.GetAPITokenInfo(context.Background())
		noError(err, "unable to retrieve StreamingFast API token")

		forkSteps := []pbbstream.ForkStep{pbbstream.ForkStep_STEP_NEW}

		credentials := oauth.NewOauthAccess(&oauth2.Token{AccessToken: tokenInfo.Token, TokenType: "Bearer"})
		stream, err := streamClient.Blocks(context.Background(), &pbbstream.BlocksRequestV2{
			StartBlockNum:     blockstart,
			StartCursor:       cursor,
			StopBlockNum:      blockend,
			ForkSteps:         forkSteps,
			IncludeFilterExpr: filter,
			Details:           pbbstream.BlockDetails_BLOCK_DETAILS_FULL,
		}, grpc.PerRPCCredentials(credentials))
		noError(err, "unable to start blocks stream")

		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break stream
				}

				fmt.Printf("Stream encountered a remote error, going to retry\n\tcursor:%s\n\terror:%v\n", cursor, err)
				break
			}

			block := &pbcodec.Block{}
			err = ptypes.UnmarshalAny(response.Block, block)
			noError(err, "should have been able to unmarshal received block payload")

			cursor = response.Cursor
			fmt.Println(block.Number)
		}
	}
}

func noError(err error, message string, args ...interface{}) {
	if err != nil {
		quit(message+": "+err.Error(), args...)
	}
}

func quit(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	os.Exit(1)
}
