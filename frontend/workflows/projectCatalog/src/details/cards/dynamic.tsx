import React from "react";

import type { DetailsCardTypes } from "../..";

import BaseCardComponent, { BaseCard } from "./base";

class DynamicCard extends BaseCardComponent {
  render() {
    return (
      <BaseCard {...this.props} loading={this.state.loading} type={"Dynamic" as DetailsCardTypes} />
    );
  }
}

export default DynamicCard;
