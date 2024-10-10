### Task-Runner

Handling CI-related tasks such as building, testing, and deploying can be automated using bash scripts, but these can quickly get messy. Managing multiple processes, handling cleanup, aborting on failure, adding pretty colored output, and more can lead to complex bash code. That's why Task-Runner was created: a simple, open source tool to replace bash scripts with cleaner, more maintainable CI logic written in Golang. Here's a sample usage:

```go
var backendDir = "../backend"
var frontendDir = "../frontend"
var acceptanceTestsDir = "../acceptance"

func TestFrontend() {
	cli.PrintTaskDescription("Testing Integrated Components")
	defer cli.Cleanup()
	cli.ExecuteInDir(backendDir, "go build")
	cli.StartDaemon(backendDir, "./backend")
	cli.WaitUntilPortIsReady("8080")

	cli.ExecuteInDir(frontendDir, "npm install")
	cli.StartDaemon(frontendDir, "npm run serve", "VITE_APP_PROFILE=TEST")
	cli.WaitForWebPageToBeReady("http://localhost:8081/")
	cli.ExecuteInDir(acceptanceTestsDir, cypressCommand, "CYPRESS_PROFILE=TEST")
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
    // Optional lines:
	go cli.HandleSignals() // If you enable this, when a process hangs, you can press "CTRL" + "C" which will call the cleanup function and try to gracefully shut down the process. If that does not work, it will forcefully exit the program.
    cli.DefaultEnvs = []string{"LOG_LEVEL=DEBUG"} // is applied to each command called
	cli.CustomCleanupFunc = MyCustomCleanupFunction // is called whenever tr.Cleanup is called

	if err := rootCmd.Execute(); err != nil {
		cli.ColoredPrintln("\nError during execution: %v\n", err)
		cli.CleanupAndExitWithError()
	}
}
```

The tool has only been tested on Linux.

### License

This project is licensed under the [Zero-Clause BSD](./LICENSE) license. 

### About

This project is a side-product of [Ocelot-Cloud](https://github.com/ocelot-cloud/ocelot-cloud).