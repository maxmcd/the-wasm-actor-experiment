package main

import (
	"fmt"
	"io/ioutil"
	"log"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	// Almost all operations in wasmtime require a contextual `store`
	// argument to share, so create that first
	store := wasmtime.NewStore(wasmtime.NewEngine())

	b, err := ioutil.ReadFile("./wasm.wasm")
	if err != nil {
		return err
	}

	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModule(store.Engine, b)
	if err != nil {
		return err
	}

	fmt.Println(module.Imports()[0].Name())

	fdWrite := wasmtime.WrapFunc(store, func(int32, int32, int32, int32) int32 {
		return 0
	})
	ticks := wasmtime.WrapFunc(store, func() int64 {
		return 0
	})
	// Next up we instantiate a module which is where we link in all our
	// imports. We've got one import so we pass that in here.
	instance, err := wasmtime.NewInstance(store, module, []*wasmtime.Extern{
		fdWrite.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
		ticks.AsExtern(),
	})
	if err != nil {
		return err
	}

	// After we've instantiated we can lookup our `run` function and call
	// it.
	run := instance.GetExport("run").Func()
	_, err = run.Call()
	if err != nil {
		return err
	}
	return nil

}
