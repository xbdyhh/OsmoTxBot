package tool

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	ctx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	asigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"google.golang.org/grpc"
)

func BuildAndSignTX(Ccontext *client.Context, priv ctypes.PrivKey, sequence uint64, accountNum uint64, chainId string, gasLimit uint64, denom string, fee int64, memo string, msg ...sdk.Msg) ([]byte, error) {
	txBuilder := Ccontext.TxConfig.NewTxBuilder()
	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin(denom, fee)))
	txBuilder.SetMemo(memo)
	err := txBuilder.SetMsgs(msg...)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	signerData := asigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accountNum,
		Sequence:      sequence,
	}
	sigV2, err = ctx.SignWithPrivKey(Ccontext.TxConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, priv, Ccontext.TxConfig, sequence)
	if err != nil {
		return nil, err
	}
	if err = txBuilder.SetSignatures(sigV2); err != nil {
		return nil, err
	}
	txBytes, err := Ccontext.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return txBytes, nil

}

func BrocastTransaction(ctx *MyContext, GRPC_SERVER_ADDRESS string, txBytes []byte) (*sdk.TxResponse, error) {
	grpcConn, err := grpc.Dial(
		GRPC_SERVER_ADDRESS, // Or your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return nil, err
	}
	txClient := tx.NewServiceClient(grpcConn)
	grpcRes, err := txClient.BroadcastTx(
		ctx,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(grpcRes.TxResponse)
	ctx.Logger.Infof("%v", grpcRes.TxResponse)
	defer grpcConn.Close()
	return grpcRes.TxResponse, nil
}

func SignAndBrocastTransaction(ctx *MyContext,
	Ccontext *client.Context, priv ctypes.PrivKey, sequence uint64, accountNum uint64,
	gasLimit uint64, denom string, fee int64, memo string,
	chainId string, GRPC_SERVER_ADDRESS string, msg ...sdk.Msg) (*sdk.TxResponse, error) {
	txBytes, err := BuildAndSignTX(Ccontext, priv, sequence, accountNum, chainId, gasLimit, denom, fee, memo, msg...)
	if err != nil {
		return nil, err
	}
	return BrocastTransaction(ctx, GRPC_SERVER_ADDRESS, txBytes)
}

func NewPrivateKeyByMnemonic(mnemonic string) (ctypes.PrivKey, error) {
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyring.SigningAlgoList{hd.Secp256k1})
	if err != nil {
		return nil, err
	}

	// create master key and derive first key for keyring
	prvbz, err := algo.Derive()(mnemonic, "", hd.CreateHDPath(118, 0, 0).String())
	if err != nil {
		return nil, err
	}
	prv := algo.Generate()(prvbz)
	return prv, nil
}
