<div align="center">
  <h3 align="center">PROCUT</h3>

  <p align="center">
    A simple and fast hash cracker written in Go.
    <br />
    <a href="https://github.com/CYBER-4RMY/procut/blob/master/README.md"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/CYBER-4RMY/procut/issues">Report Bug</a>
    ·
    <a href="https://github.com/CYBER-4RMY/procut/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

`procut` is a command-line tool for cracking password hashes using a dictionary attack. It is written in Go for performance and concurrency.

### Key Features:

*   **Fast:** Utilizes goroutines to perform concurrent hash cracking.
*   **Flexible:** Supports multiple hashing algorithms (MD5, SHA1, SHA256, SHA512).
*   **User-Friendly:** Simple command-line interface with a progress bar.
*   **Salt Support:** Crack hashes that use a salt.
*   **Batch Processing:** Crack multiple hashes from a file.
*   **Output to File:** Save found passwords to a file.

### Built With

*   [Go](https://go.dev/)

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

*   Go (1.16 or later)

### Installation

1.  Clone the repo
    ```sh
    git clone https://github.com/CYBER-4RMY/procut.git
    ```
2.  Navigate to the `procut` directory
    ```sh
    cd procut
    ```
3.  Install dependencies
    ```sh
    go mod tidy
    ```
4.  Build the project
    ```sh
    go build -o procut (if that procut file does not exist)
    ```

## Usage

After building, you can run the cracker from your terminal.

**Important:** For the Go cracker, all flags (`--alg`, `--dict`, `--salt`, `--output`) **must** come *before* the positional argument `<HASH_OR_FILE>`.

```sh
# Run directly
go run procut.go or procut --alg <alg> --dict <dict> [--salt <salt>] [--output <file>] <HASH_OR_FILE>

# Or build and run the executable
go build -o procut
./procut --alg <alg> --dict <dict> [--salt <salt>] [--output <file>] <HASH_OR_FILE>
```
Where `<HASH_OR_FILE>` can be either a hash string to crack, or a path to a file containing hashes (one per line).

**Tips for creating a hash file:**
When redirecting `sha256sum` (or `md5sum`) output to a file, it often includes a `  -` at the end. Use `awk '{print $1}'` to strip this, ensuring only the hash is in the file.

```bash
# Example: Create hash.txt with only the SHA256 hash of "mrinalini"
echo -n "mrinalini" | sha256sum | awk '{print $1}' > hash.txt
```

**Examples:**

*   **Crack a single hash directly:**
    ```sh
    ./procut --alg sha256 --dict ../passwords.txt --salt somesalt 1f3870be274f6c49b3e31a0c6728957f
    ```
*   **Crack multiple hashes from `hash.txt` and save the output:**
    ```sh
    ./procut --alg sha256 --dict ../passwords.txt --output found.txt hash.txt
    ```

---

<!-- ROADMAP -->
## Roadmap

-   [ ] Add brute-force attack mode
-   [ ] Add mask attack mode
-   [ ] Add automatic hash type detection

See the [open issues](https://github.com/CYBER-4RMY/procut/issues) for a full list of proposed features (and known issues).

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1.  Fork the Project
2.  Create your Feature
3.  Commit your Changes 
4.  Push to the Branch 
5.  Open a Pull Request

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.
