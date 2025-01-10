import type { AppBanners } from "./notification";

export interface AppConfiguration {
  /** Will override the title of the given application */
  title?: string;
  /** Supports a react node or a string representing a public assets path */
  logo?: React.ReactNode | string;
  banners?: AppBanners;
  useWorkflowLayout?: boolean;
  useFullScreenLayout?: boolean;
}
