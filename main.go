package main

import (
	"assignment2/routers"
)

const PORT = ":8000"

func main() {
	routers.StartServer().Run(PORT)
}
