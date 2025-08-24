# Go Basics Project

This small project walks through Go setup and language basics. It includes:

- A simple CLI in `cmd/hello` demonstrating flags, printing, and using a package.
- A reusable package in `pkg/utils` that shows functions, structs, methods, and interfaces.
- Unit tests in `pkg/utils` demonstrating table-driven tests.

Prerequisites
- Go 1.20+ installed (adjust if you have a different version). Verify with:

```powershell
go version
```

Quick start

```powershell
cd c:\source\Tryouts\golang\go-basics
# Download modules
go mod tidy

# Run the CLI
go run ./cmd/hello --name Alice --repeat 2

# Run tests
go test ./... -v

# Build
go build -o bin/hello ./cmd/hello
./bin/hello --name Bob
```

What to study here
- `cmd/hello/main.go` shows `package main`, `func main()`, flags, and basic control.
- `pkg/utils` shows types, methods, pointer vs value receivers, slices, maps, and error handling.
- `pkg/utils/utils_test.go` shows table-driven tests and use of `testing`.

Next steps
- Try modifying the `Person` struct or add new functions and tests.
- Run `go vet`, `gofmt`, and a linter like `golangci-lint`.

Enjoy learning Go!

Troubleshooting: Windows Defender / antivirus flags
-------------------------------------------------

On Windows you may see an error like:

```
fork/exec C:\Users\HP\AppData\Local\Temp\go-build...\exe\hello.exe: Operation did not complete successfully because the file contains a virus or potentially unwanted software.
```

This is commonly a false positive from real-time antivirus scanning that inspects the temporary executable `go run` builds in the temp folder. To work around this safely:

1. Prefer building to a project-local `bin` directory and run the produced executable instead of `go run` (avoids the temporary exe in the OS temp folder):

```powershell
cd C:\source\Tryouts\golang\go-basics
go build -o .\bin\hello.exe ./cmd/hello
.\bin\hello.exe --name Alice --repeat 2
```

2. If you must use `go run`, you can set a temporary build directory within the project folder for the current session (no admin required):

**PowerShell:**
```powershell
mkdir .\tmp  # create if missing
$env:TEMP = (Resolve-Path .\tmp).Path
$env:TMP = (Resolve-Path .\tmp).Path
go run ./cmd/hello --name Alice --repeat 2
```

**Git Bash:**
```bash
mkdir -p ./tmp
export TEMP=$(pwd)/tmp
export TMP=$(pwd)/tmp
go run ./cmd/hello --name Alice --repeat 2
```

3. If Windows Defender continues to block the build, add an exclusion for the project folder (requires an elevated PowerShell window / admin privileges):

```powershell
# Run PowerShell as Administrator
Add-MpPreference -ExclusionPath 'C:\source\Tryouts\golang\go-basics'
```

4. If you believe the detection is a false positive you can also upload the executable to VirusTotal or consult your AV vendor for a takedown/whitelist. When in doubt, inspect the code (the project here is small and visible) before whitelisting.

Notes
- Avoid permanently disabling real-time protection. Use exclusions only for trusted development folders.
- Building to a local `bin` directory (step 1) is the simplest immediate workaround.