package crypto

import (
    "bytes"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "math/big"
    "os"
    "runtime/debug"
)

var (
    ErrDataToLarge     = errors.New("message too long for RSA public key size")
    ErrDataLen         = errors.New("data length error")
    ErrDataBroken      = errors.New("data broken, first byte is not zero")
    ErrKeyPairDismatch = errors.New("data is not encrypted by the private key")
    ErrDecryption      = errors.New("decryption error")
    ErrPublicKey       = errors.New("public key error")
    ErrPrivateKey      = errors.New("private key error")
)
// rsa 加解密
type rsaCrypto struct {
    // pem 格式公钥
    publicKey []byte

    // pem 格式私钥
    privateKey []byte

    // rsa 公钥
    rsaPriKey *rsa.PrivateKey

    // rsa 私钥
    rsaPubKey *rsa.PublicKey
}

// 创建 rsa 实例
// 公钥 pubKey 和 私钥 priKey 必须传一个，没值的传 nil
// 加解密时公私钥必须是一对
func NewRsa(pubKey, priKey []byte) *rsaCrypto {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
            debug.PrintStack()
            os.Exit(-2)
        }
    }()

    if len(pubKey) == 0 && len(priKey) == 0 {
        panic("public key or private key is needed")
    }
    rc := &rsaCrypto{
        publicKey:  pubKey,
        privateKey: priKey,
    }

    var err error
    if len(pubKey) > 0 {
        rc.rsaPubKey, err = getRsaPublicKey(pubKey)
        if err != nil {
            panic(err.Error())
        }
    }
    if len(priKey) > 0 {
        rc.rsaPriKey, err = getRsaPrivateKey(priKey)
        if err != nil {
            panic(err.Error())
        }
    }
    return rc
}

// rsa 公钥加密
func (r *rsaCrypto) EncryptWithPublicKey(data []byte) ([]byte, error) {
    // 解密 pem 格式公钥
    block, _ := pem.Decode(r.publicKey)
    if block == nil {
        return nil, ErrPublicKey
    }

    // 解析公钥
    pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    pub := pubInterface.(*rsa.PublicKey)
    return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// rsa 私钥解密
func (r *rsaCrypto) DecryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
    block, _ := pem.Decode(r.privateKey)
    if block == nil {
        return nil, ErrPrivateKey
    }

    pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
        if err != nil {
            return nil, err
        }
        pri = p.(*rsa.PrivateKey)
    }
    return rsa.DecryptPKCS1v15(rand.Reader, pri, ciphertext)
}

// rsa 私钥加密
func (r *rsaCrypto) EncryptWithPrivateKey(data []byte) ([]byte, error) {
    out := bytes.NewBuffer(nil)
    err := r.privKeyIO(bytes.NewReader(data), out)
    if err != nil {
        return nil, err
    }
    return ioutil.ReadAll(out)
}

// rsa 公钥解密
func (r *rsaCrypto) DecryptWithPublicKey(ciphertext []byte) ([]byte, error) {
    out := bytes.NewBuffer(nil)
    err := r.pubKeyIO(bytes.NewReader(ciphertext), out)
    if err != nil {
        return nil, err
    }
    return ioutil.ReadAll(out)
}

// 公钥解密 reader
func (r *rsaCrypto) pubKeyIO(in io.Reader, w io.Writer) (err error) {
    k := (r.rsaPubKey.N.BitLen() + 7) / 8
    buf := make([]byte, k)
    var b []byte
    size := 0
    for {
        size, err = in.Read(buf)
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        if size < k {
            b = buf[:size]
        } else {
            b = buf
        }
        b, err = r.pubKeyDecrypt(b)
        if err != nil {
            return err
        }
        if _, err = w.Write(b); err != nil {
            return err
        }
    }
    return nil
}

// 私钥加密 reader
func (r *rsaCrypto) privKeyIO(re io.Reader, w io.Writer) (err error) {
    k := (r.rsaPriKey.N.BitLen()+7)/8 - 11
    buf := make([]byte, k)
    var b []byte
    size := 0
    for {
        size, err = re.Read(buf)
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        if size < k {
            b = buf[:size]
        } else {
            b = buf
        }
        b, err = r.priKeyEncrypt(rand.Reader, b)
        if err != nil {
            return err
        }
        if _, err = w.Write(b); err != nil {
            return err
        }
    }
    return nil
}

// 私钥加密
func (r *rsaCrypto) priKeyEncrypt(rand io.Reader, hashed []byte) ([]byte, error) {
    hl := len(hashed)
    k := (r.rsaPriKey.N.BitLen() + 7) / 8
    if k < hl+11 {
        return nil, ErrDataLen
    }
    em := make([]byte, k)
    em[1] = 1
    for i := 2; i < k-hl-1; i++ {
        em[i] = 0xff
    }
    copy(em[k-hl:k], hashed)
    m := new(big.Int).SetBytes(em)
    c, err := decrypt(rand, r.rsaPriKey, m)
    if err != nil {
        return nil, err
    }
    copyWithLeftPad(em, c.Bytes())
    return em, nil
}

// 公钥解密
func (r *rsaCrypto) pubKeyDecrypt(data []byte) ([]byte, error) {
    k := (r.rsaPubKey.N.BitLen() + 7) / 8
    if k != len(data) {
        return nil, ErrDataLen
    }
    m := new(big.Int).SetBytes(data)
    if m.Cmp(r.rsaPubKey.N) > 0 {
        return nil, ErrDataToLarge
    }
    m.Exp(m, big.NewInt(int64(r.rsaPubKey.E)), r.rsaPubKey.N)
    d := leftPad(m.Bytes(), k)
    if d[0] != 0 {
        return nil, ErrDataBroken
    }
    if d[1] != 0 && d[1] != 1 {
        return nil, ErrKeyPairDismatch
    }
    var i = 2
    for ; i < len(d); i++ {
        if d[i] == 0 {
            break
        }
    }
    i++
    if i == len(d) {
        return nil, nil
    }
    return d[i:], nil
}

// 获取 rsa 私钥
func getRsaPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode(privateKey)
    if block == nil {
        return nil, ErrPrivateKey
    }
    pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err == nil {
        return pri, nil
    }
    p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return p.(*rsa.PrivateKey), nil
}

// 设置 rsa 公钥
func getRsaPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
    block, _ := pem.Decode(publicKey)
    if block == nil {
        return nil, ErrPublicKey
    }
    // x509 parse public key
    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return pub.(*rsa.PublicKey), nil
}

// 从 crypto/rsa 复制
var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

// 从 crypto/rsa 复制
func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
    if c.Cmp(priv.N) > 0 {
        err = ErrDecryption
        return
    }
    var ir *big.Int
    if random != nil {
        var r *big.Int

        for {
            r, err = rand.Int(random, priv.N)
            if err != nil {
                return
            }
            if r.Cmp(bigZero) == 0 {
                r = bigOne
            }
            var ok bool
            ir, ok = modInverse(r, priv.N)
            if ok {
                break
            }
        }
        bigE := big.NewInt(int64(priv.E))
        rpowe := new(big.Int).Exp(r, bigE, priv.N)
        cCopy := new(big.Int).Set(c)
        cCopy.Mul(cCopy, rpowe)
        cCopy.Mod(cCopy, priv.N)
        c = cCopy
    }
    if priv.Precomputed.Dp == nil {
        m = new(big.Int).Exp(c, priv.D, priv.N)
    } else {
        m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
        m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
        m.Sub(m, m2)
        if m.Sign() < 0 {
            m.Add(m, priv.Primes[0])
        }
        m.Mul(m, priv.Precomputed.Qinv)
        m.Mod(m, priv.Primes[0])
        m.Mul(m, priv.Primes[1])
        m.Add(m, m2)

        for i, values := range priv.Precomputed.CRTValues {
            prime := priv.Primes[2+i]
            m2.Exp(c, values.Exp, prime)
            m2.Sub(m2, m)
            m2.Mul(m2, values.Coeff)
            m2.Mod(m2, prime)
            if m2.Sign() < 0 {
                m2.Add(m2, prime)
            }
            m2.Mul(m2, values.R)
            m.Add(m, m2)
        }
    }
    if ir != nil {
        m.Mul(m, ir)
        m.Mod(m, priv.N)
    }

    return
}

// 从 crypto/rsa 复制
func copyWithLeftPad(dest, src []byte) {
    numPaddingBytes := len(dest) - len(src)
    for i := 0; i < numPaddingBytes; i++ {
        dest[i] = 0
    }
    copy(dest[numPaddingBytes:], src)
}

// 从 crypto/rsa 复制
func leftPad(input []byte, size int) (out []byte) {
    n := len(input)
    if n > size {
        n = size
    }
    out = make([]byte, size)
    copy(out[len(out)-n:], input)
    return
}

// 从 crypto/rsa 复制
func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
    g := new(big.Int)
    x := new(big.Int)
    y := new(big.Int)
    g.GCD(x, y, a, n)
    if g.Cmp(bigOne) != 0 {
        return
    }
    if x.Cmp(bigOne) < 0 {
        x.Add(x, n)
    }
    return x, true
}
