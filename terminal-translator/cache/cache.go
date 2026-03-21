package cache

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"crypto/sha256"
	"strings"
)

type CacheSchema interface {
	Init() error
	Exist(k string) bool
	Rd(k string) (string, error)
	Wt(k string, v string) error
}

type LocalCache struct {
	MetaPath	string
	HomeDir 	string
}

func (m *LocalCache) Init() error {
	cache_home, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("ERROR: Invalid cache directory")
		return err
	}

	exec, err := os.Executable()
	if err != nil {
		fmt.Println("ERROR: Invalid executable path")
		return err
	}

	prog, err := filepath.EvalSymlinks(exec)
	if err != nil {
		fmt.Println("ERROR: Evaluate symlink failed")
		return err
	}

	m.HomeDir = filepath.Join(cache_home, filepath.Base(prog))
	if err := os.Mkdir(m.HomeDir, os.ModePerm); err != nil && !os.IsExist(err) {
		fmt.Println("ERROR: Create cache directory failed")
		return err
	}

	m.MetaPath = filepath.Join(m.HomeDir, "meta")
	if _, err := os.Stat(m.MetaPath); os.IsNotExist(err) {
		if err := os.WriteFile(m.MetaPath, []byte{}, 0644); err != nil {
			fmt.Println("ERROR: Create meta file failed")
			return err
		}
	}

	return nil
}

func (m *LocalCache) Exist(k string) bool {
	kx := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(k))))

	cache_path := filepath.Join(m.HomeDir, hex.EncodeToString(kx[:4]))
	if _, err := os.Stat(cache_path); err != nil {
		return false
	}

	return true
}

func (m *LocalCache) Rd(k string) (string, error) {
	kx := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(k))))

	cache_path := filepath.Join(m.HomeDir, hex.EncodeToString(kx[:4]))
	if _, err := os.Stat(cache_path); err != nil {
		return "", nil
	}

	trans, err := os.ReadFile(cache_path)
	if err != nil {
		return "", err
	}

	return string(trans), nil
}

func (m *LocalCache) Wt(k string, v string) error {
	ck := strings.ToLower(strings.TrimSpace(k))
	kx := sha256.Sum256([]byte(ck))

	cache_path := filepath.Join(m.HomeDir, hex.EncodeToString(kx[:4]))
	if err := os.WriteFile(cache_path, []byte(v), 0644); err != nil {
		return err
	}

	mapping := fmt.Sprintf("<%s> <%s>\n", ck, hex.EncodeToString(kx[:4]))

	f, err := os.OpenFile(m.MetaPath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.WriteString(mapping); err != nil {
		return err
	}

	return nil
}

