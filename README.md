# URL Shortener

This is a simple URL shortener project written in Go, without using any external frameworks. Instead, it relies on the Go standard library to handle HTTP requests and interactions with a Redis database.

## Prerequisites

Before running the project, make sure you have the following installed:

- Go (version 1.22 or higher)
- Redis (version 6.2.6 or higher)
- Docker

## Running the Project

To run the project, follow these steps:

1. Clone the repository:

    ```bash
    $ git clone https://github.com/gunawanpras/url-shortener.git
    ```
2. Change into the project directory:
    ```bash
    $ cd url-shortener
    ```
3. Check missing dependencies
    ```bash
    $ go mod tidy
    ```
4. Run the redis server
    ```bash
    $ docker-compose up --build --no-start && docker-compose start && docker-compose ps
    ```
5. Build the project:
    ```bash
    $ go build -o main
    ```
6. Run the project:
    ```bash
    $ ./main
    ```
7. Open browser and go to http://localhost:8080   