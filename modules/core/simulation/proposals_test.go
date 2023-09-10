package simulation_test

import (
	"math/rand"
	"testing"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v7/modules/core/simulation"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func TestProposalMsgs(t *testing.T) {
	// initialize parameters
	s := rand.NewSource(1)
	r := rand.New(s)

	ctx := sdk.NewContext(nil, cmtproto.Header{}, true, nil)
	accounts := simtypes.RandomAccounts(r, 3)

	// execute ProposalMsgs function
	weightedProposalMsgs := simulation.ProposalMsgs()
	require.Equal(t, 4, len(weightedProposalMsgs))

	w0 := weightedProposalMsgs[0]

	// tests w0 interface:
	require.Equal(t, simulation.OpWeightMsgUpdateParams, w0.AppParamsKey())
	require.Equal(t, simulation.DefaultWeightMsgUpdateParams, w0.DefaultWeight())

	msg := w0.MsgSimulatorFn()(r, ctx, accounts)
	msgUpdateParams, ok := msg.(*types.MsgUpdateParams)
	require.True(t, ok)

	require.Equal(t, sdk.AccAddress(address.Module("gov")).String(), msgUpdateParams.Signer)
	require.EqualValues(t, []string{"06-solomachine", "07-tendermint"}, msgUpdateParams.Params.AllowedClients)

	w1 := weightedProposalMsgs[1]

	// tests w1 interface:
	require.Equal(t, simulation.OpWeightMsgUpdateParams, w1.AppParamsKey())
	require.Equal(t, simulation.DefaultWeightMsgUpdateParams, w1.DefaultWeight())

	msg1 := w1.MsgSimulatorFn()(r, ctx, accounts)
	msgUpdateConnectionParams, ok := msg1.(*connectiontypes.MsgUpdateParams)
	require.True(t, ok)

	require.Equal(t, sdk.AccAddress(address.Module("gov")).String(), msgUpdateParams.Signer)
	require.EqualValues(t, uint64(100), msgUpdateConnectionParams.Params.MaxExpectedTimePerBlock)

	w2 := weightedProposalMsgs[2]

	// tests w2 interface:
	require.Equal(t, simulation.OpWeightMsgUpdateParams, w2.AppParamsKey())
	require.Equal(t, simulation.DefaultWeightMsgUpdateParams, w2.DefaultWeight())

	msg2 := w2.MsgSimulatorFn()(r, ctx, accounts)
	msgRecoverClient, ok := msg2.(*clienttypes.MsgRecoverClient)
	require.True(t, ok)

	require.Equal(t, sdk.AccAddress(address.Module("gov")).String(), msgRecoverClient.Signer)
	require.EqualValues(t, "07-tendermint-1", msgRecoverClient.SubjectClientId)
	require.EqualValues(t, "07-tendermint-2", msgRecoverClient.SubstituteClientId)

	w3 := weightedProposalMsgs[3]

	// tests w3 interface:
	require.Equal(t, simulation.OpWeightMsgUpdateParams, w3.AppParamsKey())
	require.Equal(t, simulation.DefaultWeightMsgUpdateParams, w3.DefaultWeight())

	msg3 := w3.MsgSimulatorFn()(r, ctx, accounts)
	msgIBCSoftwareUpgrade, ok := msg3.(*clienttypes.MsgIBCSoftwareUpgrade)
	require.True(t, ok)

	plan := upgradetypes.Plan{
		Name:   "upgrade IBC clients",
		Height: 1000,
	}

	anyClient, err := clienttypes.PackClientState(&ibctm.ClientState{})
	require.NoError(t, err)

	require.Equal(t, sdk.AccAddress(address.Module("gov")).String(), msgIBCSoftwareUpgrade.Signer)
	require.EqualValues(t, plan, msgIBCSoftwareUpgrade.Plan)
	require.EqualValues(t, anyClient, msgIBCSoftwareUpgrade.UpgradedClientState)

}
