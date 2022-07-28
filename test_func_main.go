package main

import (
	"github.com/y-omicron/util/Log"
	"github.com/y-omicron/util/Util"
)

func main() {
	Log.New(Log.LevelInfo, true, "util.log")
	Log.Info("Info!")
	Log.Success("Success!")
	Log.Debug("this is debug！")
	Log.Error("error!")
	Log.Fatal("gg!")
	Log.Trace("Trace")
	Log.Warning("Warn!")
	Util.OpenFileToWrite("test.log", []byte("hello world!"))
	Util.RandomInt(3000, 4000)
	Util.RandString(4)
}
