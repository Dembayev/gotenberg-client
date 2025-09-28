package gotenberg

import (
	"bytes"
	"context"
	"net/http"
	"testing"
)

func BenchmarkRequestChain(b *testing.B) {
	httpClient := &http.Client{}
	c, err := NewClient(httpClient, "http://example.com")
	if err != nil {
		b.Fatalf("new client: %v", err)
	}
	ctx := context.Background()
	html := []byte("<html><body>hello</body></html>")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := c.ConvertHTML(ctx, bytes.NewReader(html)).
			File(FieldFiles, FileIndexHTML, bytes.NewReader(html)).
			WebhookURLMethodPost("http://example.com/webhook").
			WebhookExtraHeaders(map[string]string{"k": "v"}).
			OutputFilename("doc.pdf").
			Bool(FieldPrintBackground, true).
			Float(FieldScale, 1.23).
			PaperSizeA4().
			Margins(0.1, 0.1, 0.1, 0.1).
			Header("X-Test", "1").
			FormField("foo", "bar")
		// touch the internal error so the compiler doesn't optimize the chain away
		_ = r.err
	}
}
