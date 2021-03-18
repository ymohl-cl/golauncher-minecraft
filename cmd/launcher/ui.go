package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/beito123/nbt"
	"github.com/jlaffaye/ftp"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/objects"
)

const (
	lancherProfileFile = "launcher_profile.json"
	serverFile         = "servers.dat"
)

const (
	// order layers of scene
	layerBackground = iota
	layer1
	layer2
	layer3
	layer4
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

// Install data to install lancher minecraft and theirs components
type Install struct {
	location string
	exec     string
	mod      string
}

type ui struct {
	config Config
	sync.Mutex
	layers   map[uint8][]objects.Object
	switcher func(string) error
	closer   func(string) error

	renderer *sdl.Renderer
	init     bool
	install  Install
}

// NewUI return a new instance of ui
func NewUI(c Config, r *sdl.Renderer) (UI, error) {
	u := ui{
		config:   c,
		layers:   make(map[uint8][]objects.Object),
		renderer: r,
	}
	return &u, nil
}

// Init the scene. Create static objects. Data is provide if you need.
func (u *ui) Init() error {
	u.init = true
	return nil
}

// IsInit return status initialize
func (u *ui) IsInit() bool {
	return u.init
}

/* url to hierarchy folder https://minecraft.gamepedia.com/.minecraft
OS 	Location
Windows 	%APPDATA%\.minecraft
macOS 	~/Library/Application Support/minecraft
Linux 	~/.minecraft
*/

// Run the scene
func (u *ui) Run() error {
	var err error

	if err = u.installData(); err != nil {
		if err != os.ErrNotExist {
			return err
		}
		u.install.location = ""
	} else {
		u.install.mod = fmt.Sprintf("%s/%s", u.install.location, "mods")
	}

	stream, err := nbt.FromFile(fmt.Sprintf("%s/%s", u.install.location, serverFile), nbt.BigEndian)
	if err != nil {
		return err
	}

	tag, err := stream.ReadTag()
	if err != nil {
		return err
	}
	var c *nbt.Compound
	fmt.Printf("tag id: %s\n", string(tag.ID()))
	switch tag.(type) {
	case *nbt.Compound:
		c = tag.(*nbt.Compound)
		fmt.Printf("C parsed: %v\n", c)
	default:
		fmt.Printf("not parsed\n")
	}

	str, err := tag.ToString()
	if err != nil {
		return err
	}
	fmt.Printf("data: %s\n", str)
	data := &nbt.Compound{
		Value: map[string]nbt.Tag{"servers": &nbt.List{
			ListType: nbt.IDTagCompound,
			Value: []nbt.Tag{
				&nbt.Compound{
					Value: map[string]nbt.Tag{
						"ip":   &nbt.String{Value: "minecraft1065.omgserv.com:10003"},
						"name": &nbt.String{Value: "bouziers"},
					},
				},
				&nbt.Compound{
					Value: map[string]nbt.Tag{
						"ip":   &nbt.String{Value: "minecraft1049.omgserv.com:10016"},
						"name": &nbt.String{Value: "mycube"},
					},
				},
			},
		}}}
	stream2 := nbt.NewStream(nbt.BigEndian)
	if err = stream2.WriteTag(data); err != nil {
		return err
	}
	buff, err := nbt.Compress(stream2, nbt.CompressGZip, nbt.DefaultCompressionLevel)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("./servers-mycube.dat", buff, os.ModePerm)
	if err != nil {
		panic(err)
	}
	// check os
	// check le dossier minecraft
	// IF OK
	// Check exec available
	// Check folder mods and create it if need

	// IF NOT
	// Get path from user
	return nil
}

func (u *ui) installData() error {
	var err error
	//	var info os.FileInfo

	switch runtime.GOOS {
	case "windows":
	case "linux":
	case "darwin":
		u.install.location = fmt.Sprintf(`/Users/%s/Library/Application Support/minecraft`, os.Getenv("USER"))
		if _, err = os.Stat(u.install.location); err != nil {
			return err
		}
	default:
		return errors.New("os not supported")
	}
	return nil
}

// Stop the scene, it's possible to Run later
func (u *ui) Stop() {}

// Close the scene at the end game
func (u *ui) Close() error {
	var err error

	u.init = false
	if err = objects.Closer(u.layers); err != nil {
		return err
	}
	u.layers = make(map[uint8][]objects.Object)
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

// actionInstall
func (u *ui) actionInstall(data ...interface{}) {
	var c *ftp.ServerConn
	var err error

	if c, err = ftp.Dial("62.210.46.77:21", ftp.DialWithTimeout(5*time.Second)); err != nil {
		panic(err)
	}
	if err = c.Login("326314_launcher", "1234"); err != nil {
		panic(err)
	}

	if _, err = os.Stat(u.install.mod); err != nil {
		if err != os.ErrNotExist {
			panic(err)
		}
		if err = os.Mkdir(u.install.mod, 755); err != nil {
			panic(err)
		}
	}
	var dir string
	if dir, err = c.CurrentDir(); err != nil {
		panic(err)
	}
	var entries []*ftp.Entry
	if entries, err = c.List(dir); err != nil {
		panic(err)
	}
	/*
		r, err := c.Retr("test-file.txt")
		if err != nil {
			panic(err)
		}
		defer r.Close()

		buf, err := ioutil.ReadAll(r)
		println(string(buf))
	*/
	for _, e := range entries {
		if e.Type == ftp.EntryTypeFile {
			if _, err = os.Stat(fmt.Sprintf("%s/%s", u.install.mod, e.Name)); err != nil {
				if err != os.ErrNotExist {
					panic(err)
				}
			} else {
				continue
			}
			func() {
				var resp *ftp.Response
				if resp, err = c.Retr(fmt.Sprintf("%s/%s", dir, e.Name)); err != nil {
					panic(err)
				}
				defer resp.Close()

				var buf []byte
				if buf, err = ioutil.ReadAll(resp); err != nil {
					panic(err)
				}
				if err = ioutil.WriteFile(fmt.Sprintf("%s/%s", u.install.mod, e.Name), buf, os.FileMode(os.O_RDWR)); err != nil {
					panic(err)
				}
			}()
		}
		fmt.Printf("%s type of %s\n", e.Name, e.Type)
	}
	fmt.Printf("acces dir %s\n", dir)

	if err := c.Quit(); err != nil {
		panic(err)
	}

}
