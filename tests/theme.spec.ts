import { expect, test } from "@playwright/test";

test.describe("theme toggle", () => {
    test("toggles dark mode class and persists choice", async ({ page }) => {
        await page.goto("/");

        const html = page.locator("html");
        await expect(html).not.toHaveClass(/dark/);

        await page.locator("#theme-toggle").click();

        await expect(html).toHaveClass(/dark/);

        const stored = await page.evaluate(() =>
            localStorage.getItem("theme"),
        );
        expect(stored).toBe("dark");
    });
});
