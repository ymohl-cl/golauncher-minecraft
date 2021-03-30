package installer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/ymohl-cl/golauncher-minecraft/pkg/version"
	"github.com/ymohl-cl/gonbt"
)

// os installer type
const (
	Windows = "windows"
	Linux   = "linux"
	OSX     = "darwin"
)

// resources minecraft
const (
	launcherProfileFile = "launcher_profiles.json"
	serversDatFile      = "servers.dat"
)

var (
	defaultOSXPathInstaller     = fmt.Sprintf(`/Users/%s/Library/Application Support/minecraft`, os.Getenv("USER"))
	defaultWindowsPathInstaller = fmt.Sprint(`%AppData%\.minecraft`)
)

// Installer minecraft mods
type Installer interface {
	Close() error
	UpdateLocation(location string) error
	MinecraftPath() Path
	ForgePath() Path
	Install() error

	ForgeVersion() string
	GameVersion() string
}

type installer struct {
	ftpDriver     *ftp.ServerConn
	minecraftPath Path
	forgePath     Path

	forgeVersion     string
	forgeFullVersion string
	serverName       string
	serverURL        string
	gameVersion      string
}

// Path type define a disk location and valid is true if location exist
type Path struct {
	Value string
	Valid bool
}

// Update record the path parameter and test if is it valid.
func (p *Path) Update(path string) error {
	var err error

	p.Value = path
	if _, err = os.Stat(p.Value); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		p.Valid = false
	} else {
		p.Valid = true
	}
	return nil
}

// New parse the config_installer.json file and create a new installer
func New() (Installer, error) {
	var err error
	var c Config
	var data []byte

	if data, err = ioutil.ReadFile("config_installer.json"); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return NewWithConfig(c)
}

// New statement of installer. Need to call the Close method
func NewWithConfig(conf Config) (Installer, error) {
	var i installer
	var err error

	i.forgeVersion = conf.Forge.Version
	i.serverName = conf.Server.Name
	i.serverURL = conf.Server.IP
	i.gameVersion = conf.GameVersion

	switch runtime.GOOS {
	case Windows:
		if err = i.UpdateLocation(defaultWindowsPathInstaller); err != nil {
			return nil, err
		}
	case OSX:
		if err = i.UpdateLocation(defaultOSXPathInstaller); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(errorOS)
	}

	if i.ftpDriver, err = ftp.Dial(conf.FTP.URL, ftp.DialWithTimeout(5*time.Second)); err != nil {
		return nil, err
	}
	if err = i.ftpDriver.Login(conf.FTP.Login, conf.FTP.Password); err != nil {
		return nil, err
	}
	return &i, nil
}

// GameVersion getter
func (i installer) GameVersion() string {
	return i.gameVersion
}

// ForgeVersion getter
func (i installer) ForgeVersion() string {
	return i.forgeVersion
}

// Close resources driver
func (i *installer) Close() error {
	if err := i.ftpDriver.Quit(); err != nil {
		return err
	}
	return nil
}

// UpdateInstallLocation : _
func (i *installer) UpdateLocation(location string) error {
	var err error

	i.forgePath = Path{}
	if err = i.minecraftPath.Update(location); err != nil {
		return err
	}
	if !i.minecraftPath.Valid {
		return nil
	}
	checker := filepath.Join(fmt.Sprintf("%s/versions", i.minecraftPath.Value))
	if _, err = os.Stat(checker); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		i.minecraftPath.Valid = false
		return nil
	}

	files := []os.FileInfo{}
	versionPath := filepath.Join(fmt.Sprintf("%s/versions", i.minecraftPath.Value))
	if files, err = ioutil.ReadDir(versionPath); err != nil {
		return err
	}

	// search the more recently version in the current version of forge
	forgeVersion := ""
	previousForgeVersion := ""
	for _, f := range files {
		if f.IsDir() && strings.HasPrefix(f.Name(), fmt.Sprintf("%s-forge", i.forgeVersion)) {
			if forgeVersion, err = version.Recent(forgeVersion, strings.TrimPrefix(f.Name(), fmt.Sprintf("%s-forge-", i.forgeVersion))); err != nil {
				return err
			}
			if previousForgeVersion != forgeVersion {
				previousForgeVersion = forgeVersion
				i.forgeFullVersion = f.Name()
			}
		}
	}
	if forgeVersion != "" {
		i.forgePath = Path{
			Value: filepath.Join(fmt.Sprintf("%s/versions/%s", i.minecraftPath.Value, i.forgeFullVersion)),
			Valid: true,
		}
	}
	return nil
}

// MinecraftPath getter
func (i installer) MinecraftPath() Path {
	return i.minecraftPath
}

// ForgePath getter
func (i installer) ForgePath() Path {
	return i.forgePath
}

// Install or update the server installation
func (i *installer) Install() error {
	var err error

	if !i.minecraftPath.Valid || !i.forgePath.Valid {
		return errors.New(errorPathInstallation)
	}
	if err = i.updateServersList(); err != nil {
		return err
	}
	// "mycube", "1.16.5-forge-36.1.3"
	if err = i.updateProfiles(); err != nil {
		return err
	}
	if err = i.updateMods(); err != nil {
		return err
	}
	return nil
}

func (i *installer) updateMods() error {
	var err error

	modFolderPath := filepath.Join(fmt.Sprintf("%s/mods", i.minecraftPath.Value))
	if _, err = os.Stat(modFolderPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		if err = os.Mkdir(modFolderPath, 755); err != nil {
			panic(err)
		}
	}

	var dir string
	var entries []*ftp.Entry
	if dir, err = i.ftpDriver.CurrentDir(); err != nil {
		return err
	}
	if entries, err = i.ftpDriver.List(dir); err != nil {
		return err
	}

	for _, e := range entries {
		if e.Type == ftp.EntryTypeFile {
			modPath := filepath.Join(modFolderPath, e.Name)
			if _, err = os.Stat(modPath); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}
			} else {
				// mod already exist
				continue
			}
			if err = func() error {
				var resp *ftp.Response
				if resp, err = i.ftpDriver.Retr(fmt.Sprintf("%s/%s", dir, e.Name)); err != nil {
					return err
				}
				defer resp.Close()

				var out *os.File
				if out, err = os.Create(modPath); err != nil {
					return err
				}
				defer out.Close()
				if _, err = io.Copy(out, resp); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *installer) updateServersList() error {
	var err error
	var tag gonbt.Tag
	var tagCompound *gonbt.CompoundT
	var tagList *gonbt.ListT
	var data []byte
	var ok bool

	// test file already exist
	serversDatFilePath := filepath.Join(i.minecraftPath.Value, serversDatFile)
	if _, err = os.Stat(serversDatFilePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		// new minecraft installation, need to create the servers.dat file
		return i.newServersList(serversDatFilePath)
	}

	// read and decode nbt data
	if data, err = ioutil.ReadFile(serversDatFilePath); err != nil {
		return err
	}
	if tag, err = gonbt.Unmarshal(data); err != nil {
		return err
	}

	if tagCompound, ok = tag.(*gonbt.CompoundT); !ok {
		return errors.New(errorServerDat)
	}

	for k, v := range tagCompound.Value {
		if k != "servers" {
			continue
		}
		if tagList, ok = v.(*gonbt.ListT); !ok {
			return errors.New(errorServerDat)
		}
	}
	if tagList == nil {
		return errors.New(errorServerDat)
	}
	find := false
	updated := false
	for _, v := range tagList.Value {
		var elem *gonbt.CompoundT
		if elem, ok = v.(*gonbt.CompoundT); !ok {
			return errors.New(errorServerDat)
		}
		var tagString *gonbt.StringT

		if tagString, ok = elem.Value["name"].(*gonbt.StringT); !ok {
			return errors.New(errorServerDat)
		}
		if tagString.Value == i.serverName {
			if tagString, ok = elem.Value["ip"].(*gonbt.StringT); !ok {
				return errors.New(errorServerDat)
			}
			find = true
			if tagString.Value != i.serverURL {
				tagString.Value = i.serverURL
				updated = true
			}
		}
	}
	if !find {
		elem := &gonbt.CompoundT{Value: make(map[string]interface{})}
		elem.Value["name"] = &gonbt.StringT{Value: i.serverName}
		elem.Value["ip"] = &gonbt.StringT{Value: i.serverURL}
		tagList.Value = append(tagList.Value, elem)
		updated = true
	}

	if updated {
		if data, err = gonbt.Marshal(tag, gonbt.CompressNone); err != nil {
			return err
		}
		if err = ioutil.WriteFile(serversDatFilePath, data, os.FileMode(os.O_RDWR)); err != nil {
			return err
		}
	}
	return nil
}

func (i *installer) newServersList(file string) error {
	var err error
	var tag gonbt.Tag
	var data []byte

	tag = &gonbt.CompoundT{
		Value: map[string]interface{}{
			"servers": &gonbt.ListT{
				Value: []interface{}{&gonbt.CompoundT{
					Value: map[string]interface{}{
						"name": &gonbt.StringT{Value: i.serverName},
						"ip":   &gonbt.StringT{Value: i.serverURL},
					},
				}},
			},
		},
	}

	if data, err = gonbt.Marshal(tag, gonbt.CompressNone); err != nil {
		return err
	}
	if err = ioutil.WriteFile(file, data, os.ModeAppend|0644); err != nil {
		return err
	}
	return nil
}

func (i *installer) updateProfiles() error {
	var err error
	var data []byte
	var l LauncherProfile

	filePath := filepath.Join(i.minecraftPath.Value, launcherProfileFile)
	if data, err = ioutil.ReadFile(filePath); err != nil {
		return err
	}
	if err = json.Unmarshal(data, &l); err != nil {
		return err
	}
	if p, ok := l.Profiles[i.serverName]; !ok || p.VersionID != i.forgeFullVersion {
		l.Profiles[i.serverName] = Profile{
			Icon:      "Lectern_Book",
			LastUsed:  time.Now(),
			VersionID: i.forgeFullVersion,
			Type:      "custom",
			Name:      i.serverName,
		}
		if data, err = json.Marshal(&l); err != nil {
			return err
		}
		if err = ioutil.WriteFile(filePath, data, os.FileMode(os.O_RDWR)); err != nil {
			return err
		}
	}
	return nil
}
