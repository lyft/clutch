import React from "react";
import BaseCardComponent, { BaseCard } from "./base";
import type { ExtendedProjectCardProps } from "./base";

class DynamicCard extends BaseCardComponent {
  static displayName = "DynamicCard";

  constructor(props: ExtendedProjectCardProps) {
    super(props);
  }

  render() {
    return <BaseCard {...this.props} />;
  }
}

export default DynamicCard;
