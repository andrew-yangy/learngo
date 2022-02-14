package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ddvkid/learngo/controllers/user"
	"github.com/ddvkid/learngo/internal/ephemeraldb"
	"github.com/ddvkid/learngo/internal/storage"
	"github.com/ddvkid/learngo/internal/storage/postgres"
	"github.com/ddvkid/learngo/internal/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	// BlueprintConnectionForTest grants all the tests in the package to clone the "blueprint"
	// database to create copies of the test dataset. See: common/ephemeraldb/README.md
	// Do not use this value in the production code
	BlueprintConnectionForTest ephemeraldb.EphemeralConnection = ephemeraldb.NewInvalid()
)

const (
	defaultContextTimeout = 10 * time.Second
	ethKey1               = "0xdc6cb44b02f57fe6660595d95d9217f9cfdec57e"
	ethKey2               = "0xdc6cb44b02f57fe6660595d95d9217f9cfdec57f"
)

func TestMain(m *testing.M) {
	postgres.WithEphemeralDB(m, &BlueprintConnectionForTest)
}

func TestListStarkKeysNoKey(t *testing.T) {
	ctx := context.Background()
	store := postgres.LocalConnection(t, BlueprintConnectionForTest)
	_, err := store.InsertAccount(ctx, &storage.Account{EtherKey: ethKey1, StarkKey: "0x024b928189711f0b0c08ffc83b17a55bef11dbc596f844beb00ecbe1d6053c95", Nonce: 1, TxHash: "sometxhash"})
	assert.NoError(t, err)
	_, err = store.InsertAccount(ctx, &storage.Account{EtherKey: ethKey2, StarkKey: "0x034b928189711f0b0c08ffc83b17a55bef11dbc596f844beb00ecbe1d6053c95", Nonce: 2, TxHash: "sometxhash"})
	assert.NoError(t, err)

	req := events.APIGatewayProxyRequest{
		Body: `{
			"ether_key": "0x01234"
		}`, /*no matching accounts*/
	}
	res, err := handle(ctx, req, store)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestListStarkKeysSingleKey(t *testing.T) {
	ctx := context.Background()
	store := postgres.LocalConnection(t, BlueprintConnectionForTest)
	_, err := store.InsertAccount(ctx, &storage.Account{EtherKey: ethKey1, StarkKey: "0x024b928189711f0b0c08ffc83b17a55bef11dbc596f844beb00ecbe1d6053c95", Nonce: 1, TxHash: "sometxhash"})
	assert.NoError(t, err)
	_, err = store.InsertAccount(ctx, &storage.Account{EtherKey: ethKey2, StarkKey: "0x034b928189711f0b0c08ffc83b17a55bef11dbc596f844beb00ecbe1d6053c95", Nonce: 2, TxHash: "sometxhash"})
	assert.NoError(t, err)

	req := events.APIGatewayProxyRequest{
		Body: `{
			"ether_key": "0xdc6cb44b02f57fe6660595d95d9217f9cfdec57e"
		}`, /*one matching account, same as ethKey1*/
	}
	res, err := handle(ctx, req, store)
	assert.NoError(t, err)

	var ur user.Response
	err = util.DecodeResponse(res.Body, &ur)
	assert.ElementsMatch(t, ur.StarkKeys, []string{"0x024b928189711f0b0c08ffc83b17a55bef11dbc596f844beb00ecbe1d6053c95"})
}
