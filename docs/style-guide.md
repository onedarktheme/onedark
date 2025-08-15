# 🎨 One Dark — Style Guide

A structured reference for using the **One Dark** palette across UIs, code editors, and terminal themes.

## 1. General Principles

- **Legibility First** — Always prioritize contrast and readability over strict adherence to this guide.
- **Consistency Matters** — Keep the same color meaning across contexts (e.g., red = danger).
- **Semantic Colors** — Reserve bright accents for meaningful actions or states.

## 2. Usage Guidelines

### UI Background Layers

| Layer                | Recommended Colors             | Notes                               |
| -------------------- | ------------------------------ | ----------------------------------- |
| Primary Background   | Base                           | Main app background                 |
| Secondary Background | Mantle / Crust                 | For panels & menus                  |
| Surface Elements     | Surface0 / Surface1 / Surface2 | Use progressively lighter for depth |
| Overlays             | Overlay0 / Overlay1 / Overlay2 | Hover states, subtle emphasis       |

### Typography

| Text Role      | Color    | Notes                              |
| -------------- | -------- | ---------------------------------- |
| Primary Text   | Text     | Body copy, headlines               |
| Secondary Text | Subtext1 | Less important labels              |
| Muted Text     | Subtext0 | Disabled states, hints             |
| Links          | Blue     | Hover states can be Cyan or Purple |
| Success        | Green    | For confirmations                  |
| Warning        | Yellow   | For caution                        |
| Error          | Red      | For destructive actions            |

### Code & Syntax Highlighting

| Element             | Color           | Example             |
| ------------------- | --------------- | ------------------- |
| Keywords            | Purple          | `function`, `class` |
| Strings             | Green           | `"Hello"`           |
| Numbers / Constants | Orange / Yellow | `42`                |
| Functions & Methods | Blue            | `myFunction()`      |
| Variables           | Text            | `let name`          |
| Parameters          | Dark Yellow     | `function(a)`       |
| Comments            | Overlay2        | `// comment`        |
| Errors              | Red             | Syntax errors       |
| Warnings            | Yellow          | Deprecation         |
| Types / Classes     | Yellow          | `MyClass`           |
| Attributes          | Orange          | HTML/JSX props      |

## 3. Accessibility Notes

- Ensure contrast ratio of at least 4.5:1 for text vs background
- Avoid placing bright saturated colors against similar brightness backgrounds without enough contrast

[WCAG on contrast ratio](https://www.w3.org/WAI/WCAG21/Techniques/general/G18)