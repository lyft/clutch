/* eslint-disable no-console */
const byteSize = require("byte-size");
const esbuild = require("esbuild");
const fs = require("fs");
const path = require("path");

const args = process.argv.slice(2);
const fileDetailsLimit = 20;

const getAllFiles = (dirPath, arrayOfFiles) => {
  const files = fs.readdirSync(dirPath);

  let tmpArrayOfFiles = arrayOfFiles || [];

  files.forEach(file => {
    if (fs.statSync(`${dirPath}/${file}`).isDirectory()) {
      if (file !== "tests" && file !== "stories" && file !== "dist" && file !== "__snapshots__") {
        tmpArrayOfFiles = getAllFiles(`${dirPath}/${file}`, tmpArrayOfFiles);
      }
    } else {
      tmpArrayOfFiles.push(path.join(dirPath, "/", file));
    }
  });

  return tmpArrayOfFiles;
};

const sizeOutputPlugin = {
  name: "sizeOutputPlugin",
  setup(build) {
    let timerStart;
    build.onStart(() => {
      timerStart = process.hrtime.bigint();
    });
    build.onEnd(result => {
      if (result?.metafile?.outputs !== undefined) {
        const files = Object.keys(result?.metafile?.outputs ?? {})
          .map(k => {
            return {
              name: k.substring(k.search("dist/"), k.length),
              size: result.metafile.outputs[k].bytes,
            };
          })
          .sort((a, b) => !b.name.includes(".map") - !a.name.includes(".map") || b.size - a.size);

        console.log("");
        files
          .slice(0, 20)
          .forEach(f => console.log(`\t${f.name.toString().padEnd(50)}${byteSize(f.size)}`));

        console.log(
          files.length > fileDetailsLimit
            ? `\t... and ${files.length - fileDetailsLimit} more output files...\n`
            : ""
        );
        const timerEnd = process.hrtime.bigint() - timerStart;
        // eslint-disable-next-line no-undef
        console.log(`\t⚡ Done in ${timerEnd / BigInt(1000000)}ms`);
      }
    });
  },
};

const options = {
  entryPoints: getAllFiles(`${process.argv[2]}/src`),
  outdir: `${process.argv[2]}/dist/`,
  target: "es2020",
  sourcemap: true,
  preserveSymlinks: true,
  color: true,
  plugins: [],
  tsconfig: `${process.argv[2]}/tsconfig.json`,
};

(async () => {
  if (args.includes("-w") || args.includes("--watch")) {
    const ctx = await esbuild.context({ ...options, logLevel: "info" });
    await ctx.watch();
  } else {
    await esbuild.build({
      ...options,
      metafile: true,
      plugins: [...options.plugins, sizeOutputPlugin],
    });
  }
})();
