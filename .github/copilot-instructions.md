# Gehør - Interval Training and Note Reading Application

Gehør is a Go-based music notation application for interval training and note reading. It uses Raylib for graphics rendering and displays musical scores using the Leland music font and SMUFL notation standard.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### System Requirements and Setup
First, install the required system dependencies for building the Raylib graphics library:
- `sudo apt update`
- `sudo apt install -y libwayland-dev libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libxkbcommon-dev libgl1-mesa-dev`

### Git Submodules (CRITICAL)
The application requires git submodules for fonts and music notation standards. ALWAYS initialize submodules:
- `git submodule update --init --recursive` -- takes ~3 seconds

### Build and Test Process
- **NEVER CANCEL**: Build takes 46 seconds from clean cache. NEVER CANCEL. Set timeout to 120+ seconds.
- `go build -v ./...` -- takes ~46 seconds from clean. Set timeout to 120+ seconds. 
- `go build -v .` -- build main application, takes ~1 second after first build
- `go test -v ./...` -- takes ~2 seconds (no test files exist, which is normal)
- `gofmt -l .` -- verify code formatting, takes <1 second

### Running the Application
- The application requires a display to run (it's a GUI application using Raylib)
- `./gehoer` -- starts the music notation GUI application
- Application window size: 1200x800 pixels
- Application loads music scores from `assets/scores/lisa_gikk_til_skolen.json`
- Uses Leland music font from `assets/fonts/Leland/Leland.otf`

### Validation Scenarios
After making changes, ALWAYS test:
1. **Build validation**: Run `go build -v .` and ensure it completes successfully
2. **Resource loading**: Verify the application can access:
   - `assets/scores/lisa_gikk_til_skolen.json` (sample score)
   - `assets/fonts/Leland/leland_metadata.json` (font metadata)  
   - `assets/fonts/Leland/Leland.otf` (music font file)
   - `external/smufl/` directory (music notation standard)
3. **Code formatting**: Run `gofmt -l .` to ensure proper formatting
4. **Module validation**: Run `go mod tidy` to verify module dependencies

## Project Structure

### Key Directories and Files
- `main.go` -- Entry point, creates 1200x800 Raylib window and starts game loop
- `game/` -- Main game logic and initialization (1 file)
- `engraver/` -- Music notation engraving system (6 files)
- `musicfont/` -- Music font loading and glyph handling (5 files)
- `music/` -- Music score parsing and data structures (2 files)
- `camera/` -- 2D camera controller for navigation (1 file)
- `grid/` -- Background grid rendering (1 file)
- `svg/`, `smufl/`, `localization/`, `metadata/`, `settings/`, `units/` -- Support modules (1 file each)

### External Dependencies
- `assets/fonts/Leland/` -- Git submodule containing Leland music font
- `external/smufl/` -- Git submodule containing SMUFL music notation standard
- `github.com/gen2brain/raylib-go/raylib v0.55.1` -- Raylib Go bindings for graphics

### Configuration Files
- `go.mod` -- Go module definition with Go 1.24 requirement
- `.github/workflows/go.yml` -- CI pipeline that builds and tests with Go 1.20+ 
- `.gitmodules` -- Git submodule definitions for fonts and SMUFL

## Build Timing and Critical Warnings

### **NEVER CANCEL BUILD OPERATIONS**
- **Initial build from clean cache**: 46 seconds. NEVER CANCEL. Set timeout to 120+ seconds.
- **Subsequent builds**: 1 second after first successful build
- **Submodule initialization**: 3 seconds for both font and SMUFL submodules
- **Tests**: 2 seconds (no test files, reports "no test files" for all packages)
- **Code formatting check**: <1 second

### Common Commands with Timing
```bash
# NEVER CANCEL: First build takes 46 seconds
time go build -v ./...  # 46 seconds from clean, set timeout 120+

# Quick rebuild after changes  
go build -v .  # ~1 second

# Test all packages
go test -v ./...  # ~2 seconds

# Initialize required submodules
git submodule update --init --recursive  # ~3 seconds

# Format check
gofmt -l .  # <1 second
```

## Dependencies and Requirements

### Go Version
- Requires Go 1.24+ (current module version)
- CI uses Go 1.20+ but 1.24+ recommended for development

### System Libraries (Linux)
All required for building Raylib graphics:
- `libwayland-dev` -- Wayland display protocol development files
- `libx11-dev` -- X11 development files  
- `libxrandr-dev` -- X11 RandR extension development files
- `libxinerama-dev` -- X11 Xinerama extension development files
- `libxcursor-dev` -- X11 cursor development files
- `libxi-dev` -- X11 Input extension development files
- `libxkbcommon-dev` -- XKB common development files
- `libgl1-mesa-dev` -- OpenGL development files

### Asset Dependencies
- Leland music font (git submodule at `assets/fonts/Leland/`)
- SMUFL music notation standard (git submodule at `external/smufl/`)
- Sample score JSON file (`assets/scores/lisa_gikk_til_skolen.json`)

## Common Troubleshooting

### Build Failures
- If build fails with missing header files, install all system dependencies listed above
- If submodules are empty, run `git submodule update --init --recursive`
- If Go version conflicts, ensure Go 1.24+ is installed

### Runtime Failures  
- Application requires a display to run (GUI application)
- In headless environments, application will crash with GLFW display errors (expected)
- Missing font or score files will cause panic on startup

## CI Pipeline
- Located at `.github/workflows/go.yml`
- Runs on push/PR to main branch
- Uses Go 1.20 minimum
- Executes `go build -v ./...` and `go test -v ./...`
- Does NOT test GUI functionality (headless CI environment)