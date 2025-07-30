package helpers

import (
	cryptorand "crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func HashPassword(password string) (string, error) {
	params := &ArgonParams{
		memory:      19 * 1024,
		iterations:  2,
		parallelism: 1,
		saltLength:  16,
		keyLength:   32,
	}

	// 1. Random salt
	salt := make([]byte, params.saltLength)
	if _, err := cryptorand.Read(salt); err != nil {
		return "", err
	}

	// 2. Hasilkan hash menggunakan Argon2id.
	hash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	// 3. Gabungkan semua informasi menjadi satu string untuk disimpan di database.
	// Format: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.memory, params.iterations, params.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func CheckPasswordHashWithArgon2(password, encodedHash string) (bool, error) {
	// 1. Uraikan (parse) hash yang tersimpan untuk mendapatkan parameter, salt, dan hash asli.
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("format hash tidak valid")
	}

	var version int
	params := &ArgonParams{}
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return false, fmt.Errorf("versi argon2 tidak kompatibel")
	}

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return false, fmt.Errorf("gagal mem-parse parameter hash")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("gagal decode salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("gagal decode hash")
	}
	params.keyLength = uint32(len(hash))

	// 2. Buat hash baru dari password yang diinput user menggunakan parameter yang sama.
	comparisonHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	// 3. Bandingkan kedua hash dengan aman (constant-time comparison).
	// Ini penting untuk mencegah timing attack.
	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return true, nil
	}
	return false, nil
}
