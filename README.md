[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/philipf/gt)
![Build workflow](https://github.com/philipf/gt/actions/workflows/go.yml/badge.svg)

# gt
GT (Go Time) is a CLI to manage time.

This is used as a learning module for Go


## Install
(There is nothing to see yet)

```bash
go install github.com/philipf/gt/cmd/gt@latest
```
## Next steps
- Revisit the code structure for the domain layer. This doesn't feel like idiomatic Go yet
    - breakdown?:
        - calendar (bounded-context) / module  (ask chat-gpt what is difference between domain, bounded-context and module). have to be careful as module might clash with golang concepts
            - model
            - app_services
            - domain_services
            - repository
        - infra??

## Done  
- Update existing tests to consistently use the asserts package
- Fix error handling, using Error types
