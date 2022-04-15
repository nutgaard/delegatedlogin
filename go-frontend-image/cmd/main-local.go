package main

import (
	. "frontend-image/internal"
	. "frontend-image/internal/config"
)

func main() {
	// Only locally
	SetupEnv()

	StartApplication()
}
