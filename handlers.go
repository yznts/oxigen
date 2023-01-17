package main

import (
	"github.com/gorilla/mux"
	"github.com/kyoto-framework/kyoto/v2"
)

// This file includes common functions to work with handlers

// Page is a custom page registering function.
// Kyoto now supports only registering pages into default net/http mux,
// that's why we are using kyoto.HandlerPage instead.
// Think about this function as a shortcut for mux.HandleFunc(..., kyoto.HandlerPage(...)).
func Page[T any](mux *mux.Router, route string, page kyoto.Component[T]) {
	mux.HandleFunc(route, kyoto.HandlerPage(page))
}

// Action is a custom action registering function.
// Kyoto now supports only registering actions into default net/http mux,
// that's why we are using kyoto.HandlerAction instead.
// Also, gorilla/mux requires a bit different handling for actions (PathPrefix).
func Action[T any](mux *mux.Router, component kyoto.Component[T]) {
	pattern := kyoto.ActionConf.Path + kyoto.ComponentName(component) + "/"
	mux.PathPrefix(pattern).HandlerFunc(kyoto.HandlerAction(component))
}
