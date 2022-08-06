package pterodactyl

import "time"

type StartupResp struct {
	Object string `json:"object"`
	Data   []struct {
		Object     string `json:"object"`
		Attributes struct {
			Env_Variable string `json:"env_variable"`
			Srv_Value    string `json:"server_value"`
		} `json:"attributes"`
	} `json:"data"`
}

type ListFilesResp struct {
	Object string `json:"object"`
	Data   []struct {
		Object     string `json:"object"`
		Attributes struct {
			Name       string    `json:"name"`
			Mode       string    `json:"mode"`
			ModeBits   string    `json:"mode_bits"`
			Size       int       `json:"size"`
			IsFile     bool      `json:"is_file"`
			IsSymlink  bool      `json:"is_symlink"`
			Mimetype   string    `json:"mimetype"`
			CreatedAt  time.Time `json:"created_at"`
			ModifiedAt time.Time `json:"modified_at"`
		} `json:"attributes"`
	} `json:"data"`
}

type SendCmdReq struct {
	Command string `json:"command"`
}

type SendPowerCmdReq struct {
	Signal string `json:"signal"`
}

type DeleteFileReq struct {
	Root  string   `json:"root"`
	Files []string `json:"files"`
}

type PteroResp struct {
	Errors []struct {
		Code   string `json:"code"`
		Status string `json:"status"`
		Detail string `json:"detail"`
	} `json:"errors"`
}
