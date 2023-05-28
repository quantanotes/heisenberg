# Heisenberg

Heisenberg is a simple vector database optimised for LLMs.

Why do we need vector databases? Simple, it is the language of AI. 

Vectors represent complex data such as images, sounds or text documents as a list of numbers - vectors. By comparing distances between vectors we can find instances of data which are similar.

## Use cases

- Image search. Say you have a database of vectorised images. You can take a text prompt, turn it in to a vector then using Heisenberg find the images which have the closest match to the prompt.
- Semantic search. Conventional text search finds phrases or words which closely match the prompt. Heisenberg however is better suited to finding pieces of texts which are conceptually related i.e. "movies" -> "The Matrix", "Guardians of the Galaxy" etc.
- LLM memory. LLMs themselves can be pretty stupid. Just like humans, they need memory. By combining semantic search, Heisenberg can facilitate more complex LLM usage (built in LLM modules are coming soon).

Currently Heisenberg is still under development so it is not recommended to use for production.

## Example

Normally to convert words in to vectors you utilise an embedding model. You can utilise the OpenAI api for example to obtain vectors for a word. Here we will use a simpler model.

Consider a two dimensional vector [x, y] where x represents gender i.e. 0 is male and 1 is female and y represents power where 1 is the most powerful.

We can represent the following words as such:

```go
func getWords() map[string][]float32 {
    words := make(map[string][]float32)
    m["man"] := [0, 0.5]
    m["woman"] := [1, 0.5]
    m["person"] := [0.5, 0.5]
    m["king"] := [1, 0]
    m["queen"] := [0, 1]
    m["servant"] := [0.5, 0]
    m["maid"] := [1, 0]
} 
```

This gives us a semantic map of words which are closely related. We can then search through them as so:

```go
import (
    "fmt"
    "heisenberg/core"
)


func main() {
    words := getWords()
    h := NewHeisenberg("./path/to/heisenberg/folder")
    // Collections are an isolated search group of vectors.
    // To create a new collection specify the name, and the size of vectors to be stored
    err := h.NewCollection("people", 2)
    if err != nil {
        panic(err)
    }

}

```