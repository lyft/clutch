const { defineConfig } = require("@yarnpkg/types");

const workspaceRegex = /^((\.|api)|((packages|workflows)\/.*))$/;

const enforcedPackages = {
  "@emotion/react": "^11.0.0",
  "@emotion/styled": "^11.0.0",
  "@mui/icons-material": "^5.8.4",
  "@mui/lab": "^5.0.0-alpha.87",
  "@mui/material": "^5.8.5",
  "@mui/styles": "^5.8.4",
  "@mui/system": "^5.8.4",
  "@types/react-dom": "^17.0.3",
  "@types/react": "^17.0.5",
  "react-dom": "^17.0.2",
  "react-router-dom": "^6.0.0-beta.0",
  "react-router": "^6.0.0-beta.0",
  esbuild: "^0.18.0",
  eslint: "^8.3.0",
  jest: "^27.0.0",
  lodash: "^4.17.0",
  react: "^17.0.2",
  typescript: "^5.5.3",
};

const enforceWorkspaceEngines = workspace => {
  workspace.set("engines.node", ">=18 <19");
  workspace.set("engines.yarn", "^4.3.1");
};

const enforceDependencies = workspace => {
  workspace.pkg.dependencies.forEach(dep => {
    if (enforcedPackages[dep.ident]) {
      workspace.set(`dependencies.${dep.ident}`, enforcedPackages[dep.ident]);
    }
  });
};

const workspaceEnforcers = [enforceWorkspaceEngines];
const dependencyEnforcers = [enforceDependencies];

const constraintFn = workspace => {
  workspaceEnforcers.forEach(enforcer => enforcer(workspace));
  dependencyEnforcers.forEach(enforcer => enforcer(workspace));
};

module.exports = defineConfig({
  constraintFn,
  dependencyEnforcers,
  enforcedPackages,
  workspaceEnforcers,
  workspaceRegex,
  async constraints({ Yarn }) {
    Yarn.workspaces().forEach(workspace => {
      if (workspaceRegex.test(workspace.cwd.trim())) {
        constraintFn(workspace);
      }
    });
  },
});
