package pterodactyl

type StartupResp struct {
	Object string `json:"object"`
	Data []struct {
		Object string `json:"object"`
		Attributes struct {
			EnvVariable string `json:"env_variable"`
			SrvValue string `json:"server_value"`
		} `json:"attributes"`
	} `json:"data"`
}

type ListFilesResp struct {
	Object string `json:"object"`
	Data []struct {
		Object string `json:"object"`
		Attributes struct {
			Name string `json:"name"`
			IsFile bool `json:"isfile"`
		} `json:"attributes"`
	} `json:"data"`
}

type SendCmdReq struct {
	Command string `json:"command"`
}

type UpdateVarReq struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

type SendPowerCmdReq struct {
	Signal string `json:"signal"`
}

type DeleteFileReq struct {
	Root string `json:"root"`
	Files []string `json:"files"`
}