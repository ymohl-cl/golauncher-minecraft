module github.com/ymohl-cl/golauncher-minecraft

go 1.15

require (
	github.com/beito123/binary v3.0.1+incompatible // indirect
	github.com/beito123/nbt v1.2.2
	github.com/jlaffaye/ftp v0.0.0-20210307004419-5d4190119067
	github.com/sqweek/dialog v0.0.0-00010101000000-000000000000
	github.com/veandco/go-sdl2 v0.4.5
	github.com/ymohl-cl/game-builder v1.0.2
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/tools v0.1.0 // indirect
)

replace github.com/ymohl-cl/game-builder => ../game-builder

replace github.com/sqweek/dialog => ../../sqweek/dialog
