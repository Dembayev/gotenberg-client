package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nativebpm/gotenberg-client"
)

func main() {
	gotenbergURL := `http://localhost:3000`

	httpClient := &http.Client{
		Timeout: 90 * time.Second,
	}

	clientBuilder := gotenberg.NewClientBuilder(httpClient, gotenbergURL)

	// Example HTML content
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Builder Pattern Demo</title>
	<link rel="stylesheet" href="styles.css">
</head>
<body>
	<div class="container">
		<h1>Gotenberg Client - Builder Pattern Demo</h1>
		<p>This PDF was generated using the fluent builder pattern!</p>
		<div class="highlight">
			<h2>Features Demonstrated:</h2>
			<ul>
				<li>Custom CSS styling</li>
				<li>A4 paper size with custom margins</li>
				<li>Background printing enabled</li>
				<li>Scale factor applied</li>
				<li>Custom filename</li>
			</ul>
		</div>
		<p>The builder pattern provides a clean, readable way to configure PDF generation options.</p>
	</div>
</body>
</html>`

	// CSS styling
	cssContent := `
.container {
	max-width: 800px;
	margin: 0 auto;
	font-family: 'Arial', sans-serif;
	line-height: 1.6;
	color: #333;
}

h1 {
	color: #2c3e50;
	border-bottom: 3px solid #3498db;
	padding-bottom: 10px;
}

h2 {
	color: #27ae60;
}

.highlight {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 20px;
	border-radius: 10px;
	margin: 20px 0;
}

ul {
	list-style-type: none;
	padding: 0;
}

li {
	background: rgba(255,255,255,0.1);
	margin: 5px 0;
	padding: 8px 12px;
	border-radius: 5px;
	border-left: 4px solid #f39c12;
}

p {
	margin: 15px 0;
	text-align: justify;
}`

	// Using the fluent builder pattern - much cleaner than functional options!
	resp, err := clientBuilder.ConvertHTML().
		WithHTML(htmlContent).
		WithFile("styles.css", bytes.NewReader([]byte(cssContent))).
		PaperSizeA4().
		Margins(1.2, 1.0, 1.2, 1.0). // top, right, bottom, left in inches
		PrintBackground(true).       // Enable CSS background printing
		Scale(0.9).                  // Slightly reduce scale for better fit
		OutputFilename("builder-pattern-demo.pdf").
		Execute(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create("./builder-pattern-demo.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Builder pattern demo PDF generated successfully",
		"pdf_size_bytes", n,
		"content_type", resp.Header.Get("Content-Type"),
		"trace", resp.Header.Get("Gotenberg-Trace"),
		"output_file", "builder-pattern-demo.pdf")
}
