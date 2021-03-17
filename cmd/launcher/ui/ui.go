package ui

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/objects"
)

const (
	defaultIndexUI = 1
)

// UI for the golauncher-minecraft (implement the scene builder from game-builder)
type UI interface {
	Build() error
	Init() error
	IsInit() bool
	Run() error
	Stop()
	Close() error
	GetLayers() (map[uint8][]objects.Object, *sync.Mutex)
	KeyboardEvent(*sdl.KeyboardEvent)
	SetSwitcher(func(string) error)
	SetCloser(func(string) error)
	Update()
}

type ui struct {
	sync.Mutex
	layers   map[uint8][]objects.Object
	switcher func(string) error
	closer   func(string) error
}

// New return a new instance of ui
func New() (UI, error) {
	u := ui{
		layers: make(map[uint8][]objects.Object),
	}
	return &u, nil
}

// Build the scene
func (u *ui) Build() error {
	return nil
}

// Init the scene. Create static objects. Data is provide if you need.
func (u *ui) Init() error {
	return nil
}

// IsInit return status initialize
func (u ui) IsInit() bool {
	return true
}

// Run the scene
func (u *ui) Run() error {
	return nil
}

// Stop the scene, it's possible to Run later
func (u *ui) Stop() {}

// Close the scene at the end game
func (u *ui) Close() error {
	return nil
}

// GetLayers get objects list by layers order
func (u *ui) GetLayers() (map[uint8][]objects.Object, *sync.Mutex) {
	return u.layers, &u.Mutex
}

// KeyboardEvent provide key down to the scene
func (u *ui) KeyboardEvent(*sdl.KeyboardEvent) {}

// SetSwitcher can be called to change scene with index scene on
// first parameter and flag on true to close old scene and false to stop it only
func (u *ui) SetSwitcher(f func(string) error) {
	u.switcher = f
}

// SetCloser quit the application
func (u *ui) SetCloser(f func(string) error) {
	u.closer = f
}

// Update : called on each frame
func (u *ui) Update() {}
