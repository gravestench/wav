package wav

import (
	"github.com/gravestench/wav/pkg"
)

func WavDecompress(data []byte, channelCount int) ([]byte, error) {
	return pkg.WavDecompress(data, channelCount)
}

func HuffmanDecompress(data []byte) []byte {
	return pkg.HuffmanDecompress(data)
}
