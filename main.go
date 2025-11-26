package main

import (
	"fmt"

	"ghost_escape/game"
	"ghost_escape/game/core"
)

func main() {
	sceneMain := &game.SceneMain{}
	game := core.GetInstance()
	if err := game.Init("GhostEscape", 1280, 720, sceneMain); err != nil {
		fmt.Println(err)
		return
	}
	game.Run()
	game.Clean()
}
