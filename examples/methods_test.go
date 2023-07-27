package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/adapters-go/pkg/tms"
)

func TestMethods_message_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "add message success",
		},
		func() {
			tms.AddMessage("test message")
			assert.True(t, true)
		})
}

func TestMethods_message_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "add message failed",
		},
		func() {
			tms.AddMessage("test message")
			assert.True(t, false)
		})
}

func TestMethods_link_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "add links success",
		},
		func() {
			tms.AddLinks(tms.Link{
				Url: "https://testit.software",
			})

			tms.AddLinks(tms.Link{
				Url:   "https://testit.software",
				Title: "TestIt",
			})

			tms.AddLinks(tms.Link{
				Url:         "https://testit.software",
				Title:       "TestIt",
				Description: "TestIt is a test management system",
			})

			tms.AddLinks(tms.Link{
				Url:         "https://testit.software",
				Title:       "TestIt",
				Description: "TestIt is a test management system",
				LinkType:    tms.LINKTYPE_RELATED,
			})
			
			assert.True(t, true)
		})
}

func TestMethods_link_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "add links failed",
		},
		func() {
			tms.AddLinks(tms.Link{
				Url: "https://testit.software",
			})

			tms.AddLinks(tms.Link{
				Url:   "https://testit.software",
				Title: "TestIt",
			})

			tms.AddLinks(tms.Link{
				Url:         "https://testit.software",
				Title:       "TestIt",
				Description: "TestIt is a test management system",
			})

			tms.AddLinks(tms.Link{
				Url:         "https://testit.software",
				Title:       "TestIt",
				Description: "TestIt is a test management system",
				LinkType:    tms.LINKTYPE_RELATED,
			})

			assert.True(t, false)
		})
}
