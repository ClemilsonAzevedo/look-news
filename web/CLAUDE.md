# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`web/` is the frontend of the **look-news** project. It is a React 19 + TypeScript SPA built with Vite 8 and styled with Tailwind CSS v4. A sibling `../api/` directory holds the backend (separate package). This is an early-stage scaffold — `src/` currently contains only the Vite entry (`main.tsx`), root component (`App.tsx`), and a Tailwind import (`index.css`).

## Commands

This project uses **Bun** (`bun.lock`). Use `bun install` and `bun run <script>`.

- `bun run dev` — Vite dev server with HMR
- `bun run build` — type-check (`tsc -b`) then production build. The build fails on type errors.
- `bun run lint` — Oxlint (the linter is `oxlint`, not ESLint)
- `bun run preview` — serve the production build locally

There is no test runner configured yet.

## Architecture & conventions

- **Tailwind v4** is wired through the `@tailwindcss/vite` plugin (`vite.config.ts`), not a `tailwind.config.js` / PostCSS pipeline. The only CSS file is `src/index.css`, which is just `@import "tailwindcss";`. Theme tokens, if added, belong in CSS via `@theme`, not a JS config.
- **TypeScript is strict** beyond defaults: `noUnusedLocals`, `noUnusedParameters`, `erasableSyntaxOnly` (no enums / non-erasable syntax), and `verbatimModuleSyntax` (use `import type` for type-only imports). `allowImportingTsExtensions` is on, so local imports include the `.tsx`/`.ts` extension (see `main.tsx`).
- **Components are named exports**, not default exports (`export function App()`). Oxlint enforces `react/rules-of-hooks` (error) and `react/only-export-components` (warn).
- The React Compiler is intentionally **not** enabled.

## Design system — read `DESIGN.md` before building UI

`DESIGN.md` defines the **"Sahara — Warm Minimalism"** visual language and is the source of truth for all UI work. Key constraints:

- **Warm palette only.** Primary `#c2652a` (burnt sienna), background `#faf5ee` (warm linen — never cold white), accent `#8c3c3c` (dusty rose). Even grays carry warm undertones.
- **Type pairing:** EB Garamond (serif headlines) + Manrope (sans body/labels).
- **Whitespace is the primary tool** — favor generous padding (28–32px cards), ultra-soft shadows, thin warm borders. Curated over cluttered.
