package main

import (
	"fmt"

	"ghost_escape/game"
	"ghost_escape/game/core"
)

func main() {
	sceneTitle := &game.SceneTitle{}
	game := core.GetInstance()
	if err := game.Init("GhostEscape", 1280, 720, sceneTitle); err != nil {
		fmt.Println(err)
		return
	}
	game.Run()
	game.Clean()
}
