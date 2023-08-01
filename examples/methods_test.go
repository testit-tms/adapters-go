package examples

import (
	"os"
	"path/filepath"
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

func TestMethods_attachments_success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "add attachments success",
		},
		func() {
			path, err := os.Getwd()
			if err != nil {
				t.Errorf("cannot get executable path: %s", err)
			}

			tms.AddAtachments(
				filepath.Join(path, "attachments", "file01.txt"),
				filepath.Join(path, "attachments", "file02.txt"),
			)

			tms.Step(
				tms.StepMetadata{
					Name: "add attachments to step",
				}, func() {
					tms.AddAtachments(
						filepath.Join(path, "attachments", "file03.txt"),
					)

					tms.AddAtachmentsFromString("myFile.txt", "content of my file")
				},
			)

			assert.True(t, true)
		})
}

func TestMethods_attachments_failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "add attachments failed",
		},
		func() {
			path, err := os.Getwd()
			if err != nil {
				t.Errorf("cannot get executable path: %s", err)
			}

			tms.AddAtachments(
				filepath.Join(path, "attachments", "file01.txt"),
				filepath.Join(path, "attachments", "file02.txt"),
			)

			tms.Step(
				tms.StepMetadata{
					Name: "add attachments to step",
				}, func() {
					tms.AddAtachments(
						filepath.Join(path, "attachments", "file03.txt"),
					)

					tms.AddAtachmentsFromString("myFile.txt", "content of my file")
				},
			)
			assert.True(t, false)
		})
}
