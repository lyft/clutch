const chokidar = require('chokidar');
const { exec } = require('child_process');

const watcher = chokidar.watch('../', {
  ignored: /_website/,
  persistent: true
});

const generateDocs = () => {
  exec("go run . -progressive", {cwd: "generator"});
};

const log = console.log.bind(console);
watcher
  .on('add', path => {console.log(`Detected new file ${path}]`); generateDocs();})
  .on('change', path => {console.log(`Detected change in ${path}]`); generateDocs();})
  .on('unlink', path => {console.log(`Detected removal of ${path}]`); generateDocs();});

watcher.once('ready', () => {
  console.error('Watching', `"${JSON.stringify(watcher.getWatched())}" ..`);
});