### 1. 加解密说明
由于非对称加密算法的运行速度比对称加密算法的速度慢很多，当需要加密大量的数据时，建议采用对称加密算法，提高加解密速度。

对称加密算法不能实现签名，所以签名只能非对称算法。
由于对称加密算法的密钥管理是一个复杂的过程，密钥的管理直接决定着他的安全性，因此当数据量很小时，我们可以考虑采用非对称加密算法。  

在实际的操作过程中，通常采用的方式是：用**非对称加密算法**管理**对称算法的密钥**，然后**用对称加密算法加密数据**，这样就集成了两类加密算法的优点，既实现了加密速度快的优点，又实现了安全方便管理密钥的优点。

常用的对称加密算法 **AES**（AES 可以被硬件支持），非对称加密算法 RSA、**ECC**(算法各项指标都优于 RSA)。

### 2. 常见哈希和加解密方式：
#### 2.1 哈希 Hash
- **FNV Hash**
- adler
- CRC32/64
- (Message-Digest Algorithm 信息-摘要算法)
    - MD4
    - **MD5**
- SHA (Secure Hash Algorithm 安全散列算法)
    - **SHA-1**
    - SHA-2 (SHA-224、**SHA-256**、SHA-384、SHA-512)

MD5 和 SHA1 已经被攻破，安全性高的地方建议采用 SHA-256。

#### 2.2 加解密 Encryption
- 2.2.1 对称加密 (Symmetric)
    - 凯撒密码（Caesar cipher 变换加密）
    - **AES (Advanced Encryption Standard)**
    - SM1、SM4
    - **DES**、3DES、Blowfish、RC2、RC4、RC5、RC6 等
- 2.2.2 非对称加密 (Asymmetry)
    - 基于因数分解的算法 (**RSA**：Ron Rivest、Adi Shamir、Leonard Adleman 人名)
    - 基于整数有限域离散对数的算法（DSA：Digital Signature Algorithm）
    - 基于椭圆曲线的算法 (**ECC**：Elliptic Curve Cryptography)
    - 商密/国密2（SM2）
