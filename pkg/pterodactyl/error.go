package pterodactyl

import (
	"encoding/json"
)

func IsError(body string) bool {
	var ptero_resp PteroResp

	err := json.Unmarshal([]byte(body), &ptero_resp)

	// Likely not an error if it can't parse into the error structure (though, I guess it could be an error).
	if err != nil {
		return false
	}

	// If we have more than one error, just return true.
	if len(ptero_resp.Errors) > 0 {
		return true
	}

	return false
}
