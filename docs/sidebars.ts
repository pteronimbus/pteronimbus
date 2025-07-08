import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

/**
 * Pteronimbus documentation structure.
 * Only includes components that have been implemented.
 */
const sidebars: SidebarsConfig = {
  // Architecture - system design and concepts
  architectureSidebar: [
    {
      type: 'category',
      label: 'Architecture',
      items: [
        'architecture/overview',
      ],
    },
  ],

  // Backend - what we've actually built
  backendSidebar: [
    {
      type: 'category',
      label: 'Backend Service',
      items: [
        'backend/health-checks',
      ],
    },
  ],
};

export default sidebars;
