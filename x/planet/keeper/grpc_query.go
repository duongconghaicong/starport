package keeper

import (
	"github.com/hello/planet/x/planet/types"
)

var _ types.QueryServer = Keeper{}
