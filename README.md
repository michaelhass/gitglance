# gitglance

Simple git terminal UI written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea).
![gitglance demo](assets/gitglance_demo.gif)

## Features
- Stage files ✔️
- Unstage files ✔️
- Reset files ✔️
- View diffs ✔️
- Commit ✔️
- Refresh Status ✔️
- Open Editor ✔️

## Installation

### Go
```
go install github.com/michaelhass/gitglance@latest
```

### Homebrew
Tap:

```
brew install michaelhass/gitglance/gitglance
```

### Build & run locally
Build & run:
```
make build
make run
```

Debugging:

You can run gitglance in debug mode, which will write logs to a **debug.log** file.
```
// Run the application in debug mode
make debug
// Attach to the debug log in another terminal session
make observe_log
```

More details in the [Makefile](Makefile).
## Configuration

### Editor
Gitglance can try to open an editor for selected files.
If not configured, **vi** will be used. The editor is chosen from a list of options with the following priority.
1. **git config core.editor**
2. **git config --global core.editor**
3. env **VISUAl**
4. env **EDITOR**

While the editor is open, gitglance will be paused. Once the editor process finishes, gitglance resumes and updates the current status.
For external editors, it may not be possible to correctly determine, if the editor process has finsihed, unless correctly configured.
Example
```
// open zed editor in a new window (-n) and wait (-w)
export VISUAL="zed -w -n"
```

## Inspiration
- [lazygit](https://github.com/jesseduffield/lazygit)
- [GitUI](https://github.com/extrawurst/gitui)
