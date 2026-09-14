# go

Various tools written in GO.

## Repository Contents

### Core Modules

#### **audio**
Audio playback utilities using the Beep library.
- **player.go** - Audio player with support for:
  - Loading and playing OGG Vorbis audio files
  - Loading audio cues from directories with prefix organization
  - Streaming audio (looped or one-shot playback)
  - Random audio selection from multiple cues
  - Speaker initialization and management

#### **convo**
Conversation and dialogue tree system.
- **conversation.go** - Conversation node and dialogue management with:
  - Conversation nodes with NPC text, player options, and effects
  - Variable-driven dialogue branching (using govaluate expressions)
  - Opening branches with conditional evaluation
  - Node manipulation (add/remove options and effects, reorder options)
- **parse.go** - Parsing conversation data
- **option.go** - Conversation option handling
- **play.go** - Conversation execution and playback

#### **core**
Core utilities and types.
- **coloredChar.go** - Colored character representation for terminal rendering

#### **cview**
Terminal-based user interface toolkit (forked from tview).
- Complete UI widget library including:
  - Input forms with fields, dropdowns, checkboxes, and buttons
  - Text views and navigable multi-color displays
  - Selectable lists with context menus
  - Modal dialogs
  - Progress bars (horizontal and vertical)
  - Layout managers (Grid, Flexbox, tabbed panels)
  - Table and tree views
  - Draggable and resizable windows
  - Application wrapper with mouse support

#### **dice_curve**
RPG character sheet and game mechanics system (compatible with GURPS-like rules).
- **charsheet.go** - Character sheet management with:
  - Attributes (Strength, Dexterity, Intelligence, Health, Will, Perception, etc.)
  - Character point tracking and allocation
  - Skill management with levels and defaults
  - Resource tracking (Hit Points, Fatigue Points)
  - Modifier system for temporary and persistent effects
  - Active defense calculations (dodge, block, parry)
- **attributes.go** - Attribute definitions and calculations
- **dice.go** - Dice rolling mechanics
- **skills.go** - Skill definitions, leveling, and calculations
- **melee.go** - Melee combat mechanics
- **rpg.go** - General RPG rules and calculations
- **rules.go** - Game rule implementations

#### **fview**
Extended UI components.
- **modals.go** - Modal dialog implementations

#### **fxtools**
General utility library with graphics, math, and functional programming tools.
- **anim.go** - Animation utilities
- **brushes.go** - Drawing/rendering brush implementations
- **circle.go** - Circle drawing and manipulation
- **colors.go** - Color utilities and palettes
- **colorcodes.go** - Color code conversions
- **cp437.go** - Code page 437 character set utilities
- **dda_raycast.go** - DDA (Digital Differential Analyzer) raycasting for graphics
- **easing.go** - Easing functions for animations
- **files.go** - File I/O utilities
- **functional.go** - Functional programming helpers
- **gameloop.go** - Game loop management
- **interval.go** - Time interval utilities
- **logic.go** - Boolean logic and bitwise operations
- **logic_test.go** - Tests for logic utilities
- **math.go** - Math helpers
- **menuitem.go** - Menu item representation
- **murmur.go** - Murmur hash implementation
- **net.go** - Networking utilities
- **predicate.go** - Predicate/filter functions
- **string_flags.go** - String flag parsing
- **strings.go** - String manipulation utilities
- **terminal.go** - Terminal utilities
- **tuple.go** - Tuple type for pairs of values

#### **geometry**
Geometric algorithms and spatial data structures for game development and pathfinding.
- **point.go** - 2D point operations
- **rect.go** - Rectangle operations and intersection tests
- **distance.go** - Distance calculations (Euclidean, Manhattan, etc.)
- **astar.go** - A* pathfinding algorithm implementation
- **dijkstra.go** - Dijkstra's algorithm for shortest paths
- **jps.go** - Jump Point Search pathfinding (optimization of A*)
- **breadthfirst.go** - Breadth-first search pathfinding
- **fov.go** - Field of View calculations (likely shadowcasting)
- **camera.go** - Camera/viewport management
- **neighbors.go** - Neighbor finding utilities
- **pathrange.go** - Path range and reachability calculations
- **cc.go** - Connected components or constraint checking
- **heap.go** - Priority queue heap implementation
- **utils.go** - General geometric utilities

#### **recfile**
Record file format parser and serializer (GNU recutils format support).
- **rec.go** - Main record file parser with:
  - Support for GNU recfiles format (field-based record database)
  - Schema parsing and validation
  - Record type definitions
  - Type system (int, string, enum, references)
  - Multi-line field support (escaped newlines)
  - CSV export capability
- **encoder.go** - Record file encoding and writing
- **decoder.go** - Record file decoding and reading
- **schema.go** - Schema definition and validation
- **repo.go** - Repository management for record files
- **string.go** - String field handling

#### **textiles**
Tile and sprite management system.
- **tiles.go** - Tile definitions and management
- **icons.go** - Icon/sprite handling
- **iconRecord.go** - Record-based icon definitions
- **palette.go** - Color palette management
- **zones.go** - Spatial zone/region management

### Example Files

- **main.go** - Example program demonstrating recfile loading and processing using cview
- **test.rec** - Sample recfile for testing

### Build Files

- **go.mod** - Go module definition with dependencies
- **go.sum** - Go module checksum verification
- **.gitignore** - Git ignore patterns
