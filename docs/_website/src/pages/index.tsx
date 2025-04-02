import React from "react";
import classnames from "classnames";
import Layout from "@theme/Layout";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import useBaseUrl from "@docusaurus/useBaseUrl";
import styles from "./styles.module.css";
import {
  ConsolidationConfig,
  DemoConfig,
  FeatureConfig,
  FeaturesConfig,
  HeroConfig,
  SiteConfig,
} from "../theme/types";

function ArchivalNotice(): JSX.Element | null {
  const context = useDocusaurusContext();
  const siteConfig = context.siteConfig as SiteConfig;
  const notice = siteConfig.customFields.archivalNotice;

  if (notice === undefined || notice === null || !notice.enabled) {
    return null;
  }

  return (
    <div className={classnames("alert alert--warning", styles.archivalNotice)}>
      <div className="container">
        <h4 style={{ marginBottom: "0.5rem" }}>{notice.title}</h4>
        <p style={{ marginBottom: 0 }}>{notice.message}</p>
      </div>
    </div>
  );
}

interface HeroProps extends Pick<SiteConfig, "tagline"> {
  config: HeroConfig;
}
function Hero({ tagline, config }: HeroProps): JSX.Element {
  return (
    <header className={classnames("hero hero--primary", styles.heroSection)}>
      <div
        className={classnames(
          "container",
          styles.container,
          styles.heroContainer
        )}
      >
        <h1 className="hero__title">{tagline}</h1>
        <h4 className={classnames(styles.heroDescription)}>
          {config.description}
        </h4>
        <div className={styles.buttons}>
          <Link
            className={classnames(
              "button button--outline button--lg",
              styles.button,
              styles.blueBtn
            )}
            to={useBaseUrl(config.buttons.first.url)}
          >
            {config.buttons.first.text}
          </Link>
          <Link
            className={classnames(
              "button button--outline button--lg",
              styles.button,
              styles.greenBtn
            )}
            to={useBaseUrl(config.buttons.second.url)}
          >
            {config.buttons.second.text}
          </Link>
        </div>
      </div>
      <div style={{ width: "20%", margin: "0 5% 5% 5%" }}>
        <img src={useBaseUrl("img/microsite/home.svg")} alt="home icon" />
      </div>
    </header>
  );
}

function Feature({ imageUrl, title, description }: FeatureConfig): JSX.Element {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={classnames("col col--4")}>
      {description === undefined ? (
        <div style={{ display: "flex", justifyContent: "center" }}>
          <img src={imgUrl} alt={title} height="300" />
        </div>
      ) : (
        <>
          {imgUrl !== "" && (
            <div className={classnames("text--center", styles.featureIcon)}>
              <hr className={classnames(styles.featureAccent)} />
              <img className={styles.featureImage} src={imgUrl} alt={title} />
              <hr className={classnames(styles.featureAccent)} />
            </div>
          )}
          <div className="text--center">
            <h2>{title}</h2>
          </div>
          <div className={classnames(styles.featureText)}>
            <p>{description}</p>
          </div>
        </>
      )}
    </div>
  );
}

function Features({ title, featureList }: FeaturesConfig): JSX.Element {
  return (
    <section className={classnames(styles.section, styles.features)}>
      <div className={classnames("container", styles.container)}>
        <h3 className={classnames("hero__title", styles.sectionHeadingDark)}>
          {title}
        </h3>
        {featureList?.length !== 0 && (
          <div className="row">
            {featureList.map((props, idx) => (
              <Feature key={idx} {...props} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
}

function Demo({ lines, cta }: DemoConfig): JSX.Element {
  return (
    <div className={styles.darkBackground}>
      <section
        className={classnames(
          "text--center",
          styles.section,
          styles.demoSection
        )}
      >
        <div
          className={classnames(
            "container",
            styles.container,
            styles.demoContainer
          )}
        >
          {lines.map((line, idx) => (
            <h1
              key={idx}
              className={classnames("hero__title", styles.sectionHeadingDark)}
            >
              {line}
            </h1>
          ))}
          <div style={{ display: "flex", justifyContent: "center" }}>
            <Link
              className={classnames(
                "button button--outline button--lg",
                styles.button,
                styles.greenBtn,
                styles.demoBtn
              )}
              to={useBaseUrl(cta.link)}
            >
              {cta.text}
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}

function Consolidation({ snippets }: ConsolidationConfig): JSX.Element {
  return (
    <section
      className={classnames(
        styles.section,
        styles.consolidation,
        styles.darkBackground
      )}
    >
      <div className={classnames("container", styles.container)}>
        <img src={useBaseUrl("img/microsite/consolidation.gif")} />
      </div>
      <div
        className={classnames(
          "container",
          styles.container,
          styles.textContainer
        )}
      >
        <div>
          {snippets.map((snippet, idx) => (
            <p key={idx}>{snippet}</p>
          ))}
        </div>
      </div>
    </section>
  );
}

function Home(): JSX.Element {
  const context = useDocusaurusContext();
  const siteConfig = context.siteConfig as SiteConfig;
  const sections = siteConfig.customFields.sections;
  return (
    <Layout
      title={siteConfig.title}
      description={siteConfig.customFields.tagDescription}
    >
      <ArchivalNotice />
      <Hero
        tagline={siteConfig.tagline}
        config={siteConfig.customFields.hero}
      />
      <main>
        <Features
          title={sections.features.title}
          featureList={sections.features.featureList}
        />
        <Demo lines={sections.demo.lines} cta={sections.demo.cta} />
        <Consolidation snippets={sections.consolidation.snippets} />
      </main>
    </Layout>
  );
}

export default Home;
