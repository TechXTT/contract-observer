# Contract Observer

> A tool to observe Lock events on a contract and vote on them.

## Installing / Getting started

```shell
go run cmd/main.go
```

If `.env` file is not present, the program will not run. The `.env` file should contain information from the `.env.example` file.
With the `.env` file in place, the program will run and print the chain ID it's listening to.

### Initial Configuration

```shell
go mod download
```

## Developing

Here's a brief intro about what a developer must do in order to start developing
the project further:

```shell
git clone https://github.com/TechXTT/contract-observer.git
cd contract-observer/
go mod download
```

<!--
### Deploying / Publishing

In case there's some step you have to take that publishes this project to a
server, this is the right time to state it.

```shell
packagemanager deploy awesome-project -s server.com -u username -p password
```

And again you'd need to tell what the previous code actually does. -->

## Features

What's all the bells and whistles this project can perform?

- Listens to Lock events on a contract
- Votes on Lock events

## Links

Even though this information can be found inside the project on machine-readable
format like in a .json file, it's good to include a summary of most useful
links to humans using your project. You can include links like:

- Repository: https://github.com/TechXTT/contract-observer
- Related projects:
  - Bridge contract: https://github.com/SamBorisov/Bridge/tree/Bridge2
