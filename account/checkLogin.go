package account

import (
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/lemon-mint/vbox"
)

func DecodeTocken(tocken string) (userId []byte, _err bool) {
	//get key form config folder
	file, _ := os.Open("config/accessKey.txt")
	key, _ := ioutil.ReadAll(file)

	//decode userInfo
	auth := vbox.NewBlackBox(key)
	decodedTocken, _ := hex.DecodeString(tocken)
	userId, boolean := auth.Open(decodedTocken)

	if !boolean {
		return nil, true
	}

	return userId, false
}
