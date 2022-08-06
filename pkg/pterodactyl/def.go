package pterodactyl

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
			Name    string `json:"name"`
			Is_File bool   `json:"isfile"`
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
