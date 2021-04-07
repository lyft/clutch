import React from 'react';
import classnames from 'classnames';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

function Hero({ tagline, config }) {
  return (
    <header className={classnames('hero hero--primary', styles.heroSection)}>
      <div className={classnames("container", styles.container, styles.heroContainer)}>
        <h1 className="hero__title">{tagline}</h1>
        <h4 className={classnames(styles.heroDescription)}>{config.description}</h4>
        <div className={styles.buttons}>
          <Link
            className={classnames(
              'button button--outline button--lg',
              styles.button,
              styles.blueBtn,
              )}
            to={useBaseUrl(config.buttons.first.url)}>
            {config.buttons.first.text}
          </Link>
          <Link
            className={classnames(
              'button button--outline button--lg',
              styles.button,
              styles.greenBtn,
              )}
            to={useBaseUrl(config.buttons.second.url)}>
            {config.buttons.second.text}
          </Link>
        </div>
      </div>
      <div style={{width: "20%", margin: "0 5% 5% 5%"}}>
        <img src={useBaseUrl("img/microsite/home.svg")} alt="home icon" />
      </div>
    </header>
  );
};

function Feature({imageUrl, title, description}) {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={classnames('col col--4')}>
      {description === undefined ? (
        <div style={{display: "flex", justifyContent: "center"}} >
          <img src={imgUrl} alt={title} height="300" />
        </div>
      ) : (
        <>
          {imgUrl && (
            <div className={classnames("text--center", styles.featureIcon)}>
              <hr className={classnames(styles.featureAccent)} />
                <img className={styles.featureImage} src={imgUrl} alt={title} />
              <hr className={classnames(styles.featureAccent)} />
            </div>
          )}
          <div className="text--center"><h2>{title}</h2></div>
          <div className={classnames(styles.featureText)}>
            <p>{description}</p>
          </div>
        </>
      )}
    </div>
  );
};

function Features({ config }) {
  return (
    <section className={classnames(styles.section, styles.features)}>
      <div className={classnames("container", styles.container)}>
        <h3 className={classnames("hero__title", styles.sectionHeadingDark)}>{config.title}</h3>
        {config.featureList && config.featureList.length && (
          <div className="row">
            {config.featureList.map((props, idx) => (
              <Feature key={idx} {...props} />
              ))}
          </div>
        )}
      </div>
    </section>
  );
};

function Demo({ config }) {
  return (
    <div className={styles.darkBackground}>
      <section className={classnames("text--center", styles.section, styles.demoSection)}>
        <div className={classnames("container", styles.container, styles.demoContainer)}>
          {config.lines.map((line, idx) => (
            <h1 key={idx} className={classnames("hero__title", styles.sectionHeadingDark)}>{line}</h1>
          ))}
          <div style={{display: "flex", justifyContent: "center"}}>
            <Link
              className={classnames(
                'button button--outline button--lg',
                styles.button,
                styles.greenBtn,
                styles.demoBtn,
                )}
              to={useBaseUrl(config.cta.link)}>
                {config.cta.text}
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
};

function Consolidation({ config }) {
  return (
    <section className={classnames(styles.section, styles.consolidation, styles.darkBackground)}>
      <div className={classnames("container", styles.container)}>
        <img src={useBaseUrl('img/microsite/consolidation.gif')} />
      </div>
      <div className={classnames("container", styles.container, styles.textContainer)}>
        <div>
          {config.snippets.map((snippet, idx) => (
          <p key={idx}>{snippet}</p>  
          ))}
        </div>
      </div>
    </section>
  );
};

function Home() {
  const context = useDocusaurusContext();
  const {siteConfig = {}} = context;
  const sections = siteConfig.customFields.sections;
  return (
    <Layout title={siteConfig.title} description={siteConfig.customFields.tagDescription}>
      <Hero tagline={siteConfig.tagline} config={siteConfig.customFields.hero} />
      <main>
        <Features config={sections.features} />
        <Demo config={sections.demo} />
        <Consolidation config={sections.consolidation} />
      </main>
    </Layout>
  );
}

export default Home;
