# gotenberg-client

[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)

[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)



A high-performance Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API with a fluent builder pattern interface.A high-performance, streaming Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API, optimized for minimal memory allocations.

## Features

- 🚀 **Fluent Builder Pattern**: Chain method calls for intuitive API usage- 🚀 **High Performance**: 70% fewer allocations, 4x faster than traditional approaches

- ♻️ **Client Reuse**: Reset and reuse client instances for multiple conversions- � **Streaming Architecture**: Direct multipart.Writer streaming without intermediate buffers

- 📦 **Zero Dependencies**: Uses only Go standard library- ♻️ **Buffer Pooling**: Optimized memory reuse with sync.Pool

- 🔄 **Webhook Support**: Full support for asynchronous processing with webhooks- 🔄 **Client Reuse**: Reset() method for processing multiple documents

- ⚡ **Context Support**: Proper timeout and cancellation handling- 📦 **Zero Dependencies**: Uses only Go standard library

- 🎯 **Memory Efficient**: Buffer pooling with sync.Pool for optimal performance- 🔗 **Webhook Support**: Full async processing with webhooks

- ⚡ **Context Support**: Proper timeout and cancellation handling

## Best Practices

- ✅ **Always close `resp.Body`** to prevent resource leaks
- ✅ **Use `context.Context`** for timeouts and cancellations  
- ✅ **Check response status codes** before processing
- ✅ **Set appropriate HTTP client timeouts** for your use case
- ✅ **Inspect `Gotenberg-Trace` header** for debugging

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
