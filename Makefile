build:
	GOOS=js GOARCH=wasm go build -o public/game.wasm github.com/jodios/minijamegame18
run: build
