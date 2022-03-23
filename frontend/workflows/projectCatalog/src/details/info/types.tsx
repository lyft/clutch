import type { CHIP_VARIANTS } from "@clutch-sh/core";
import type { IconProp } from "@fortawesome/fontawesome-svg-core";

export interface ProjectInfoChip {
  text: string;
  title?: string;
  icon?: any;
  url?: string;
  variant?: typeof CHIP_VARIANTS[number];
}

export interface ProjectMessenger {
  text: string;
  icon?: IconProp;
  url?: string;
}

interface ProjectRequests {
  number: number;
  url?: string;
  type: string;
}

export interface ProjectRepository {
  name: string;
  repo: string;
  url?: string;
  icon?: IconProp;
  requests?: ProjectRequests;
}

export interface ProjectInfo {
  owner: string;
  name: string;
  disabled?: boolean;
  description?: string;
  repository: ProjectRepository;
  languages?: string[];
  messenger?: ProjectMessenger;
  chips?: ProjectInfoChip[];
  serviceIds?: string[];
}
