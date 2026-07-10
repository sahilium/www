-- D1 schema for Obsidian CMS
-- Run this in Cloudflare Dashboard → D1 → Query

CREATE TABLE IF NOT EXISTS feed_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  slug TEXT UNIQUE NOT NULL DEFAULT 'feed',
  title TEXT NOT NULL DEFAULT '',
  content TEXT NOT NULL,
  metadata TEXT NOT NULL DEFAULT '{}',
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_feed_slug ON feed_entries(slug);
CREATE INDEX IF NOT EXISTS idx_feed_updated ON feed_entries(updated_at DESC);
