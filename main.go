package main

import "face-restoration/internal/logic"

func main() {
	impl := logic.NewFaceRestorationImpl()

	impl.Cron.Start()

	if err := impl.Engine.Run(":80"); err != nil {
		panic(err)
	}
}
