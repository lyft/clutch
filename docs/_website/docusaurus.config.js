module.exports = {
  title: 'Clutch · An extensible platform for infrastructure management.',
  tagline: 'Shifting infrastructure management to a friendlier place.',
  url: 'https://clutch.sh',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'lyft', // Usually your GitHub org/user name.
  projectName: 'clutch', // Usually your repo name.
  stylesheets: [
    'https://fonts.googleapis.com/css2?family=Roboto&family=Open+Sans&display=swap',
    'https://cdn.rawgit.com/luizbills/feather-icon-font/v4.7.0/dist/feather.css',
  ],
  plugins: [],
  customFields: {
    tagDescription: 'An extensible platform for infrastructure management.',
    hero: {
      description: "Clutch provides everything you need to improve your developers' experience and operational capabilities. It comes with several out-of-the-box features for managing cloud-native infrastructure, but is easily configured or extended to interact with whatever you run, wherever you run it.",
      buttons: {
        first: {
          url: "docs/about/what-is-clutch",
          text: "Learn More",
        },
        second: {
          url: "docs/getting-started/build-guides",
          text: "Get Started",
        },
      },
    },
    sections: {
      features: {
        title: "Why Clutch?",
        featureList: [
          {
            title: 'Secure',
            imageUrl: 'img/microsite/reasons/secure.svg',
            description: `
                Clutch has first-class support for role based access control down to the
                individual resource level. In addition, it ships with rich auditing so you
                can see what's happening in Slack, email notifications, or logs.
            `,
          },
          {
            title: '',
            imageUrl: 'img/microsite/logo.svg',
          },
          {
            title: 'Extensible. Really.',
            imageUrl: 'img/microsite/reasons/extensible.svg',
            description: `
                Highly configurable. No forks. Private extensions. Clutch's abstractions
                make it work for your environment without messy hacks or rewrites.
                Adding new features is easy too.
            `,
          },
          {
            title: 'One Entrypoint',
            imageUrl: 'img/microsite/reasons/single-entrypoint.svg',
            description: `
                Access your company's tech stack through a single pane of glass. But
                don't worry, it's not too fragile or breakable. Clutch is easy to maintain
                as your infrastructure evolves.
            `,
          },
          {
            title: 'Straightforward',
            imageUrl: 'img/microsite/reasons/user-experience.svg',
            description: `
                Consistent and clear design with built-in safegaurds throughout to
                turn your complex processes into simple and safe operations that anyone
                can understand.
            `,
          },
          {
            title: 'The Long Tail',
            imageUrl: 'img/microsite/reasons/file.svg',
            description: `
                Infrastructure as code is great, we love it too, but there's a lot
                of your infrastructure not covered by it.
            `,
          },
        ],
      },
      demo: {
        lines: [
          "Don't take our word for it.",
          "See what Clutch has to offer.",
        ],
        cta: {
          text: "Workflows & Components",
          link: "docs/components",
        },
      },
      consolidation: {
        snippets: [
          `
            Stop putting your team through an endless stream of high-friction tools and user interfaces.
            Clutch allows you to combine many tools into one, in the form that your developers use most.
          `,
          `
            We grow with you. Clutches extensible platform means you can integrate as many tools as
            you need, even if they are specific to you.
          `,
        ]
      }
    },
  },
  themeConfig: {
    colorMode: {
      disableSwitch: true,
    },
    googleAnalytics: {
      trackingID: 'UA-170615678-4',
      anonymizeIP: true,
    },
    algolia : {
      apiKey: '32f1f7956b3d2c3c90fbe259c7901d94',
      indexName: 'lyft_clutch',
    },
    prism : {
      additionalLanguages: ['protobuf', 'typescript'],
      theme: require('prism-react-renderer/themes/vsDark'),
    },
    navbar: {
      title: 'Clutch',
      logo: {
        alt: 'Clutch Logo',
        src: 'img/navigation/logoMark.svg',
      },
      items: [], // items are defined directly in the swizzled component so they can have an icon attr.
    },
    footer: {
      style: 'light',
      logo: {
        src: "img/navigation/logo.svg"
      },
      links: [
        {
          title: 'About',
          items: [
            {
              label: 'What is Clutch?',
              to: 'docs/about/what-is-clutch',
            },
            {
              label: 'Roadmap',
              to: 'docs/about/roadmap',
            },
            {
              label: 'Architecture',
              to: 'docs/about/architecture',
            },
          ],
        },
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: 'docs/getting-started/build-guides',
            },
            {
              label: 'Development',
              to: 'docs/development/guide',
            },
            {
              label: 'Configuration',
              to: 'docs/configuration',
            },
          ],
        },
        {
          title: 'Components',
          items: [
            {
              label: 'Frontend',
              to: 'docs/components#frontend',
            },
            {
              label: 'Backend',
              to: 'docs/components#backend',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub',
              to: 'https://github.com/lyft/clutch',
            },
            {
              label: 'Slack',
              to: 'https://join.slack.com/t/lyftoss/shared_invite/zt-casz6lz4-G7gOx1OhHfeMsZKFe1emSA',
            },
            {
              label: 'Twitter',
              to: 'https://twitter.com/clutchdotsh',
            },
            {
              label: 'More',
              to: 'docs/community',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} <a href="https://lyft.com" target="blank">Lyft, Inc.</a>`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          path: "generated/docs",
          sidebarPath: require.resolve('../sidebars.json'),
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
