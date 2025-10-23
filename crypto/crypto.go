package crypto

import (
	"github.com/lengzhao/govm/types"
)

// Crypto 加密接口
type Crypto interface {
	// GenerateKeyPair 生成密钥对
	GenerateKeyPair() (PrivateKey, PublicKey, error)

	// Sign 使用私钥签名数据
	Sign(data []byte, privateKey PrivateKey) ([]byte, error)

	// Verify 使用公钥验证签名
	Verify(data []byte, signature []byte, publicKey PublicKey) bool

	// Hash 计算数据哈希
	Hash(data []byte) types.Hash

	// GenerateAddress 从公钥生成地址
	GenerateAddress(publicKey PublicKey) types.Address
}

// PrivateKey 私钥接口
type PrivateKey interface {
	// Bytes 获取私钥字节表示
	Bytes() []byte

	// PublicKey 获取对应的公钥
	PublicKey() PublicKey

	// Sign 使用私钥签名数据
	Sign(data []byte) ([]byte, error)
}

// PublicKey 公钥接口
type PublicKey interface {
	// Bytes 获取公钥字节表示
	Bytes() []byte

	// Address 获取公钥对应的地址
	Address() types.Address

	// Verify 使用公钥验证签名
	Verify(data []byte, signature []byte) bool
}
