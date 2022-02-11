package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ddvkid/learngo/internal/ephemeraldb"
	"github.com/ddvkid/learngo/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	// BlueprintConnectionForTest grants all the tests in the package to clone the "blueprint"
	// database to create copies of the test dataset. See: common/ephemeraldb/README.md
	// Do not use this value in the production code
	BlueprintConnectionForTest ephemeraldb.EphemeralConnection = ephemeraldb.NewInvalid()
)

func TestMain(m *testing.M) {
	postgres.WithEphemeralDB(m, &BlueprintConnectionForTest)
}

func TestPayload(t *testing.T) {
	ctx := context.Background()
	store := postgres.LocalConnection(t, BlueprintConnectionForTest)
	req := events.APIGatewayProxyRequest{
		Body: `{
				"ether_key": "0x3c9fdd0863cd6d5df75dda12e3f55d6d44c07dc0",
				"stark_key": "0x04e6df51094f4b134a45059fd8becca127ee958228fa3dadd36a09c4e0010102",
				"nonce": 0,
				"stark_signature": "0x0302ed3176d1dc92eb8b47253eb83d0afab068f45ba84dd36d38cee7e210d053065ed5711c421c46650ab7b81ea1bd37cee7d3b07ebbe47f1dac6692c50ebe1e"
			}`,
	}
	response, err := handle(ctx, req, store)
	assert.NoError(t, err)
	assert.Equal(t, "{\"transaction_hash\":\"\"}", response.Body)
}
