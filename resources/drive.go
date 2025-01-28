package resources

type Drive struct {
	Id          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	DriveType   string `json:"driveType,omitempty"`
	Name        string `json:"name,omitempty"`
	Quota       Quota  `json:"quota,omitempty"`
}

type Quota struct {
	Total     int64  `json:"total,omitempty"`
	Used      int64  `json:"used,omitempty"`
	Remaining int64  `json:"remaining,omitempty"`
	Deleted   int64  `json:"deleted,omitempty"`
	State     string `json:"state,omitempty"`
	FileCount int64  `json:"fileCount,omitempty"`
}
