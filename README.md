# Gomon  

A simple file watcher that automatically restarts your application when code changes are detected.  
For now, it can only watch one file at a time.  

## Installation  

```bash
go install github.com/srivatsa-bot/gomon@latest
```

Add this path to your `.bashrc`:  

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

source the `.bashrc` file to save the changes:  

```bash
source ~/.bashrc
```

## Features  

- 🔄 Auto-reloads your application when file changes are detected (on save)  
- 🌐 Supports Go, JavaScript, and Python  
- 💻 Works only on UNIX machines.(Note: On macOS, you may encounter issues when running applications that involve ports, such as HTTP servers).

## Usage  

### Watch a Go file (default):  
```bash
gomon server.go
```

### Watch other files:  
```bash
gomon server.py
gomon server.js
```

### Show version:  
```bash
gomon --version
```

### Show help:  
```bash
gomon --help
```

---
Go to test branch to get the binaries directly.

Feel free to contribute and improve **Gomon**! 🚀
