# Gehør - Interval Training and Note Reading Application

Gehør is a Go-based music education application for interval training and note reading (Norwegian: "Intervalltrening og notelesing"). It uses raylib for graphics rendering and displays musical notation with interactive features.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Bootstrap and Build the Repository
Run these commands in exact order:

1. **Initialize git submodules** (required for fonts and SMUFL specification):
   ```bash
   git submodule update --init --recursive
   ```
   - Takes ~10-15 seconds
   - Downloads Leland music font and SMUFL specification

2. **Install system dependencies** (required for raylib compilation):
   ```bash
   sudo apt update && sudo apt install -y libwayland-dev libxkbcommon-dev xorg-dev
   ```
   - Required for GLFW/Wayland support in raylib
   - Takes ~60-120 seconds depending on network

3. **Download Go dependencies**:
   ```bash
   go mod download
   ```
   - Takes ~3 seconds
   - Downloads raylib-go and other Go dependencies

4. **Build the application**:
   ```bash
   go build -v ./...
   ```
   - **NEVER CANCEL: Build takes 40-60 seconds. NEVER CANCEL. Set timeout to 90+ minutes.**
   - Compiles C libraries (raylib, GLFW) which takes significant time
   - Creates `gehoer` binary in root directory

### Testing and Validation
```bash
go test -v ./...
```
- Takes ~1-2 seconds
- **NOTE**: No test files exist in this codebase (all packages show "[no test files]")
- This is normal and expected

### Code Quality Checks
Always run these before committing:
```bash
go vet ./...         # Static analysis - takes ~2 seconds
gofmt -l .          # Check formatting - takes ~1 second
```
- Both should produce no output if code is clean
- **CRITICAL**: The CI (.github/workflows/go.yml) will fail if these produce errors

## Application Structure

### Key Directories and Purpose
- **`main.go`** - Entry point, creates 1200x800 raylib window
- **`game/`** - Main game logic, camera controls, rendering loop
- **`engraver/`** - Musical notation rendering (staff, notes, clefs, rests)
- **`musicfont/`** - SMUFL-compliant music font handling and glyph rendering
- **`music/`** - Music data structures and JSON score loading
- **`assets/scores/`** - Sample musical scores in JSON format
- **`assets/fonts/Leland/`** - Git submodule with Leland music font
- **`external/smufl/`** - Git submodule with SMUFL specification
- **Supporting modules**: `camera/`, `grid/`, `svg/`, `settings/`, `units/`, `localization/`

### Common File Locations
```bash
# Repository root structure
ls -la
# Shows: main.go, go.mod, gehoer (binary), assets/, game/, engraver/, etc.

# Sample score file
cat assets/scores/lisa_gikk_til_skolen.json
# Norwegian children's song "Lisa gikk til skolen" in JSON format

# Main application entry
cat main.go
# Creates raylib window and starts game loop
```

## Running the Application

### **CRITICAL LIMITATION**: Graphics Display Required
The application **CANNOT** run in headless environments:
```bash
./gehoer
# ERROR: "GLFW: Failed to initialize GLFW" followed by segmentation fault
```

- **Requires X11/Wayland display server**
- **Cannot run in CI environments without virtual display**
- Application initializes 1200x800 pixel window for musical notation display

### Expected Runtime Behavior
When run with proper display:
- Opens 1200x800 window titled "Gehør"
- Displays musical staff with grid background
- Loads and renders "Lisa gikk til skolen" musical score
- Supports camera controls for navigation
- Runs at 60 FPS using raylib game loop

## Validation Scenarios

Since the application requires graphics display and no tests exist:

### Build Validation (Always Required)
```bash
# Full validation sequence
git submodule update --init --recursive
go mod download
go build -v ./...
go vet ./...
gofmt -l .
```

### Code Change Validation
1. **Always build after changes**: `go build -v ./...`
2. **Check for compilation errors** - especially in raylib integration
3. **Verify JSON score loading** - check `assets/scores/` if modifying music parsing
4. **Test font loading** - verify SMUFL font handling if modifying `musicfont/`

### CI Compatibility
The GitHub Actions workflow (`.github/workflows/go.yml`):
- Uses Go 1.20 (codebase uses Go 1.24+)
- Runs `go build -v ./...` and `go test -v ./...`
- **Does NOT install system dependencies** - will fail without manual setup
- **Does NOT initialize submodules** - will fail without manual setup

## Common Tasks and Timing

### Repository Operations
- **Clone fresh repo**: `git clone` + submodule init ~30-45 seconds
- **Clean build from scratch**: ~45-60 seconds total
- **Incremental builds**: ~5-15 seconds depending on changes

### Development Workflow
1. Make code changes
2. **Always run**: `go build -v ./...` - **NEVER CANCEL: 40-60 seconds**
3. **Always run**: `go vet ./...` - 2 seconds
4. **Before commit**: `gofmt -l .` - 1 second
5. Test manually with graphics if possible

### Key Dependencies
- **Go 1.24+** (higher than CI uses Go 1.20)
- **raylib-go v0.55.1** - Graphics and window management
- **System packages**: libwayland-dev, libxkbcommon-dev, xorg-dev
- **Git submodules**: Leland font and SMUFL specification

## Troubleshooting

### Build Failures
- **"wayland-client-core.h: No such file"** → Install libwayland-dev
- **"Failed to detect any supported platform"** → Missing xorg-dev
- **Submodule errors** → Run `git submodule update --init --recursive`

### Runtime Issues
- **Segmentation fault on startup** → No graphics display available
- **Font loading errors** → Verify submodules initialized
- **Score loading errors** → Check JSON format in `assets/scores/`

### Performance Notes
- **Builds are slow** - raylib C compilation takes 30-40 seconds minimum
- **No parallel builds** - Go builds raylib sequentially
- **Memory usage** - Loads entire font files and SMUFL specification

## CRITICAL Reminders
- **NEVER CANCEL builds** - C compilation appears to hang but is working
- **ALWAYS set 90+ minute timeouts** for build commands
- **CANNOT run application** without graphics display
- **NO TESTS EXIST** - validation must be done through building and manual testing
- **CI WILL FAIL** without dependency installation and submodule setup