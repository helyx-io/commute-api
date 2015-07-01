package utils

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"log"
    "strings"
    "net/http"
    "encoding/hex"
    "runtime/debug"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Helper Functions
////////////////////////////////////////////////////////////////////////////////////////////////

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("[ERROR] %s: %s", msg, err)
		panic(err)
	}
}

func RecoverFromError(w http.ResponseWriter) {
    if r := recover(); r != nil {
        err, _ := r.(error)

        log.Printf("Err: % - Stack: %v", err.Error(), string(debug.Stack()));
        http.Error(w, err.Error(), 500)
        return
    }
}


func Int32ToColor(intColor int32) string {
    var b = make([]byte, 3)

    b[0] = uint8(intColor)
    b[1] = uint8(intColor >> 8)
    b[2] = uint8(intColor >> 16)

    return strings.ToUpper(hex.EncodeToString(b))
}
