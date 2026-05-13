import { test, expect } from '@playwright/test';

test('has title or main view', async ({ page }) => {
  // Wait for the app to load
  await page.goto('/');

  // As a Wails app, it might load as a visualizer. Let's check for the body or main div.
  await expect(page.locator('#app')).toBeVisible();
  
  // Let's assert there's some content, even if it's an upload button.
  // We'll just verify the page didn't crash.
  const appLoc = page.locator('#app');
  await expect(appLoc).toBeVisible();
});
