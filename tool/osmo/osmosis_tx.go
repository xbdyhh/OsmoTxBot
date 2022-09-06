package osmo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	ctx "github.com/cosmos/cosmos-sdk/client/tx"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	asigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	atypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/osmosis-labs/osmosis/v11/app"
	"github.com/osmosis-labs/osmosis/v11/x/gamm/types"
	ptypes "github.com/osmosis-labs/osmosis/v11/x/gamm/types"
	"github.com/xbdyhh/OsmoTxBot/tool"
	"github.com/xbdyhh/OsmoTxBot/tool/module"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

const (
	GAS_LIMIT           = 100000
	GRPC_SERVER_ADDRESS = "grpc.osmosis.zone:9090"
	CHAIN_ID            = "osmosis-1"
	ACCOUNT_ADDR        = "osmo16kydz6vznpgtpgws733panrs6atdsefcfxa97j"
	GAS_FEE             = 1000
)

var Ccontext = client.Context{}.WithChainID(CHAIN_ID)

func InitCcontext() {
	encodingConfig := app.MakeEncodingConfig()
	Ccontext = Ccontext.WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(atypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithViper("OSMOSIS").
		WithSignModeStr(flags.SignModeDirect)
	conf := sdk.GetConfig()
	conf.SetBech32PrefixForAccount("osmo", "osmopub")
}
func QueryOsmoAccountInfo(ctx *tool.MyContext, address string) (*module.AccountInfo, error) {
	myAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
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

func QueryOsmoBalanceInfo(ctx *tool.MyContext, address string) (*module.Balances, error) {
	myAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
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

func QueryOsmoPoolInfo(ctx *tool.MyContext) (*module.Pools, error) {
	// Create a connection to the gRPC server.
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
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
func NewOsmoBuildSwapExactAmountInMsg(addr string, tokenIn sdk.Coin, tokenOutMinAmtStr string, routerids, routerdenoms []string) (sdk.Msg, error) {
	routes, err := swapAmountInRoutes(routerids, routerdenoms)
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

func SignTx(txBuilder client.TxBuilder, priv ctypes.PrivKey, sequence uint64, accountNum uint64, msg sdk.Msg) error {
	var err error
	if err = txBuilder.SetMsgs(msg); err != nil {
		return err
	}
	sigV2 := signing.SignatureV2{
		PubKey: priv.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  Ccontext.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return err
	}
	signerData := asigning.SignerData{
		ChainID:       CHAIN_ID,
		AccountNumber: accountNum,
		Sequence:      sequence,
	}
	sigV2, err = ctx.SignWithPrivKey(Ccontext.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, priv, Ccontext.TxConfig, sequence)
	if err != nil {
		return err
	}
	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return err
	}
	return nil
}

func SendOsmoTx(ctx *tool.MyContext, mnemonic, tokenInDemon, tokenOutMinAmtStr string, tokenInAmount int64, sequence, accnum uint64,
	routerids, routerdenoms []string) (*sdk.TxResponse, error) {
	txBuilder := Ccontext.TxConfig.NewTxBuilder()
	txBuilder.SetGasLimit(GAS_LIMIT)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("uosmo", GAS_FEE)))
	priv, err := tool.NewPrivateKeyByMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	tokenInStr := sdk.NewInt64Coin(tokenInDemon, tokenInAmount)
	msg, err := NewOsmoBuildSwapExactAmountInMsg(ACCOUNT_ADDR, tokenInStr, tokenOutMinAmtStr, routerids, routerdenoms)
	if err != nil {
		return nil, err
	}
	if err = SignTx(txBuilder, priv, sequence, accnum, msg); err != nil {
		return nil, err
	}
	txBytes, err := Ccontext.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	txJSONBytes, err := Ccontext.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	fmt.Println(string(txJSONBytes))
	res, err := tool.BrocastTransaction(ctx, GRPC_SERVER_ADDRESS, txBytes)
	return res, nil
}

func SendOsmoTx2(ctx *tool.MyContext, Ccontext *client.Context, mnemonic string, accAddr string, sequence uint64, accnum uint64,
	tokenInDemon, tokenOutMinAmtStr string, tokenInAmount int64, routerids, routerdenoms []string) (*sdk.TxResponse, error) {
	priv, err := tool.NewPrivateKeyByMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	tokenInStr := sdk.NewInt64Coin(tokenInDemon, tokenInAmount)
	msg, err := NewOsmoBuildSwapExactAmountInMsg(accAddr, tokenInStr, tokenOutMinAmtStr, routerids, routerdenoms)
	if err != nil {
		return nil, err
	}
	msg2, err := NewOsmoBuildSwapExactAmountInMsg(accAddr, tokenInStr, tokenOutMinAmtStr, routerids, routerdenoms)
	if err != nil {
		return nil, err
	}
	return tool.SignAndBrocastTransaction(ctx, Ccontext, priv, sequence, accnum, 300000, "uosmo", 2500, "", "osmosis-1", GRPC_SERVER_ADDRESS, msg, msg2)
}

func PrivateToOsmoAddress(prv ctypes.PrivKey) (string, error) {
	addr, err := sdk.AccAddressFromHex(prv.PubKey().Address().String())
	if err != nil {
		return "", err
	}
	addr2, err := bech32.ConvertAndEncode("osmo", addr)
	if err != nil {
		return "", err
	}
	return addr2, nil
}
