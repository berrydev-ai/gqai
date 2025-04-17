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

- 🧰 Define tools using GraphQL operations
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
schema: https://graphql.org/graphql/
documents: .
```

This file tells gqai where to find your GraphQL schema and operations.

*Note: The schema also tells gqai where to execute the operations. This must be a live server rather than a static schema file*

2. Write a GraphQL operation

`get_all_films.graphql`:
```graphql
# Get all Star Wars films
query get_all_films {
  allFilms {
    films {
      title
      episodeID
    }
  }
}
```

3. Add gqai to your `mcp.json` file:

```
  "gqai": {
    "command": "gqai",
    "args": [
      "run",
      "--config"
      ".graphqlrc.yml"
    ]
  }
```

## 🧪 CLI Testing
### Call a tool via CLI to test:

```bash
gqai tools/call get_all_films
```

This will execute the `get_all_films` tool and print the result.

```shell
{
  "data": {
    "allFilms": {
      "films": [
        {
          "id": 4,
          "title": "A New Hope"
        },
        {
          "id": 5,
          "title": "The Empire Strikes Back"
        },
        {
          "id": 6,
          "title": "Return of the Jedi"
        },
        ...
      ]
    }
  }
}
```
### Call a tool with arguments:

Create a GraphQL operation that takes arguments, and these will be the tool inputs:

`get_film_by_id.graphql`:
```graphql
query get_film_by_id($id: ID!) {
  film(filmID: $id) {
    episodeID
    title
    director
    releaseDate
  }
}
```

Call the tool with arguments:

```bash
gqai tools/call get_film_by_id '{"id": "1"}'
```

This will execute the `get_film_by_id` tool with the provided arguments.

```shell
{
  "data": {
    "film": {
      "episodeID": 1,
      "title": "A New Hope",
      "director": "George Lucas",
      "releaseDate": "1977-05-25"
    }
  }
}
```

## 📦 Tool Metadata
Auto-generated tool specs for each operation, so you can plug into any LLM that supports tool use.

## 🤖 Why gqai?
gqai makes it easy to turn your GraphQL backend into a model-ready tool layer — no code, no extra infra. Just define your operations and let AI call them.

## 📜 License
MIT — fork it, build on it, all the things.

## 👋 Author
Made with ❤️ and 🤖vibes by Stephen Spalding && <your-name-here>
