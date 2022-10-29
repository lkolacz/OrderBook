package rest

type Version struct {
	BuildVersion  string `json:"BuildVersion"`
	BuildType     string `json:"BuildType"`
	BuildDateTime string `json:"BuildDateTime"`
}

func DefaultVersion() *Version {
	return &Version{
		BuildType: "dev",
	}
}
