# goforge

![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)

`goforge` is a CLI tool built with Go to accelerate the creation of **Modular Monolith** projects. This tool automatically generates a directory structure, boilerplate code, and commonly used dependencies, allowing you to focus immediately on business logic.

## Key Features

-   **Rapid Project Initialization**: Create the foundation of a new project with a single command.
-   **Modular Structure**: Generates a directory structure designed for scalability, separating domains, APIs, and configurations.
-   **Web Server Setup**: Comes pre-configured with a web server using Gin, complete with a "Hello World" example endpoint.
-   **Dependency Injection**: Integrated with Google Wire for clean dependency injection. The `go generate` process is run automatically.
-   **Automatic Dependency Installation**: All essential dependencies (Gin, Viper, GORM, Wire) are automatically downloaded during initialization.
-   **Database Selection**: Provides an option to choose between PostgreSQL and MySQL during the project creation process.

## Installation

1.  **Clone the Repository**
    ```sh
    git clone [https://github.com/Ryftri/goforge.git](https://github.com/Ryftri/goforge.git)
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

To create a new project, run the following command from the `goforge` directory:

```sh
# For Linux/macOS
./goforge init your-project-name

# For Windows
goforge.exe init your-project-name
```

The tool will prompt you to select a database and will automatically create a new directory named `your-project-name`, install all dependencies, and run `go generate`.

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
├── config.yaml
├── go.mod
└── go.sum
```

## Next Steps

After the project is successfully created, everything is ready to go! You do not need to run `go generate` or `go get` again.

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

<script type="text/javascript" src="https://cdnjs.buymeacoffee.com/1.0.0/button.prod.min.js" data-name="bmc-button" data-slug="ryftri" data-color="#171b3b" data-emoji="☕"  data-font="Poppins" data-text="Buy me a coffee" data-outline-color="#ffffff" data-font-color="#ffffff" data-coffee-color="#FFDD00" ></script>