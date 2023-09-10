package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgUpdateParams int = 100

	OpWeightMsgUpdateParams = "op_weight_msg_update_params" // #nosec
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs() []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateClientMsgUpdateParams,
		),
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateConnectionMsgUpdateParams,
		),
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgRecoverClient,
		),
	}
}

// SimulateClientMsgUpdateParams returns a random MsgUpdateParams for 02-client
func SimulateClientMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	var signer sdk.AccAddress = address.Module("gov")
	params := types.DefaultParams()
	params.AllowedClients = []string{"06-solomachine", "07-tendermint"}

	return &types.MsgUpdateParams{
		Signer: signer.String(),
		Params: params,
	}
}

// SimulateConnectionMsgUpdateParams returns a random MsgUpdateParams 03-connection
func SimulateConnectionMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	var signer sdk.AccAddress = address.Module("gov")
	params := connectiontypes.DefaultParams()
	params.MaxExpectedTimePerBlock = uint64(100)

	return &connectiontypes.MsgUpdateParams{
		Signer: signer.String(),
		Params: params,
	}
}

// SimulateRecoverClient returns a random MsgRecoverClient 02-client
func SimulateMsgRecoverClient(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	var signer sdk.AccAddress = address.Module("gov")

	subjectClientId := "07-tendermint-1"
	substituteClientId := "07-tendermint-2"

	return &clienttypes.MsgRecoverClient{
		Signer:             signer.String(),
		SubjectClientId:    subjectClientId,
		SubstituteClientId: substituteClientId,
	}
}

// SimulateRecoverClient returns a random MsgRecoverClient 02-client
func SimulateMsgIBCSoftwareUpgrade(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	var signer sdk.AccAddress = address.Module("gov")

	plan := upgradetypes.Plan{
		Name:   "upgrade IBC clients",
		Height: 1000,
	}

	anyClient, err := clienttypes.PackClientState(&ibctm.ClientState{})
	if err != nil {
		return nil
	}

	return &clienttypes.MsgIBCSoftwareUpgrade{
		Signer:              signer.String(),
		Plan:                plan,
		UpgradedClientState: anyClient,
	}

}
