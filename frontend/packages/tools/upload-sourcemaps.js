/* eslint-disable no-console */
const { browser } = require("@bugsnag/source-maps");
const { config } = require("dotenv");

const srcDir = process.argv[2];
const buildDir = process.argv[3];
const envFile = process.argv[4];

const dotEnvFile = `${srcDir}/${envFile}`;

config({ path: dotEnvFile });

const uploadToBugsnag = async ({ apiKey, baseUrl, distDir: directory }) => {
  try {
    await browser.uploadMultiple({
      apiKey,
      baseUrl,
      directory,
      overwrite: true,
    });
    console.info(`[BugSnag] Successfully uploaded source maps ${directory} to BugSnag`);
  } catch (err) {
    console.error(`[BugSnag] Error uploading source maps to BugSnag: ${err}`);
  }
};

const uploadBugsnagSourcemaps = () => {
  const apiKey = process.env.REACT_APP_BUGSNAG_API_TOKEN || "";
  const baseUrl = process.env.REACT_APP_BASE_URL || "";
  if (!apiKey) {
    console.error(`[BugSnag] No API token found in ${dotEnvFile} file. Skipping upload.`);
    return Promise.reject(new Error("[BugSnag] API Key missing"));
  }

  if (!baseUrl) {
    console.error(
      `[BugSnag] No BaseUrl defined in process.env.BASE_URL in ${dotEnvFile}. Skipping Upload`
    );
    return Promise.reject(new Error("[BugSnag] BaseUrl missing"));
  }

  return uploadToBugsnag({
    apiKey,
    baseUrl: `${baseUrl}${process.env.REACT_APP_BASE_URL_PATH ?? "/static/js/"}`,
    distDir: process.env.SOURCEMAPS_DIR || `${srcDir}/${buildDir}`,
  });
};

uploadBugsnagSourcemaps();
