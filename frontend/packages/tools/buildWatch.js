import childProcess from "child_process";
import path from "path";
import chalk from "chalk";

const { WORKSPACES } = process.env;
const buildCmd = process.argv[2];
const startCmd = process.argv[3];
const PROJECT_CWD = process.argv[4];

if (!WORKSPACES && !WORKSPACES.length) {
  throw new Error("WORKSPACES environment variable is required");
}
if (!buildCmd && !buildCmd.length) {
  throw new Error("buildCmd argument is required - example 'compile:watch'");
}
if (!startCmd && !startCmd.length) {
  throw new Error("startCmd argument is required - example 'start'");
}
if (!PROJECT_CWD && !PROJECT_CWD.length) {
  throw new Error("PROJECT_CWD argument is required - example '$INIT_CWD'");
}

/* eslint-disable no-console */
const log = (msg, color) => {
  if (color) {
    console.log(chalk.hex(color).bold(msg));
  } else {
    console.log(msg);
  }
};
/* eslint-enable no-console */

const storedColors = new Set();
const genRandom = () => {
  const randomColor = Math.floor(Math.random() * 16777215).toString(16);
  if (storedColors.has(randomColor)) {
    return genRandom();
  }
  return `#${randomColor}`;
};

const parsedWorkspaces = WORKSPACES.split("\n").map(workspace => ({
  ...JSON.parse(workspace),
  color: genRandom(),
}));

Promise.all(
  parsedWorkspaces
    .filter(({ location }) => location !== ".")
    .map(workspace => {
      return new Promise(resolve => {
        const child = childProcess.spawn("yarn", ["run", buildCmd], {
          cwd: path.join(PROJECT_CWD, workspace.location),
        });

        child.stdout.on("data", data => {
          if (!data.includes(`Couldn't find a script named "${buildCmd}".`)) {
            log(`[${workspace.name}]: ${data.toString().trim()}`, workspace.color);
          }
          resolve();
        });

        child.stderr.on("data", data => {
          log(`[${workspace.name}]: ${data.toString().trim()}`, workspace.color);
          resolve();
        });
      });
    })
).then(() => {
  log("\nStarting Server...\n");

  const child = childProcess.spawn("yarn", ["run", startCmd], { cwd: PROJECT_CWD });

  child.stdout.on("data", data => {
    log(`${data.toString().trim()}\n`);
  });

  child.stderr.on("data", data => {
    log(`${data.toString().trim()}\n`);
  });
});
