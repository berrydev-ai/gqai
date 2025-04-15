# gqai
*graphql → ai*

**gqai** is a lightweight proxy that exposes GraphQL operations as [Model Context Protocol (MCP)](https://platform.openai.com/docs/guides/function-calling) tools for AI models like ChatGPT.  
It lets you define tools using regular GraphQL queries/mutations and run them locally or expose them over HTTP.

🔌 Powered by your GraphQL backend  
⚙️ Driven by `.graphqlrc.yml` + plain `.graphql` files  
🧠 Model-ready inputs/outputs
📍 Built in Go

---

## ✨ Features

- 🧰 Run GraphQL operations as tools via CLI
- 🌐 Serve tools via HTTP for AI agents
- 🗂 Automatically discover operations from `.graphqlrc.yml`
- 🧾 Tool metadata compatible with OpenAI function calling / MCP

---

## 🛠️ Installation

```bash
go install github.com/fotoetienne/gqai@latest
```


## 🚀 Quick Start
1. Create a .graphqlrc.yml:

```yaml
schema: "https://graphql.org/graphql/"
documents: "operations"
```

This file tells gqai where to find your GraphQL schema and operations.

*Note: The schema also tells gqai where to execute the operations. This must be a live server rather than a static schema file*

2. Write a GraphQL operation (operations/getAllFilms.graphql):

```graphql
query getAllFilms {
  allFilms {
    films {
      title
      episodeID
    }
  }
}
```

3. Run it via CLI:

```bash
gqai run getAllFilms
```

## 🌐 HTTP Server
Run a local server that exposes your tools via HTTP:

```bash
gqai serve
```

Call the tools via HTTP POST requests:

```bash
curl -X POST http://localhost:8080/tools/call  -d '{"toolName": "getAllFilms"}'
```

## 📦 Tool Metadata
Auto-generated tool specs for each operation, so you can plug into any LLM that supports tool use.


## 🤖 Why gqai?
gqai makes it easy to turn your GraphQL backend into a model-ready tool layer — no wrappers, no extra infra. Just define your operations and let AI call them.

## 🧪 Roadmap
  - [x] CLI tool runner

  - [x] HTTP server mode (gqai serve)

  - [x] Tool metadata generation

  - [ ] JSON Schema validation

  - [ ] Plug-and-play with OpenAI, Fireworks, etc.

## 📜 License
MIT — fork it, build on it, model all the things.

## 👋 Author
Made with ❤️ and 🤖vibes by Stephen Spalding 
