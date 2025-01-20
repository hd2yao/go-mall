package resources

import (
	"bytes"
	"embed"
	"io"
)

// Go 程序打包的时候只会打 .go 文件，这种静态资源文件需要我们使用 go embed 功能才能打包到二进制文件里去x

//go:embed *
var f embed.FS

// LoadResourceFile 把资源文件 filePath 加载到内存，以 io.Reader 的形式返回
func LoadResourceFile(filePath string) (io.Reader, error) {
	_bytes, err := f.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(_bytes), nil
}
