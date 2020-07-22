const checker = require("license-checker"); // eslint-disable-line import/no-extraneous-dependencies

const allowLicenses = [
  "Apache-1.0",
  "Apache-1.1",
  "Apache-2.0",
  "Artistic-2.0",
  "BSD",
  "BSL-1.0",
  "bzip2-1.0.5",
  "bzip2-1.0.6",
  "CC-BY-1.0",
  "CC-BY-2.0",
  "CC-BY-2.5",
  "CC-BY-3.0",
  "CC-BY-4.0",
  "CC0-1.0",
  "curl",
  "ICU",
  "IJG",
  "ISC",
  "JSON",
  "libtiff",
  "Ms-PL",
  "MIT",
  "MPL-1.0",
  "MPL-1.1",
  "MPL-2.0",
  "OLDAP-1.1",
  "OLDAP-1.2",
  "OLDAP-1.3",
  "OLDAP-1.4",
  "OLDAP-2.0",
  "OLDAP-2.0.1",
  "OLDAP-2.1",
  "OLDAP-2.2",
  "OLDAP-2.2.1",
  "OLDAP-2.2.2",
  "OLDAP-2.3",
  "OLDAP-2.4",
  "OLDAP-2.5",
  "OLDAP-2.6",
  "OLDAP-2.7",
  "OLDAP-2.8",
  "OpenSSL",
  "PHP-3.0",
  "PHP-3.1",
  "PostgreSQL",
  "Public Domain",
  "Python-2.0",
  "Ruby",
  "TCL",
  "Unlicense",
  "W3C",
  "WTFPL",
  "Xnet",
  "X11",
  "libpng",
  "Zlib",
  "Zlib-acknowledgment",
  "ZPL",
];

const ignorePackages = ["@clutch-sh/clutch"];

const checkerArgs = {
  start: ".",
  exclude: allowLicenses.join(","),
  excludePackages: ignorePackages.join(";"),
};

checker.init(checkerArgs, (err, packages) => {
  if (err) {
    console.log(`Encountered exception: ${err}`); // eslint-disable-line no-console
    process.exit(1);
  }
  if (Object.keys(packages).length !== 0) {
    console.log("Error: Found node dependencies with unapproved licenses..."); // eslint-disable-line no-console
    console.log(checker.asTree(packages)); // eslint-disable-line no-console
    process.exit(1);
  }
  console.log("Success: All node dependencies have approved licenses"); // eslint-disable-line no-console
});
