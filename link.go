package tms

type Link struct {
	Url         string
	Title       string
	Description string
	LinkType    LinkType
}

type LinkType string

const (
	LINKTYPE_RELATED     LinkType = "Related"
	LINKTYPE_BLOCKED_BY  LinkType = "BlockedBy"
	LINKTYPE_DEFECT      LinkType = "Defect"
	LINKTYPE_ISSUE       LinkType = "Issue"
	LINKTYPE_REQUIREMENT LinkType = "Requirement"
	LINKTYPE_REPOSITORY  LinkType = "Repository"
)
