SRC = main.go field.go drawing.go
OUT = astar

all: build

build:
	go build -o $(OUT) $(SRC)

debug:
	go build -ldflags "-s" -o $(OUT) $(SRC)

clean:
	rm -f $(OUT)
