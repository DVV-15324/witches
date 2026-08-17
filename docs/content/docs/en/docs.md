---
title: "Witches Documentation"
weight: 1
bookCollapseSection: true
---

## Quick Start

### 1. Installation

```bash
go install github.com/DVV-15324/witches@latest
```

### 2. Create New Project

```bash
witches create example
cd example
```

### 3. Configure

Edit `witches.env`

### 4. Database Setup

```bash
witches database generate
```

### 5. Init & Install Dependencies

```bash
witches init      # Generate template files
witches install   # Install dependencies
```

### 6. Run the Server

```bash
witches run
```

---

## Commands

| Command | Description |
|---------|-------------|
| `witches create <name>` | Create a new project |
| `witches init` | Generate template files and structure |
| `witches install` | Install dependencies |
| `witches run` | Start the server |
| `witches migrate up` | Apply all migrations |
| `witches migrate down` | Rollback last migration |
| `witches database generate` | Generate database URL |
| `witches add <your_service>` | Add a new service  |

---

## License

MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <br>
  Made with ❤️ by <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>