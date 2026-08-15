# Security Policy

## Supported Versions

The deployed site (the `main` branch) is the only actively supported release.

| Version | Supported          |
| ------- | ------------------ |
| `main`  | :white_check_mark: |

## Reporting a Vulnerability

Please report security issues **privately** — do not open a public issue.

- **GitHub:** use the [private vulnerability reporting](https://github.com/sahilium/www/security/advisories/new) flow on the repository.
- **PGP:** for sensitive reports, encrypt to the site's PGP key, published at [`/pgp.asc`](https://sahil.im/pgp.asc).

Please include:
- The affected component (site, API, or Obsidian plugin).
- A minimal description of the issue and, if possible, a reproduction.
- Impact and any suggested fix, if known.

We aim to acknowledge reports within a few days and keep you updated on progress. Please allow time for a fix before public disclosure.

## Scope

This policy applies to the [Astro site](https://sahil.im), the [Go API](api/), and the [Obsidian plugin](obsidian-plugin/).

## Out of scope

- Dependencies managed by your package manager (report these upstream).
- General questions or feature requests (open a normal issue or discussion instead).
