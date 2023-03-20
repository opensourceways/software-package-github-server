package messageserver

type Config struct {
	Group          string         `json:"group"    required:"true"`
	Topics         Topics         `json:"topics"`
	TopicsToNotify TopicsToNotify `json:"topics_to_notify"`
}

type Topics struct {
	ApprovedPkg string `json:"approved_pkg" required:"true"`
	MergedPR    string `json:"merged_pr"    required:"true"`
}

type TopicsToNotify struct {
	CreatedRepo string `json:"created_repo" required:"true"`
}
