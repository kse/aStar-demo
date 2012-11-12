all: graphics

graphics: world
	go build -o display graphics.go

world:
	go build -o world.a world.go
