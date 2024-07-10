import { createLogger, defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import viteTsconfigPaths from "vite-tsconfig-paths";

const customLogger = createLogger();
const customLoggerWarn = customLogger.warn;

customLogger.warn = (message, ...args) => {
  if (message.includes("Is the variable mistyped?")) {
    return;
  }
  customLoggerWarn(message, ...args);
};

export default defineConfig({
  // depending on your application, base can also be "/"
  base: "",
  plugins: [react(), viteTsconfigPaths()],
  customLogger,
  server: {
    // this ensures that the browser opens upon server start
    open: true,
    // this sets a default port to 3000
    port: 3000,
    strictPort: true,
  },
});
