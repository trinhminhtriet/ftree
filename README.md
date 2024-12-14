# 🚀 FTree

```text
 _____  _____
|  ___||_   _| _ __   ___   ___
| |_     | |  | '__| / _ \ / _ \
|  _|    | |  | |   |  __/|  __/
|_|      |_|  |_|    \___| \___|

```

🚀 FTree: Terminal-based file tree manipulation tool for navigating, viewing, and managing directories and files efficiently.

## ✨ Features

- **File Navigation**: Easily navigate through directories and files using keyboard shortcuts.
- **File Viewing**: View file contents directly within the terminal.
- **File Operations**: Perform basic file operations like create, delete, rename, and move.
- **Search**: Quickly search for files and directories.
- **File System Watcher**: Automatically update the file tree on file system changes.
- **Customizable UI**: Customize the appearance with various themes and styles.
- **Cross-Platform**: Works on macOS, Linux, and Windows.
- **Lightweight**: Minimal dependencies and fast performance.
- **Extensible**: Easily extend functionality with plugins.

## 🚀 Installation

### Binary

Download from [latest releases ](https://github.com/trinhminhtriet/ftree/releases)

### From source

```sh
git clone https://github.com/trinhminhtriet/ftree.git
cd ftree

make install
```

## 💡 Usage

```bash
ftree [flags] [directory]

Flags:
  -i    In-place render (without alternate screen)
  -pad uint
        Edge padding for top and bottom (default 5)
```

Key bindings:

| key           | desc                                                   |
| ------------- | ------------------------------------------------------ |
| j / arr down  | Select next child                                      |
| k / arr up    | Select previous child                                  |
| h / arr left  | Move up a dir                                          |
| l / arr right | Enter selected directory                               |
| d             | Move selected child (then 'p' to paste)                |
| y             | Copy selected child (then 'p' to paste)                |
| D             | Delete selected child                                  |
| if / id       | Create file (if) / directory (id) in current directory |
| r             | Rename selected child                                  |
| e             | Edit selected file in $EDITOR                          |
| gg            | Go to top most child in current directory              |
| G             | Go to last child in current directory                  |
| enter         | Collapse / expand selected directory                   |
| esc           | Clear error message / stop current operation           |
| ?             | Toggle help                                            |
| q / ctrl+c    | Exit                                                   |

## 🤝 How to contribute

We welcome contributions!

- Fork this repository;
- Create a branch with your feature: `git checkout -b my-feature`;
- Commit your changes: `git commit -m "feat: my new feature"`;
- Push to your branch: `git push origin my-feature`.

Once your pull request has been merged, you can delete your branch.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
