import esbuild from "esbuild";
import fs from "fs";
import path from "path";

const args = process.argv.slice(2);

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

const options = {
  entryPoints: getAllFiles(`${process.argv[2]}/src`),
  outdir: `${process.argv[2]}/dist/`,
  target: "es2019",
  sourcemap: true,
  tsconfig: `${process.argv[2]}/tsconfig.json`,
};

if (args.includes("-w") || args.includes("--watch")) {
  const ctx = await esbuild.context(options);
  await ctx.watch();
} else {
  await esbuild.build(options);
}
