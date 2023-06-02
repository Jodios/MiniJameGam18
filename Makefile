GAME_NAME=game

build:
	GOOS=js GOARCH=wasm go build -o public/${GAME_NAME}.wasm github.com/jodios/minijamegame18

run: build
	./${GAME_NAME}

clean:
	go clean
	rm public/${GAME_NAME}.wasm
