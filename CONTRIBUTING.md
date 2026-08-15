# Contributing

First off: thank you for considering contributing to a website that describes
itself in a single sentence and occasionally turns into a blog. You are the
hero this repository doesn't deserve.

## Ground rules

- **It's a personal site.** The maintainer lives here. Treat it like someone's
  living room: take your shoes off, don't rearrange the furniture, and don't
  be weird about the hero status line.
- **Small PRs > heroic ones.** A 400-line "while I was here" refactor will be
  reviewed with the enthusiasm of a morning meeting.
- **Test the thing.** If you changed how a fetcher parses XML, run the tests:
  `task api:test`. If you changed the plugin, `task obsidian:test`. If you
  changed a `biome.json`, question your life choices.
- **Follow the existing vibe.** Astro, Tailwind, tabs, double quotes, no
  comments unless they earn their keep. When in doubt, `task lint` will gently
  (or not-so-gently) remind you.

## Getting started

1. Fork the repo.
2. `git clone` your fork and run:
   ```sh
   git submodule update --init   # the blog posts live in a submodule, because of course they do
   bun install                   # site deps
   cd api && go mod download     # Go deps
   cd ../obsidian-plugin && npm install
   ```
3. Create a branch with a name that describes what you're doing, not your
   emotional state: `fix/letterboxd-year` > `please-dont-break-it-again`.
4. Make your change. Touch only the parts you meant to touch. Resist the urge
   to reformat the entire file because your editor "helped".
5. Run the checks:
   ```sh
   task lint       # site: biome + astro check
   task api:lint   # go vet + gofmt
   task api:test   # go tests + coverage
   task obsidian:test
   ```
6. Open a PR. Link it to an issue if one exists. Add a one-line summary of
   *what* you did and *why*. "Fixed it" is not a summary.

## Submitting PRs

- Keep the diff focused. If your PR is more than ~150 changed lines, consider
  whether it's actually two PRs.
- Don't be alarmed if feedback comes slowly. There is one maintainer, and they
  occasionally go outside.
- If a CI check fails, fix it before you ask. The robots can smell fear.
- Celebrate small wins: every merged PR moves the coverage number slightly
  closer to a number that doesn't embarrass us.

## Code style

- **Go:** `gofmt` is law. `go vet` is the constable.
- **TypeScript / Astro:** Biome is the referee. Don't argue with the referee.
- **CSS:** Tailwind utilities, in a `global.css` that stays small.
- **Commit messages:** concise, present tense, maybe a conventional prefix
  (`feat:`, `fix:`, `optimise:`). Emojis permitted but rationed.

## Testing

- Go API: `task api:test` — real unit tests, mocked HTTP, no real requests.
- Obsidian plugin: `task obsidian:test` — vitest with a lovingly crafted
  `obsidian` mock.
- Site e2e: `task test` — Playwright against the dev server.

If you add a feature, add a test. If you can't test it, at least whisper to
yourself that you should have.

## Questions?

Open an issue. There are no stupid questions, only stubs left waiting for
answers. That's a joke about this repo having at least one `coming soon`.

Happy hacking.