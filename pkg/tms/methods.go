package tms

import (
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

func AddMessage(m string) {
	if node, ok := ctxMgr.GetValue(testResultKey); ok {
		tr := node.(*testResult)
		tr.message = m
	}
}

func AddLinks(l Link) {
	if node, ok := ctxMgr.GetValue(testResultKey); ok {
		tr := node.(*testResult)
		tr.resultLinks = append(tr.resultLinks, l)
	}
}

func AddAtachments(paths ...string) {
	atachs := client.writeAttachments(paths...)

	if node, ok := ctxMgr.GetValue(nodeKey); ok {
		n := node.(hasAttachments)
		for _, a := range atachs {
			n.addAttachments(a)
		}
	}
}

func AddAtachmentsFromString(name, content string) {
	path, err := os.Getwd()
	if err != nil {
		logger.Error("cannot get executable path: %s", err, slog.With("op", "AddAtachmentsFromString"))
	}

	fp := filepath.Join(path, name)
	err = os.WriteFile(fp, []byte(content), 0644)
	if err != nil {
		logger.Error("cannot write file: %s", err, slog.With("op", "AddAtachmentsFromString"))
	}

	atachs := client.writeAttachments(fp)

	if node, ok := ctxMgr.GetValue(nodeKey); ok {
		n := node.(hasAttachments)
		for _, a := range atachs {
			n.addAttachments(a)
		}
	}

	err = os.Remove(fp)
	if err != nil {
		logger.Error("cannot remove file: %s", err, slog.With("op", "AddAtachmentsFromString"))
	}
}
