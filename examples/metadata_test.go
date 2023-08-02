package examples

import (
	"testing"

	"github.com/testit-tms/adapters-go/pkg/tms"
)

func TestMetadata_without_metadata_success(t *testing.T) {
	tms.Test(t, tms.TestMetadata{},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_without_metadata_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_external_id_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with external id success",
			ExternalId:  "with_external_id_success",
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_external_id_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with external id failed",
			ExternalId:  "with_external_id_failed",
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_title_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with title success",
			Title:       "with title",
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_title_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with title failed",
			Title:       "with title",
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_work_item_ids_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with work item ids success",
			WorkItemIds: []string{"12345", "54321"},
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_work_item_ids_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with work item ids failed",
			WorkItemIds: []string{"12345", "54321"},
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_description_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with description success",
			Description: "with description",
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_description_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with description failed",
			Description: "with description",
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_labels_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with labels success",
			Labels:      []string{"label01", "label02"},
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_labels_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with labels failed",
			Labels:      []string{"label01", "label02"},
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_class_name_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with class name success",
			ClassName:   "CustomClassName",
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_class_name_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with class name failed",
			ClassName:   "CustomClassName",
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_name_space_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with name space success",
			NameSpace:   "CustomNameSpace",
		},
		func() {
			tms.True(t, true)
		})

}

func TestMetadata_with_name_space_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with name space failed",
			NameSpace:   "CustomNameSpace",
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_links_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with links success",
			Links: []tms.Link{
				{
					Url: "http://google.com",
				},
				{
					Url:         "http://google.com",
					Title:       "Google",
					Description: "Google search engine",
					LinkType:    tms.LINKTYPE_RELATED,
				},
			},
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_links_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with links failed",
			Links: []tms.Link{
				{
					Url: "http://google.com",
				},
				{
					Url:         "http://google.com",
					Title:       "Google",
					Description: "Google search engine",
					LinkType:    tms.LINKTYPE_RELATED,
				},
			},
		},
		func() {
			tms.True(t, false)
		})
}

func TestMetadata_with_all_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with all success",
			ExternalId:  "with_all_success",
			Title:       "with title",
			WorkItemIds: []string{"12345", "54321"},
			Description: "with description",
			Labels:      []string{"label01", "label02"},
			ClassName:   "CustomClassName",
			NameSpace:   "CustomNameSpace",
			Links: []tms.Link{
				{
					Url: "http://google.com",
				},
				{
					Url:         "http://google.com",
					Title:       "Google",
					Description: "Google search engine",
					LinkType:    tms.LINKTYPE_RELATED,
				},
			},
		},
		func() {
			tms.True(t, true)
		})
}

func TestMetadata_with_all_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "with all failed",
			ExternalId:  "with_all_failed",
			Title:       "with title",
			WorkItemIds: []string{"12345", "54321"},
			Description: "with description",
			Labels:      []string{"label01", "label02"},
			ClassName:   "CustomClassName",
			NameSpace:   "CustomNameSpace",
			Links: []tms.Link{
				{
					Url: "http://google.com",
				},
				{
					Url:         "http://google.com",
					Title:       "Google",
					Description: "Google search engine",
					LinkType:    tms.LINKTYPE_RELATED,
				},
			},
		},
		func() {
			tms.True(t, false)
		})
}
