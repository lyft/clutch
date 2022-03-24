import React from "react";

import type { DetailsCardTypes } from "../..";

import BaseCardComponent, { BaseCard } from "./base";

class MetaCard extends BaseCardComponent {
  render() {
    return (
      <BaseCard
        {...this.props}
        error={this.state.error}
        loading={this.state.loading}
        type={"Metadata" as DetailsCardTypes}
      />
    );
  }
}

export default MetaCard;
