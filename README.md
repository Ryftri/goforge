# goforge

[![Buy Me a Coffee](https://img.buymeacoffee.com/button-api/?text=Buy%20me%20a%20coffee&emoji=%E2%98%95&slug=ryftri&button_colour=271f47&font_colour=ffffff&font_family=Poppins&outline_colour=ffffff&coffee_colour=FFDD00)](https://www.buymeacoffee.com/ryftri)

![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)

`goforge` is a CLI tool built with Go to accelerate the creation of **Modular Monolith** projects. This tool automatically generates a directory structure, boilerplate code, and commonly used dependencies, allowing you to focus immediately on business logic.

## Key Features

-   **Rapid Project Initialization**: Create the foundation of a new project with a single command.
-   **Modular Structure**: Generates a directory structure designed for scalability, separating domains, APIs, and configurations.
-   **Web Server Setup**: Comes pre-configured with a web server using Gin, complete with a "Hello World" example endpoint.
-   **Opinionated Tooling**: Comes with Google Wire for dependency injection and a pre-configured setup for Mockery to encourage good testing practices.
-   **Automatic Dependency Installation**: All essential dependencies (Gin, Viper, GORM, Wire, Mockery) are automatically downloaded during initialization.
-   **Custom Module Name**: Prompts for a custom Go module name during setup.
-   **Database Selection**: Provides an option to choose between PostgreSQL and MySQL during the project creation process.

## Prerequisites

Before using `goforge`, you need to have Go installed on your system (version 1.21 or higher).

Additionally, this tool generates a project that relies on <a href="https://github.com/google/wire" target="_blank">Google Wire</a>. While `goforge` installs the Wire CLI for you inside the generated project, it's a good practice to have it installed globally:
```sh
go install github.com/google/wire/cmd/wire@latest
```

## Installation

### For Users (Recommended)

You can install the `goforge` tool with a single command:
```sh
go install github.com/Ryftri/goforge
```
This will compile and place the `goforge` binary in your Go bin directory, allowing you to run it from any terminal.

### For Developers (From Source)

If you wish to contribute to `goforge`, you can build it from source.
1.  **Clone the Repository**
    ```sh
    git clone https://github.com/Ryftri/goforge.git
    ```
2.  **Navigate to the Directory**
    ```sh
    cd goforge
    ```
3.  **Build the Binary**
    ```sh
    go build
    ```

## Usage

Run the following command to create a new project:
```sh
goforge init your-project-name
```
The tool will prompt you for a Go module name and your preferred database. It will then automatically create a new directory named `your-project-name`, install all dependencies, and run `go generate` to create the dependency injection code.

### A Note on Included Tools

-   **Google Wire (Required)**: The generated project is fundamentally built around Google Wire for dependency injection. It is required for the application to compile and run.
-   **Mockery (Optional)**: We've included a default `.mockery.yaml` configuration and installed the <a href="https://vektra.github.io/mockery/v3.0/installation/" target="_blank">Mockery tool v3</a> to make setting up test mocks easier. If you do not wish to use Mockery, you can simply delete the `.mockery.yaml` file. This will not affect the running application.

## Generated Project Structure

`goforge` will generate a project structure like the following:
```
your-project-name/
├── api/
│   └── v1/
│       ├── handler/
│       │   └── hello.go
│       └── router.go
├── cmd/
│   └── api/
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── internal/
│   └── category/
├── migrations/
├── pkg/
│   ├── config/
│   │   └── config.go
│   └── database/
├── .mockery.yaml
├── config.yaml
├── go.mod
└── go.sum
```

## Next Steps

After the project is successfully created, everything is ready to go!

1.  **Navigate to Your Project Directory**
    ```sh
    cd your-project-name
    ```
2.  **Run the Server**
    ```sh
    go run cmd/api/main.go
    ```
3.  **Access Your Endpoint**
    The server will be running at `localhost:8080`. You can test the default endpoint with `curl` or a web browser:
    ```sh
    curl http://localhost:8080/api/v1/hello-world
    ```
    
[![Buy Me a Coffee](https://img.buymeacoffee.com/button-api/?text=Buy%20me%20a%20coffee&emoji=%E2%98%95&slug=ryftri&button_colour=271f47&font_colour=ffffff&font_family=Poppins&outline_colour=ffffff&coffee_colour=FFDD00)](https://www.buymeacoffee.com/ryftri)
