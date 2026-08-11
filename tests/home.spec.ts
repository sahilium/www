import { expect, test } from "@playwright/test";

test.describe("home page", () => {
    test("loads with the site title and greeting", async ({ page }) => {
        await page.goto("/");

        await expect(page).toHaveTitle("Sahil A.");
        await expect(page.locator("h1").first()).toBeVisible();
        await expect(
            page.getByText("Latest writing", { exact: true }),
        ).toBeVisible();
    });

    test("has working navigation links in the navbar", async ({ page }) => {
        await page.goto("/");

        const aboutLink = page.locator('#main-nav a[href="/about"]');
        await expect(aboutLink).toBeVisible();
        await aboutLink.click();

        await expect(page).toHaveURL(/\/about$/);
        await expect(page.getByRole("heading", { name: "Sahil A." })).toBeVisible();
    });

    test("footer contains copyright and motto", async ({ page }) => {
        await page.goto("/");

        const footer = page.locator("footer");
        await expect(footer).toContainText("Sahil A.");
        await expect(footer).toContainText("here be dragons");
        await expect(footer).toContainText("CREDITS");
    });
});
