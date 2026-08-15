import { beforeEach, describe, expect, it, vi } from "vitest";

const notices: string[] = [];

vi.mock("obsidian", () => {
  class Plugin {
    settings: Record<string, unknown> = {};
    async loadData() {
      return {};
    }
    async saveData(_data: unknown) {}
    addSettingTab() {}
    registerEvent() {}
    addCommand() {}
  }
  class Notice {
    constructor(msg: string) {
      notices.push(msg);
    }
  }
  return {
    Plugin,
    Notice,
    App: class {},
    PluginSettingTab: class {
      containerEl = { empty() {}, createEl() {} };
      constructor(_app: unknown, _plugin: unknown) {}
    },
    Setting: class {
      constructor(public containerEl: unknown) {}
      setName() {
        return this;
      }
      setDesc() {
        return this;
      }
      setPlaceholder() {
        return this;
      }
      setValue() {
        return this;
      }
      setClass() {
        return this;
      }
      addText(fn: (t: unknown) => void) {
        const text = {
          setPlaceholder() {
            return text;
          },
          setValue() {
            return text;
          },
          onChange() {
            return text;
          },
        };
        fn(text);
        return this;
      }
    },
    TAbstractFile: class {},
    TFile: class {},
    requestUrl: vi.fn(),
  };
});

import SahilCMS from "./main.ts";
import { requestUrl } from "obsidian";

function makePlugin(): SahilCMS {
  const plugin = new SahilCMS({} as any, {} as any);
  (plugin as any).app = {
    vault: {
      read: vi.fn().mockResolvedValue("# Hello"),
      getAbstractFileByPath: vi.fn(),
      on: vi.fn(),
    },
    metadataCache: { getFileCache: vi.fn() },
  };
  plugin.settings = {
    apiUrl: "https://api.sahil.im",
    apiToken: "tok",
    feedFolder: "",
    feedFile: "feed.md",
  };
  return plugin;
}

describe("SahilCMS", () => {
  let plugin: SahilCMS;

  beforeEach(() => {
    notices.length = 0;
    (requestUrl as any).mockReset();
    (requestUrl as any).mockResolvedValue({ status: 200, json: {} });
    plugin = makePlugin();
  });

  it("loadSettings merges defaults", async () => {
    await plugin.loadSettings();
    expect(plugin.settings.apiUrl).toBe("https://api.sahil.im");
    expect(plugin.settings.feedFile).toBe("feed.md");
  });

  it("buildFeedPath without folder returns filename", () => {
    plugin.settings.feedFolder = "";
    expect(plugin.buildFeedPath()).toBe("feed.md");
  });

  it("buildFeedPath trims slashes and joins folder", () => {
    plugin.settings.feedFolder = "/site-items/";
    expect(plugin.buildFeedPath()).toBe("site-items/feed.md");
  });

  it("isFeedFile matches the target path", () => {
    plugin.settings.feedFolder = "site";
    expect(plugin.isFeedFile({ path: "site/feed.md" } as any)).toBe(true);
    expect(plugin.isFeedFile({ path: "other.md" } as any)).toBe(false);
  });

  it("enqueueSync reads the file and derives the slug", async () => {
    (plugin as any).syncing = true;
    await plugin.enqueueSync({ name: "Feed.md" } as any);
    expect((plugin as any).syncQueue).toEqual([{ slug: "feed", content: "# Hello" }]);
  });

  it("syncFile queues and de-dupes by slug", async () => {
    (plugin as any).syncing = true;
    await plugin.syncFile({ name: "feed.md" } as any);
    await plugin.syncFile({ name: "feed.md" } as any);
    expect((plugin as any).syncQueue).toEqual([{ slug: "feed", content: "# Hello" }]);
  });

  it("processQueue posts and succeeds", async () => {
    (plugin as any).syncQueue.push({ slug: "feed", content: "# Hi" });
    await plugin.processQueue();
    expect(requestUrl).toHaveBeenCalledTimes(1);
    expect((plugin as any).syncQueue).toHaveLength(0);
    expect(notices[0]).toContain("Synced feed");
  });

  it("processQueue retries and re-queues on failure", async () => {
    (requestUrl as any).mockResolvedValue({ status: 500, json: {} });
    // guard against the recursive re-process loop: only run one full pass
    const realProcess = plugin.processQueue.bind(plugin);
    let calls = 0;
    plugin.processQueue = (async () => {
      calls++;
      if (calls === 1) await realProcess();
    }) as any;

    (plugin as any).syncQueue.push({ slug: "feed", content: "# Hi" });
    await plugin.processQueue();
    expect(requestUrl).toHaveBeenCalledTimes(3);
    expect((plugin as any).syncQueue).toHaveLength(1);
    expect(notices[0]).toContain("Sync failed");
  });
});