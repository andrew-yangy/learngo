package ephemeraldb

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

/* #nosec */
func RandomNumericString() string {
	if u, err := uuid.NewRandom(); err != nil {
		return u.String()
	} else {
		// this error technically never happens but since it is a pure function we have to find an
		// alternative method to generate the random string; we fall back to using the rand package.
		rand.Seed(time.Now().UnixNano())
		return fmt.Sprintf("%d", rand.Uint64())
	}
}
