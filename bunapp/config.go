package bunapp

import (
	"embed"
	"io/fs"
	"sync"
)
var (
	//go:embed embed
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

func FS() fs.FS {
	unwrapFSOnce.Do(func() {
		fsys, err := fs.Sub(embedFS, "embed")
		if err != nil {
			panic(err)
		}
		unwrappedFS = fsys
	})
	return unwrappedFS
}


type AppConfig	struct {
	Dev bool
	Port int
	Db struct {
		Host string
		Port int
		User string
		Password string
		Database string
	}
	Jwt struct {
		Secret string
		Expiration string
	}
}

func ReadConfig(fsys fs.FS, service, env string) (*AppConfig, error) {}

