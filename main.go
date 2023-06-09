package main

import "face-restoration/internal/logic"

func main() {
	impl := logic.NewFaceRestorationImpl()
	impl.Engine.Run(":80")
}
