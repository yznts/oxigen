package pages

import (
	"go.kyoto.codes/v3/component"
)

type HomeState struct {
	common
}

func Home(ctx *component.Context) component.State {
	state := &HomeState{}
	return state
}
