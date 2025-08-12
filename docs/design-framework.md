# 🏗 Design Language Framework

This document defines the design framework used to structure the OneDark palette.

## Table of Contents
- [🏗 Design Language Framework](#-design-language-framework)
  - [Table of Contents](#table-of-contents)
  - [1. Palette Structure](#1-palette-structure)
  - [2. Color Roles \& Semantics](#2-color-roles--semantics)
    - [Text, Subtext1, Subtext0](#text-subtext1-subtext0)
    - [Overlay2, Overlay1, Overlay0](#overlay2-overlay1-overlay0)
    - [Surface2, Surface1, Surface0](#surface2-surface1-surface0)
    - [Base, Mantle, Crust](#base-mantle-crust)
  - [3. Layered Composition](#3-layered-composition)

## 1. Palette Structure

OneDark’s palette is made of:
- 10 **accent colors** – bright, semantic colors for syntax, highlights, and emphasis.
- 12 **functional colors** – a monochromatic range for building layered UIs.

The functional set runs from the brightest text to the deepest background:

```text
Text → Subtext1 → Subtext0 → Overlay2 → Overlay1 → Overlay0 → Surface2 → Surface1 → Surface0 → Base → Mantle → Crust
```

Think of it as a sliding scale: from high-contrast elements you read or interact with, down to background layers that hold everything together.

## 2. Color Roles & Semantics

### Text, Subtext1, Subtext0

Define typographic hierarchy:
- `Text` – main foreground (highest contrast)
- `Subtext1` – secondary text
- `Subtext0` – tertiary or muted text

### Overlay2, Overlay1, Overlay0

Light overlays for popups, tooltips, diff highlights, and temporary layers:
- `Overlay2` – highest emphasis
- `Overlay0` – softest overlay

### Surface2, Surface1, Surface0

Backgrounds for UI components like panels, cards, and sidebars:
- `Surface2` – most elevated surface
- `Surface0` – flattest, least prominent

### Base, Mantle, Crust

The foundation of the UI:
- `Base` – main canvas or editor background
- `Mantle` – secondary background (sidebars, supporting panels)
- `Crust` – framing or structural background (status bars, outer edges)

## 3. Layered Composition

OneDark is built like a stage set — from bright foregrounds to deep backdrops:
- **Foreground**: Text, Subtext, Overlay – content and interactive elements.
- **Midground**: Surface – containers and grouped sections.
- **Background**: Base, Mantle, Crust – the structure that holds everything.

This layering ensures contrast, depth, and a sense of order in every interface.