package documentloaders

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

// NotionDirectoryLoader is a document loader that reads content from pages within a Notion Database.
type NotionDirectoryLoader struct {
	filePath string
	encoding string
}

// NewNotionDirectory creates a new NotionDirectoryLoader with the given file path and encoding.
func NewNotionDirectory(filePath string, encoding ...string) *NotionDirectoryLoader {
	defaultEncoding := "utf-8"

	if len(encoding) > 0 {
		return &NotionDirectoryLoader{
			filePath: filePath,
			encoding: encoding[0],
		}
	}

	return &NotionDirectoryLoader{
		filePath: filePath,
		encoding: defaultEncoding,
	}
}

// Load retrieves data from a Notion directory and returns a list of schema.Document objects.
func (n *NotionDirectoryLoader) Load(ctx context.Context) ([]schema.Document, error) {
	return n.load(ctx, n.filePath)
}

func (n *NotionDirectoryLoader) load(ctx context.Context, path string) ([]schema.Document, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	documents := make([]schema.Document, 0, len(files))
	for _, file := range files {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() { // || filepath.Ext(file.Name()) != ".md" {
			doc, err := n.load(ctx, filePath)
			if err != nil {
				slog.Info("cannot load sub dir", "path", filePath, "err", err)
			}
			documents = append(documents, doc...)
			continue
		}

		if strings.HasPrefix(file.Name(), "Modelfile") {
			continue
		}
		slog.Info("Processing file", "file", filePath)
		text, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		metadata := map[string]interface{}{"source": filePath}
		documents = append(documents, schema.Document{PageContent: string(text), Metadata: metadata})
	}

	return documents, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple
// documents using a text splitter.
func (n NotionDirectoryLoader) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := n.Load(ctx)
	if err != nil {
		return nil, err
	}
	return textsplitter.SplitDocuments(splitter, docs)
}
