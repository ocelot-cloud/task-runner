### Task-Runner

Handling CI-related tasks such as building, testing, and deploying can be automated using bash scripts, but these can quickly get messy. Managing multiple processes, handling cleanup, aborting on failure, adding pretty colored output, and more can lead to complex bash code. That's why Task-Runner was created: a simple, open source tool to replace bash scripts with cleaner, more maintainable code written in Go.

### Installation

```bash
go get github.com/ocelot-cloud/task-runner
```

### Usage Example

```go
var backendDir = "../backend"
var frontendDir = "../frontend"
var acceptanceTestsDir = "../acceptance"

func TestFrontend() {
    tr.PrintTaskDescription("Testing Integrated Components")
    defer tr.Cleanup() // shuts down the daemon processes at the end
    tr.ExecuteInDir(backendDir, "go build")
    tr.StartDaemon(backendDir, "./backend")
    tr.WaitUntilPortIsReady("8080")

    tr.ExecuteInDir(frontendDir, "npm install")
    tr.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE=TEST")
    tr.WaitForWebPageToBeReady("http://localhost:8081/")
    tr.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE=TEST")
}
```

The idea is to write simple functions like this and build a CLI tool, e.g., by using [cobra](https://github.com/spf13/cobra), to call these functions. The final use of the CLI tool might look like this: 

```bash
go build
./my-task-runner test frontend
```

This approach helps you create a modern and scalable CI infrastructure. Here is an example of how to initialize the tool in the main function: 

```go
func main() {
    // Optional. If you enable this, when a process hangs, you can press "CTRL" + "C" which will 
    // call the cleanup function and try to gracefully shut down the process. If that does not work, 
    // it will forcefully exit the program.
    go tr.HandleSignals()

    // Optional. Environment variables are applied to each command called.
    tr.DefaultEnvs = []string{"LOG_LEVEL=DEBUG"}

    // Optional. Is called as sub-function whenever tr.Cleanup() is called. 
    // Can be used to add custom post-task cleanup functionality.
    tr.CustomCleanupFunc = MyCustomCleanupFunction

    // Sample usage of cobra library
    if err := rootCmd.Execute(); err != nil {
        tr.ColoredPrintln("\nError during execution: %v\n", err)
        tr.CleanupAndExitWithError()
    }
}

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "some brief description",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("hello world")
    },
}
```

### Compatibility

The tool has been tested under Linux and Windows. In order to be independent of the Linux tools `cp`, `rm`, `mkdir` and `move`, which are quite useful for handling files and folders on the system, operating system independent functions have been implemented:

```
tr.Copy(...)
tr.Remove(...)
tr.MakeDir(...)
tr.Move(...)
```

### License

This project is open source licensed under the [Zero-Clause BSD](./LICENSE) license.
