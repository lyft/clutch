/* eslint-disable no-restricted-syntax */
module.exports = {
  async constraints({ Yarn }) {
    for (const dep of Yarn.dependencies({ ident: "react" })) {
      dep.update(`^17.0.2`);
    }

    for (const dep of Yarn.dependencies({ ident: "react-dom" })) {
      dep.update(`^17.0.2`);
    }

    for (const dep of Yarn.dependencies({ ident: "typescript" })) {
      dep.update(`^4.2.3`);
    }

    for (const workspace of Yarn.workspaces()) {
      workspace.set("engines.node", `<19`);
    }
  },
};
