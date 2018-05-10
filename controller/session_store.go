package controller

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/*
	tipe struct untuk menyimpan session
*/
type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

// variabel SessionStore public/Singleton
var globalSessionStore SessionStore

/*
	fungsi untuk load daftar session dari json file
	daftar session yang masih aktif akan tersimpan di
	sini untuk beberapa lama, sesuai dengan lamanya
	waktu yang telah ditentukan pada file session.go
*/
// func init() {
// 	sessionStore, err := NewFileSessionStore("./data/sessions.json")
// 	if err != nil {
// 		panic(fmt.Errorf("Error creating session store: %s", err))
// 	}
// 	globalSessionStore = sessionStore
// }

func SetglobalSessionStore(sessionStore SessionStore) SessionStore {
	globalSessionStore = sessionStore

	return globalSessionStore
}

/*
	tipe struct untuk menulis dan membaca session
	dari dan ke file json yang menyimpan list session.
*/
type FileSessionStore struct {
	filename string
	Sessions map[string]Session
}

/*
	inisialisasi session yang berasal dari File JSON
	yang menyimpan session yang masih aktif
*/
func NewFileSessionStore(name string) (*FileSessionStore, error) {

	/*
		pointer FileSessionStore untuk menampung list session
		yang masih aktif.
	*/
	store := &FileSessionStore{
		Sessions: map[string]Session{},
		filename: name,
	}

	// menbaca konten dari file session.json
	contents, err := ioutil.ReadFile(name)

	if err != nil {
		// If it's a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}

	/*
		memampping konten yang merupakan struktur json
		menjadi objek FileSessionStore
	*/
	err = json.Unmarshal(contents, store)

	if err != nil {
		return nil, err
	}

	return store, err
}

func (s *FileSessionStore) Find(id string) (*Session, error) {
	session, exists := s.Sessions[id]
	if !exists {
		return nil, nil
	}

	return &session, nil
}

func (store *FileSessionStore) Save(session *Session) error {
	store.Sessions[session.ID] = *session
	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)
	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}
