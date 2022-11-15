import React from "react";
import { MDXProvider } from "@mdx-js/react";
import MDXComponents from "@theme/MDXComponents";

import "./styles.css";

interface RoadmapItemProps {
  title: string;
  description: string;
  children: React.ReactChildren;
}

const RoadmapItem = ({
  title,
  description,
  children,
}: RoadmapItemProps): JSX.Element => {
  const [displayLong, setDisplayLong] = React.useState(false);

  return (
    <div className="ri-container">
      <div className="ri-description-short">
        <div className="ri-icon">
          <span className="fe fe-zap" />
        </div>
        <div className="ri-detail">
          <div className="ri-title">{title}</div>
          <div className="ri-description">
            {description}
            {React.Children.count(children) > 0 && (
              <span
                className="ri-more fe fe-more-horizontal"
                onClick={() => setDisplayLong(!displayLong)}
              />
            )}
          </div>
        </div>
      </div>
      {displayLong && (
        <div className="ri-description-long">
          <MDXProvider components={MDXComponents}>{children}</MDXProvider>
        </div>
      )}
    </div>
  );
};

export default RoadmapItem;
