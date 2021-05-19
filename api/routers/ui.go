package routers

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jweny/pocassist/web"
	"net/http"
	"strings"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) > len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

// BinaryFileSystem ...
func BinaryFileSystem(root string) *binaryFileSystem {
	return &binaryFileSystem{
		fs: &assetfs.AssetFS{
			Asset:     web.Asset,
			AssetDir:  web.AssetDir,
			AssetInfo: web.AssetInfo,
			Prefix:    root,
		},
	}
}
