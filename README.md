# raf-orphan-cleaner

A small command-line tool that removes orphan Fujifilm **`.RAF`** raw files when there is no matching JPEG with the same base name in the same folder.

Typical use: after culling JPEGs from a camera dump (keeping only the keepers), delete the leftover RAW files that no longer have a JPEG partner.

## What it does

For a single folder (non-recursive), the tool:

1. Collects every `.JPG` / `.JPEG` file (case-insensitive).
2. Looks at every `.RAF` file (case-insensitive).
3. Treats a `.RAF` as **orphan** if no JPEG shares the same stem (file name without extension).
4. **Dry-runs by default** — lists orphans and prints a summary.
5. Deletes orphans only when you pass **`-d`**.

Example pairing:

| Kept | Reason |
|------|--------|
| `FUJI0285.RAF` | `FUJI0285.JPG` exists |
| `FUJI0527.RAF` | no `FUJI0527.JPG` / `.JPEG` → orphan |

## Safety

- **Dry-run is the default.** Nothing is deleted unless you pass `-d`.
- Scans **one folder only** (does not walk subdirectories).
- Matching is by file stem only; it does not inspect EXIF or image content.
- Always run a dry-run first and review the file list before deleting.

Deletion is permanent (`os.Remove`). There is no trash / recycle-bin step.

## Requirements

- [Go](https://go.dev/dl/) **1.22** or newer (to build from source)
- No third-party Go modules (stdlib only)

## Build

```bash
cd raf-orphan-cleaner
go build -o raf-orphan-cleaner .
```

On Windows:

```bash
go build -o raf-orphan-cleaner.exe .
```

### Cross-compile

From macOS or Linux, build a Windows executable:

```bash
GOOS=windows GOARCH=amd64 go build -o raf-orphan-cleaner.exe .
```

Apple Silicon → Intel macOS binary:

```bash
GOOS=darwin GOARCH=amd64 go build -o raf-orphan-cleaner .
```

## Usage

```text
raf-orphan-cleaner [flags] [folder]
```

| Argument / flag | Meaning |
|-----------------|--------|
| `folder` | Directory to scan (default: current directory) |
| `-d` | Actually delete orphan `.RAF` files |

### Dry-run (safe)

```bash
./raf-orphan-cleaner /path/to/photos
```

Lists each orphan file name, then a summary:

```text
FUJI0527.RAF
FUJI0528.RAF
would delete 100 orphan .RAF (87 kept) in /path/to/photos
```

### Delete orphans

```bash
./raf-orphan-cleaner -d /path/to/photos
```

Prints `deleted <name>` for each removed file, then:

```text
deleted 100 orphan .RAF (87 kept) in /path/to/photos
```

### Help

```bash
./raf-orphan-cleaner -h
```

## How matching works

- Extensions compared case-insensitively: `.raf`, `.RAF`, `.jpg`, `.JPG`, `.jpeg`, `.JPEG`.
- Stem match is case-insensitive (`fuji0285.raf` matches `FUJI0285.JPG`).
- Only files directly in the given folder are considered.
- Non-RAF / non-JPEG files are ignored.

## Project layout

```text
raf-orphan-cleaner/
├── main.go       # CLI + orphan detection
├── go.mod        # module raf-orphan-cleaner
├── README.md
└── .gitignore
```

## License

Use and modify freely for personal or project use unless you add a more specific license later.
