import { expect, test } from "@playwright/test";

test.describe("command palette", () => {
    test("opens via the Explore button and shows nav commands", async ({
        page,
    }) => {
        await page.goto("/");

        await page.locator("#palette-toggle").click();

        const overlay = page.locator("#palette-overlay");
        await expect(overlay).toBeVisible();

        await expect(page.locator("#palette-results")).toContainText(
            "/portfolio",
        );
        await expect(page.locator("#palette-results")).toContainText("/blog");
    });

    test("closes on Escape", async ({ page }) => {
        await page.goto("/");

        await page.locator("#palette-toggle").click();
        await expect(page.locator("#palette-overlay")).toBeVisible();

        await page.keyboard.press("Escape");
        await expect(page.locator("#palette-overlay")).toBeHidden();
        await expect(page.locator("#palette-overlay")).toHaveAttribute(
            "inert",
            "",
        );
        await expect(page.locator("#palette-overlay")).toHaveAttribute(
            "aria-hidden",
            "true",
        );
    });
});
