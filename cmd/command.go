package cmd

import (
	goflag "flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type flag struct {
	name  string
	value string
	usage string
}

type command struct {
	name        string
	description string
	flags       []flag
	function    func(Flags)
}

type Flags interface {
	BoolVar(*bool, string, bool, string)
	String(string, string, string) *string
	Duration(string, time.Duration, string) *time.Duration
	Parse()
	Var(goflag.Value, string, string)
	Int(string, int, string) *int
}

type AppConfig struct {
	whoami      string
	help        bool
	globalFlags []flag
	global      func(Flags)
	commands    []command
}

func App() *AppConfig {
	app := new(AppConfig)
	app.Command("help", "Show this message", func(f Flags) { f.Parse(); app.Usage() })
	return app
}

type flagDummy struct {
	list *[]flag
}

func (f *flagDummy) BoolVar(_ *bool, name string, value bool, usage string) {
	flag := flag{name, fmt.Sprintf("%+v", value), usage}
	*f.list = append(*f.list, flag)
}

func (f *flagDummy) Duration(name string, value time.Duration, usage string) *time.Duration {
	flag := flag{name, fmt.Sprintf("%+v", value), usage}
	*f.list = append(*f.list, flag)
	return nil
}

func (f *flagDummy) String(name string, value string, usage string) *string {
	flag := flag{name, fmt.Sprintf("%#v", value), usage}
	*f.list = append(*f.list, flag)
	return nil
}
func (f *flagDummy) Int(name string, value int, usage string) *int {
	flag := flag{name, fmt.Sprintf("%+v", value), usage}
	*f.list = append(*f.list, flag)
	return nil
}
func (f *flagDummy) Var(value goflag.Value, name string, usage string) {
	flag := flag{name, value.String(), usage}
	*f.list = append(*f.list, flag)
}

type doneTracing struct{}

func (f *doneTracing) Error() string {
	return "Tracing flags"
}

func (f *flagDummy) Parse() {
	panic(doneTracing{})
}

func (app *AppConfig) Global(f func(Flags)) {
	if app.global != nil {
		panic("Already set")
	}
	app.global = f
	f(&flagDummy{&app.globalFlags})
}

func (app *AppConfig) Command(name string, description string, f func(Flags)) {
	var flags []flag
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case doneTracing:
			default:
				panic(r)
			}
		}
		c := command{name, description, flags, f}
		app.commands = append(app.commands, c)
	}()

	f(&flagDummy{&flags})
}

type flagSet struct {
	*goflag.FlagSet
	app *AppConfig
}

func newFlagSet(app *AppConfig) *flagSet {
	return &flagSet{goflag.NewFlagSet("", goflag.ContinueOnError), app}
}

type flagSetFailure struct {
	err error
}

func (f *flagSet) Parse() {
	if err := f.FlagSet.Parse(os.Args); err != nil {
		switch err {
		case goflag.ErrHelp:
			f.app.Usage()
		default:
			fmt.Println("run with command 'help' for usage information")
			os.Exit(2)
		}
	}
	if len(f.FlagSet.Args()) > 0 {
		fmt.Println("arguments provided but not defined:", strings.Join(f.FlagSet.Args(), " "))
		fmt.Println("run with command 'help' for usage information")
		os.Exit(2)
	}
}

func (app *AppConfig) Run() {
	app.whoami = path.Base(os.Args[0])
	endOfGlobalFlags := 0
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-") {
			endOfGlobalFlags++
		} else {
			break
		}
	}
	commandArgs := os.Args[endOfGlobalFlags+1:]
	os.Args = os.Args[1 : endOfGlobalFlags+1]
	set := newFlagSet(app)
	set.FlagSet.Usage = func() {}
	app.global(set)
	set.Parse()
	if len(commandArgs) == 0 {
		app.Usage()
	}
	command := commandArgs[0]
	os.Args = commandArgs[1:]
	for _, c := range app.commands {
		if c.name == command {
			set := newFlagSet(app)
			set.FlagSet.Usage = func() {}
			c.function(set)
			os.Exit(0)
		}
	}
	app.Usage()
}

func (app *AppConfig) Usage() {
	fmt.Println("Usage:", app.whoami, "[global flags]", "command", "[flags]")
	fmt.Println("Global:")
	for _, f := range app.globalFlags {
		fmt.Printf("\t-%s=%v\n", f.name, f.value)
	}
	fmt.Println("Commands:")
	for _, c := range app.commands {
		fmt.Printf("\t%s: %s\n", c.name, c.description)
		for _, f := range c.flags {
			fmt.Printf("\t\t-%s=%s\n", f.name, f.value)
		}
	}

	os.Exit(2)
}