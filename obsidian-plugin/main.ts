import { App, Plugin, PluginSettingTab, Setting, TAbstractFile, TFile, Notice, requestUrl } from "obsidian";

interface CMSSettings {
  apiUrl: string;
  apiToken: string;
  feedFolder: string;
  feedFile: string;
}

const DEFAULT_SETTINGS: CMSSettings = {
  apiUrl: "https://api.sahil.im",
  apiToken: "",
  feedFolder: "",
  feedFile: "feed.md",
};

const DEBOUNCE_MS = 2000;
const MAX_RETRIES = 3;

interface QueueItem {
  slug: string;
  content: string;
}

export default class SahilCMS extends Plugin {
  settings: CMSSettings = DEFAULT_SETTINGS;
  private debounceTimer: number | null = null;
  private syncQueue: QueueItem[] = [];
  private syncing = false;

  async onload() {
    await this.loadSettings();

    this.addSettingTab(new CMSSettingTab(this.app, this));

    this.registerEvent(
      this.app.vault.on("modify", (file: TAbstractFile) => {
        if (!(file instanceof TFile)) return;
        if (!this.settings.apiToken) return;
        if (!this.isFeedFile(file)) return;

        const frontmatter = this.app.metadataCache.getFileCache(file)?.frontmatter;
        if (frontmatter && frontmatter.publish === false) return;

        if (this.debounceTimer) window.clearTimeout(this.debounceTimer);
        this.debounceTimer = window.setTimeout(() => this.enqueueSync(file), DEBOUNCE_MS);
      })
    );

    this.addCommand({
      id: "sync-feed",
      name: "Sync feed now",
      callback: async () => {
        const targetPath = this.buildFeedPath();
        const feed = this.app.vault.getAbstractFileByPath(targetPath);
        if (feed instanceof TFile) {
          await this.syncFile(feed);
        } else {
          new Notice(`Feed file "${targetPath}" not found`);
        }
      },
    });
  }

  async enqueueSync(file: TFile) {
    const content = await this.app.vault.read(file);
    const slug = file.name.replace(/\.md$/i, "").toLowerCase();
    console.log(`[sahil-cms] queued ${slug} (${content.length} chars)`);
    this.syncQueue.push({ slug, content });
    this.processQueue();
  }

  async processQueue() {
    if (this.syncing || this.syncQueue.length === 0) return;
    this.syncing = true;

    const item = this.syncQueue.shift()!;
    const url = `${this.settings.apiUrl}/api/cms/feed`;
    let success = false;
    let lastErr = "";

    for (let attempt = 0; attempt < MAX_RETRIES && !success; attempt++) {
      try {
        console.log(`[sahil-cms] POST ${url} (attempt ${attempt + 1})`);
        const resp = await requestUrl({
          url,
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${this.settings.apiToken}`,
          },
          body: JSON.stringify(item),
        });
        console.log(`[sahil-cms] response status=${resp.status}`, resp.json || "");
        if (resp.status >= 200 && resp.status < 300) {
          success = true;
        } else {
          lastErr = `HTTP ${resp.status}`;
        }
      } catch (e: unknown) {
        lastErr = e instanceof Error ? e.message : String(e);
        console.error(`[sahil-cms] attempt ${attempt + 1} failed:`, e);
        if (attempt < MAX_RETRIES - 1) {
          await this.sleep((attempt + 1) * 2000);
        }
      }
    }

    if (success) {
      new Notice(`Synced feed → ${this.settings.apiUrl}`);
    } else {
      new Notice(`Sync failed: ${lastErr}`);
      this.syncQueue.unshift(item);
    }

    this.syncing = false;
    this.processQueue();
  }

  async syncFile(file: TFile) {
    const content = await this.app.vault.read(file);
    const slug = file.name.replace(/\.md$/i, "").toLowerCase();
    this.syncQueue = this.syncQueue.filter((q) => q.slug !== slug);
    this.syncQueue.push({ slug, content });
    if (!this.syncing) this.processQueue();
  }

  buildFeedPath(): string {
    const folder = this.settings.feedFolder.replace(/^\/|\/$/g, "");
    const file = this.settings.feedFile;
    return folder ? `${folder}/${file}` : file;
  }

  isFeedFile(file: TFile): boolean {
    const targetPath = this.buildFeedPath();
    return file.path === targetPath;
  }

  sleep(ms: number) {
    return new Promise((r) => setTimeout(r, ms));
  }

  async loadSettings() {
    this.settings = Object.assign({}, DEFAULT_SETTINGS, await this.loadData());
  }

  async saveSettings() {
    await this.saveData(this.settings);
  }
}

class CMSSettingTab extends PluginSettingTab {
  plugin: SahilCMS;

  constructor(app: App, plugin: SahilCMS) {
    super(app, plugin);
    this.plugin = plugin;
  }

  display() {
    const { containerEl } = this;
    containerEl.empty();

    containerEl.createEl("h2", { text: "Sahil CMS Settings" });

    new Setting(containerEl)
      .setName("API URL")
      .setDesc("Your sahil-api base URL")
      .addText((text) =>
        text
          .setPlaceholder("https://api.sahil.im")
          .setValue(this.plugin.settings.apiUrl)
          .onChange(async (v) => {
            this.plugin.settings.apiUrl = v;
            await this.plugin.saveSettings();
          })
      );

    new Setting(containerEl)
      .setName("API Token")
      .setDesc("Bearer token for authentication")
      .addText((text) =>
        text
          .setPlaceholder("cms_...")
          .setValue(this.plugin.settings.apiToken)
          .onChange(async (v) => {
            this.plugin.settings.apiToken = v;
            await this.plugin.saveSettings();
          })
      );

    new Setting(containerEl)
      .setName("Feed folder")
      .setDesc("Folder path inside your vault (e.g. site-items). Leave empty for root.")
      .addText((text) =>
        text
          .setPlaceholder("site-items")
          .setValue(this.plugin.settings.feedFolder)
          .onChange(async (v) => {
            this.plugin.settings.feedFolder = v;
            await this.plugin.saveSettings();
          })
      );

    new Setting(containerEl)
      .setName("Feed filename")
      .setDesc("Filename to watch inside the folder")
      .addText((text) =>
        text
          .setPlaceholder("feed.md")
          .setValue(this.plugin.settings.feedFile)
          .onChange(async (v) => {
            this.plugin.settings.feedFile = v;
            await this.plugin.saveSettings();
          })
      );

    new Setting(containerEl)
      .setName("Watched path")
      .setDesc(this.plugin.buildFeedPath())
      .setClass("setting-item-info");
  }
}
