---
title: "Get start Witches"
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

Edit `witches.env` file with your database credentials.

### 4. Generate database URL from env

```bash
witches database generate
```

### 5. Run database migrations

```bash
witches migrate up
```

### 6. Init & Install Dependencies

```bash
witches init      # Generate template files
witches install   # Install dependencies
```

### 7. Run the Server

```bash
witches run
```
---

<div align="center">
  <br>
  Made with ❤️ by <a href="https://github.com/DVV-15324">DVV-15324</a>
  <br><br>
</div>