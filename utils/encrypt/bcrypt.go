package encrypt

import "golang.org/x/crypto/bcrypt"

func Brypt(payload string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(payload), bcrypt.DefaultCost)
	var hashErr error
	if err != nil {
		hashErr = err
		panic(err)
	}
	return string(hash), hashErr
}

func CompareHashBrypt(payload string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(payload))
	return err == nil
}
