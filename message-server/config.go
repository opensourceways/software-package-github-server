package messageserver

import "github.com/opensourceways/software-package-github-server/softwarepkg/infrastructure/messageimpl"

type Config struct {
	Group   string             `json:"group"    required:"true"`
	Topics  Topics             `json:"topics"`
	Message messageimpl.Config `json:"message"`
}

type Topics struct {
	ApprovedPkg string `json:"approved_pkg" required:"true"`
	MergedPR    string `json:"merged_pr"    required:"true"`
}
