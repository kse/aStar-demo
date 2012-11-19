SRC = src.go
OUT = display

all: build

build:
	go build -o $(OUT) $(SRC)

debug:
	go build -ldflags "-s" -o $(OUT) $(SRC)

clean:
	rm -f $(OUT)
