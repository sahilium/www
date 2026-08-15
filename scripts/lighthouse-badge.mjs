import { spawnSync } from "node:child_process";
import { mkdirSync, writeFileSync, existsSync, readFileSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const lighthouseBin = join(root, "node_modules", ".bin", "lighthouse");
const outDir = join(root, "lighthouse");
const targetUrl = process.env.LIGHTHOUSE_URL || "https://sahil.im/";
const chromePath = process.env.CHROME_PATH || "";

const categories = ["performance", "accessibility", "best-practices", "seo"];
const presets = ["desktop", "mobile"];

function colorFor(score) {
  if (score >= 90) return "brightgreen";
  if (score >= 80) return "green";
  if (score >= 70) return "yellowgreen";
  if (score >= 60) return "yellow";
  return "red";
}

function run(preset) {
  const args = [
    targetUrl,
    "--quiet",
    "--output=json",
    `--output-path=${join(outDir, `report-${preset}.json`)}`,
    `--only-categories=${categories.join(",")}`,
    "--chrome-flags=--headless=new --no-sandbox",
  ];
  if (chromePath) args.push(`--chrome-path=${chromePath}`);
  if (preset === "desktop") args.push("--preset=desktop");

  const r = spawnSync(lighthouseBin, args, { stdio: "inherit" });
  if (r.status !== 0) {
    console.error(`Lighthouse failed for ${preset} (exit ${r.status})`);
    return null;
  }
  const report = readFileSync(join(outDir, `report-${preset}.json`), "utf8");
  const data = JSON.parse(report);
  const scores = {};
  for (const cat of categories) {
    scores[cat] = Math.round((data.categories[cat]?.score ?? 0) * 100);
  }
  return scores;
}

mkdirSync(outDir, { recursive: true });

const labels = {
  performance: "Performance",
  accessibility: "Accessibility",
  "best-practices": "Best Practices",
  seo: "SEO",
};

for (const preset of presets) {
  const scores = run(preset);
  if (!scores) process.exitCode = 1;
  else {
    console.log(`${preset}:`, scores);
    for (const cat of categories) {
      const score = scores[cat];
      const file = join(outDir, `${cat}-${preset}.json`);
      const payload = {
        schemaVersion: 1,
        label: labels[cat],
        message: String(score),
        color: colorFor(score),
      };
      const next = `${JSON.stringify(payload)}\n`;
      if (!existsSync(file) || readFileSync(file, "utf8") !== next) {
        writeFileSync(file, next);
      }
    }
  }
}