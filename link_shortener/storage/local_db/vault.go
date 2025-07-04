package local_db

import (
	"encoding/json"
	"fmt"
	"io"
	"link_shortener/pkg/validate"
	"slices"
)

type Vault struct {
	Bins []Bin `json:"bins"`
}

type Db interface {
	Read() (io.ReadCloser, error)
	Write() (io.WriteCloser, error)
}

type VaultWithDb struct {
	Vault
	db Db
}

func NewVaultWithDb(db Db) (*VaultWithDb, error) {
	vault := &VaultWithDb{
		Vault: Vault{},
		db:    db,
	}
	err := vault.Read()
	if err != nil {
		return nil, err
	}
	return vault, nil
}

func (v *VaultWithDb) Create(request map[string]string) error {
	err := v.Read()
	if err != nil {
		return err
	}
	bin := NewBin(request["email"], request["hash"])
	err = validate.StructValidator(bin)
	if err != nil {
		return err
	}
	v.Bins = append(v.Bins, *bin)
	err = v.Update()
	if err != nil {
		return err
	}
	return nil
}

func (v *VaultWithDb) Read() (err error) {
	var file io.ReadCloser
	defer func() {
		if file != nil {
			err = file.Close()
		}
		if err != nil {
			err = fmt.Errorf("error in 'sync': %w", err)
		}
	}()
	file, err = v.db.Read()
	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(&v.db)
	if err != nil {
		return err
	}
	return nil
}

func (v *VaultWithDb) Update() (err error) {
	var file io.WriteCloser
	defer func() {
		if file != nil {
			err = file.Close()
		}
		if err != nil {
			err = fmt.Errorf("error in 'update': %w", err)
		}
	}()
	file, err = v.db.Write()
	if err != nil {
		return err
	}
	err = json.NewEncoder(file).Encode(v.Bins)
	if err != nil {
		return err
	}
	return nil
}

func (v *VaultWithDb) Delete(hash string) error {
	err := v.Read()
	if err != nil {
		return err
	}
	for i, bin := range v.Bins {
		if bin.Hash == hash {
			v.Bins = slices.Delete(v.Bins, i, i+1)
			break
		}
	}
	err = v.Update()
	if err != nil {
		return err
	}
	return nil
}
