# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project does not yet follow [Semantic Versioning](https://semver.org).

This changelog covers the current Astro-based site, from the Astro migration
onwards. Earlier versions of the site are not tracked here.

## [Unreleased]

### Added
- Astro-based site, migrating from the previous (SCSS/Jekyll-style) version.
- Command palette.
- Hero section with dynamic status lines; new favicon.
- Japanese microseasons (七十二候) — the `/koo` page with dynamic kō logic.
- Click-sound interaction.
- IndieWeb: h-card and xxiivv webring membership.
- Selection highlighting.
- Initial `sahil.im` Go API (`api/`).
- "Now" page backed by the API (GitHub commit + portfolio graph).
- Goodreads integration (replacing GitHub commit + Hardcover).
- Obsidian-managed feed page (publish posts from Obsidian to the CMS API).
- Content moved into a git submodule (away from GitHub).
- Playwright UI tests and test specs.
- Restructured "Now" page with moon phase and images.
- Dynamic Washi theme colors for the `koo` page.
- Taskfiles (Task) covering the site, API, and Obsidian plugin.
- Lighthouse CI with live badge scores.
- Lint CI (Biome + Astro check + Go).

### Changed
- Design and typography overhaul.
- Footer improvements.
- Greetings and statuses now served from the API.
- Feed links and post cleanup.
- Lightened dark-mode background.
- README and specifications updated.

### Fixed
- `@tailwindcss/vite` downgrade / replacement with PostCSS for the Pages build.
- CTA behaviour on navigation.
- Music artist handling.
- esbuild bump in the Obsidian plugin.
- Lighthouse and Lint CI action/typecheck issues.

### Performance
- Core Web Vitals optimizations (`optimise: 1`–`5`, including font loading).

## Older versions

Changes from before the Astro migration (2026-06-19) are not listed here.
