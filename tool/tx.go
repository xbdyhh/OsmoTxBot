package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	atypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/osmosis-labs/osmosis/v11/app"
	"github.com/osmosis-labs/osmosis/v11/x/gamm/types"
	ptypes "github.com/osmosis-labs/osmosis/v11/x/gamm/types"
	"github.com/xbdyhh/OsmoTxBot/tool/module"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

var Ccontext = client.Context{}.WithChainID("osmosis-1")

func InitContext() {
	encodingConfig := app.MakeEncodingConfig()
	Ccontext = Ccontext.WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(atypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithViper("OSMOSIS")
	conf := sdk.GetConfig()
	conf.SetBech32PrefixForAccount("osmo", "osmopub")
}

func QueryAccountInfo(ctx context.Context, address string) (*module.AccountInfo, error) {
	myAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		"grpc.osmosis.zone:9090", // your gRPC server address.
		grpc.WithInsecure(),      // The SDK doesn't support any transport security mechanism.
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)

	// This creates a gRPC client to query the x/bank service.
	authClient := atypes.NewQueryClient(grpcConn)
	authRes, err := authClient.Account(ctx, &atypes.QueryAccountRequest{myAddress.String()})
	if err != nil {
		return nil, err
	}
	res, err := Ccontext.Codec.MarshalJSON(authRes)
	if err != nil {
		return nil, err
	}
	acc := &module.AccountInfo{}
	err = json.Unmarshal(res, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func QueryBalanceInfo(ctx context.Context, address string) (*module.Balances, error) {
	myAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		"grpc.osmosis.zone:9090", // your gRPC server address.
		grpc.WithInsecure(),      // The SDK doesn't support any transport security mechanism.
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)

	// This creates a gRPC client to query the x/bank service.
	authClient := banktypes.NewQueryClient(grpcConn)
	authRes, err := authClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
		Address:    myAddress.String(),
		Pagination: nil,
	})
	if err != nil {
		return nil, err
	}
	res, err := Ccontext.Codec.MarshalJSON(authRes)
	if err != nil {
		return nil, err
	}
	bal := &module.Balances{}
	err = json.Unmarshal(res, bal)
	if err != nil {
		return nil, err
	}
	return bal, nil
}

func QueryPoolInfo(ctx context.Context) (*module.Pools, error) {
	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		"grpc.osmosis.zone:9090", // your gRPC server address.
		grpc.WithInsecure(),      // The SDK doesn't support any transport security mechanism.
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)

	// This creates a gRPC client to query the x/bank service.
	poolClient := ptypes.NewQueryClient(grpcConn)
	authRes, err := poolClient.Pools(ctx, &ptypes.QueryPoolsRequest{
		&query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      1000,
			CountTotal: false,
			Reverse:    false,
		},
	})
	if err != nil {
		return nil, err
	}
	res, err := Ccontext.Codec.MarshalJSON(authRes)
	if err != nil {
		return nil, err
	}
	pool := &module.Pools{}
	err = json.Unmarshal(res, pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// NewBuildSwapExactAmountInMsg tokenInStr Expected format: "{amount}{denomination}"
func NewBuildSwapExactAmountInMsg(addr, tokenInStr, tokenOutMinAmtStr string, routerids, routerdenoms []string) (sdk.Msg, error) {
	routes, err := swapAmountInRoutes(routerids, routerdenoms)
	if err != nil {
		return nil, err
	}

	tokenIn, err := sdk.ParseCoinNormalized(tokenInStr)
	if err != nil {
		return nil, err
	}

	tokenOutMinAmt, ok := sdk.NewIntFromString(tokenOutMinAmtStr)
	if !ok {
		return nil, errors.New("invalid token out min amount")
	}
	msg := &types.MsgSwapExactAmountIn{
		Sender:            addr,
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmt,
	}

	return msg, nil
}
func swapAmountInRoutes(ids, denoms []string) ([]types.SwapAmountInRoute, error) {
	if len(ids) != len(denoms) {
		return nil, errors.New("swap route pool ids and denoms mismatch")
	}

	routes := []types.SwapAmountInRoute{}
	for index, poolIDStr := range ids {
		pID, err := strconv.Atoi(poolIDStr)
		if err != nil {
			return nil, err
		}
		routes = append(routes, types.SwapAmountInRoute{
			PoolId:        uint64(pID),
			TokenOutDenom: denoms[index],
		})
	}
	return routes, nil
}
