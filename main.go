package main

import (
	"./streamer"
)

func main() {
	stream := streamer.NewStreamer("discipline.ogg")
	stream.Stream(8000)
}
