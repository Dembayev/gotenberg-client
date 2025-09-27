package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/nativebpm/gotenberg-client"
)

func main() {
	gotenbergURL := `http://localhost:3000`

	slog.Info("Example 1: Complete Builder Pattern with HTML conversion")

	htmlContent := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Advanced Builder Pattern Demo</title>
	<link rel="stylesheet" href="styles.css">
</head>
<body>
	<div class="container">
		<h1>Advanced Gotenberg Client - Builder Pattern</h1>
		<p>This demonstrates the complete builder pattern implementation!</p>
		<div class="feature-box">
			<h2>Key Features:</h2>
			<ul>
				<li>Fluent interface for client creation</li>
				<li>Method chaining for configuration</li>
				<li>Type-safe option building</li>
				<li>Readable and maintainable code</li>
			</ul>
		</div>
		<div class="code-example">
			<h3>Code Example:</h3>
			<pre>
client := gotenberg.NewClientBuilder("http://localhost:3000").
    WithTimeout(60 * time.Second).
    Build()

resp, err := client.ConvertHTML().
    WithHTML(htmlContent).
    WithCSS(cssContent).
    PaperSizeA4().
    Margins(1.0, 1.0, 1.0, 1.0).
    PrintBackground(true).
    OutputFilename("advanced-demo.pdf").
    Execute(ctx)
			</pre>
		</div>
	</div>
</body>
</html>`

	cssContent := `
.container {
	max-width: 900px;
	margin: 0 auto;
	font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	line-height: 1.6;
	color: #2c3e50;
	padding: 20px;
}

h1 {
	color: #2980b9;
	border-bottom: 3px solid #3498db;
	padding-bottom: 15px;
	margin-bottom: 30px;
}

h2 {
	color: #27ae60;
	margin-top: 25px;
}

h3 {
	color: #8e44ad;
}

.feature-box {
	background: linear-gradient(135deg, #74b9ff 0%, #0984e3 100%);
	color: white;
	padding: 25px;
	border-radius: 12px;
	margin: 25px 0;
	box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.code-example {
	background: #f8f9fa;
	border-left: 4px solid #007bff;
	padding: 20px;
	margin: 20px 0;
	border-radius: 4px;
}

pre {
	background: #2d3748;
	color: #e2e8f0;
	padding: 15px;
	border-radius: 8px;
	overflow-x: auto;
	font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
	font-size: 12px;
	line-height: 1.4;
}

ul {
	list-style-type: none;
	padding: 0;
}

li {
	background: rgba(255,255,255,0.15);
	margin: 8px 0;
	padding: 10px 15px;
	border-radius: 6px;
	border-left: 4px solid #f39c12;
}

p {
	margin: 20px 0;
	text-align: justify;
	font-size: 16px;
}`

	resp, err := gotenberg.NewClientBuilder(http.DefaultClient, gotenbergURL).
		ConvertHTML().
		WithFile("styles.css", strings.NewReader(cssContent)).
		PaperSizeA4().
		Margins(1.2, 1.0, 1.2, 1.0).
		PrintBackground(true).
		Scale(0.95).
		OutputFilename("advanced-builder-demo.pdf").
		Execute(context.Background(), strings.NewReader(htmlContent))

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create("./advanced-builder-demo.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Advanced builder demo PDF generated successfully",
		"pdf_size_bytes", n,
		"content_type", resp.Header.Get("Content-Type"),
		"output_file", "advanced-builder-demo.pdf")
}
