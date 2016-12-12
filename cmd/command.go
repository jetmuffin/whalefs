package cmd

import "time"

import (
	goflag "flag"
)

type flag struct {
	name 	string
	value 	string
	usage 	string
}

type Command struct {
	name 		string
	description 	string
	flags 		[]flag
	function	func(Flags)
}

type Flags interface {
	BoolVar(*bool, string, bool, string)
	String(string, string, string) *string
	Duration(string, time.Duration, string) *time.Duration
	Parse()
	Var(goflag.Value, string, string)
	Int(string, int, string) *int
}

