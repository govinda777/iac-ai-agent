import os
from playwright.sync_api import sync_playwright

def run(playwright):
    browser = playwright.chromium.launch()
    page = browser.new_page()
    # Use os.path.abspath to get the full path to the file
    file_path = os.path.abspath('index.html')
    page.goto(f'file://{file_path}')
    page.screenshot(path='jules-scratch/verification/verification.png')
    browser.close()

with sync_playwright() as playwright:
    run(playwright)