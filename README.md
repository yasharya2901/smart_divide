# Smart Divide

## Description
This is the backend for the `Smart Divide` project.

## Installation
1. Clone the repository
    ```bash
    git clone git@github.com:<username>/smart_divide.git
    ```
2. Install the dependencies
    ```bash
    cd smart_divide
    go mod download
    ```
3. Copy the `.env.example` file to `.env` and update the values
    ```bash
    cp .env.example .env
    ```

4. Load the environment variables
    
    For Linux
    ```bash
    ./load_env.sh
    ```

    For Windows
    ```powershell
    ./load_env.ps1
    ```

5. Run the server
    ```bash
    go run main.go
    ```

