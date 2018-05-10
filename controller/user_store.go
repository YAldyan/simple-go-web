package controller

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type FileUserStore struct {
	filename string
	Users    map[string]User
}

var GlobalUserStore UserStore

// func init() {
// 	/*
// 		membaca data user yang telah pernah dibuat
// 		dan disimpan dalam users.json
// 	*/
// 	store, err := NewFileUserStore("./data/users.json")
// 	if err != nil {
// 		panic(fmt.Errorf("Error creating user store: %s", err))
// 	}
// 	GlobalUserStore = store
// }

func NewFileUserStore(filename string) (*FileUserStore, error) {

	/*
		inisialisasi store dengan reference/pointer
		FileUserStore.
	*/
	store := &FileUserStore{
		Users:    map[string]User{},
		filename: filename,
	}

	/*
		read content dari file json
	*/
	contents, err := ioutil.ReadFile(filename)

	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}

	/*
		konversio data content dari json file
		ke variabel store
	*/
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (store FileUserStore) Save(user User) error {
	/*
		data yang di-insert disimpan dalam
		variabel
	*/
	store.Users[user.ID] = user

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

/*
	untuk perbaikan performance adalah dengan mengubah metode
	searching User by Username, menjadi searching by key dari
	sebuah variabel berupa Map. Key-nya adalah username, dgn
	asumsi username bersifat unique
*/
func (store FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}

	for _, user := range store.Users {
		if strings.ToLower(username) == strings.ToLower(user.Username) {
			return &user, nil
		}
	}
	return nil, nil
}

/*
	untuk perbaikan performance adalah dengan mengubah metode
	searching User by Email, menjadi searching by key dari
	sebuah variabel berupa Map. Key-nya adalah email, dengan
	asumsi email bersifat unique
*/
func (store FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}

	for _, user := range store.Users {
		if strings.ToLower(email) == strings.ToLower(user.Email) {
			return &user, nil
		}
	}
	return nil, nil
}
