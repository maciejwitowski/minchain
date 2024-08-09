package core

import (
	"minchain/core/types"
)

type Chainhead interface {
	SetHead(block *types.Block)
	GetHead() *types.Block
}
