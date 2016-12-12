package cmd

import (
	"testing"
)

func TestApp(t *testing.T) {
	cli := App()

	// global var
	var global_var bool
	cli.Global(func(flag Flags) {
		flag.BoolVar(&global_var, "global_var", false, "global_var")
	})

	cli.Command("test command", "test command description", func(flag Flags) {
		var boolFlag *bool

		stringFlag := flag.String("test flag", "default flag", "test flag description")
		flag.BoolVar(boolFlag, "test flag", true, "test flag description")
		intFlag := flag.Int("test flag", 1, "test flag description")

		flag.Parse()

		if *stringFlag != "default flag" || *boolFlag != true || *intFlag != 1 {
			t.Error("parse flag error.")
		}
	})
}
