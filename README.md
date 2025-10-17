# ReLive - Go project process reloader

## Table of Contents

- [Motivation](#motivation)
- [Information](#information)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Motivation

While working on my API, I noticed that I change the code all the time, so I have to restart my server constantly. I knew there are ways to fix this easily with another Go projects, like [Air](https://github.com/air-verse/air), but I came up with an idea and was able to make it my own.

It might not be the best tool possible for this purpose, but it's working, and and it fulfills its purpose. If you know, you could make it better, feel free to contribute.

## Information

> [!NOTE]
> **Requirements:**
> - `go>=1.25.1`

## Installation

```bash
curl -sS https://webi.sh/golang | sh # Install Go
go install github.com/MedrekIT/relive@latest # Install repository as program for global execution
```

## Usage

### Run
```bash
relive <project_path>
# OR
relive # in project's directory
```

## Contributing

### Install project

```bash
git clone https://github.com/MedrekIT/relive.git # Clone repository
cd relive # Move to repository directory
go mod tidy
```

### Add you changes and submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## License

This project is licensed under the [MIT License](LICENSE).

You are free to use, modify, and distribute this software under the terms of that license.
