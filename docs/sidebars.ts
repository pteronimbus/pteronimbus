import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  tutorialSidebar: [
    'intro',
    {
      type: 'category',
      label: 'Architecture',
      items: ['architecture/overview'],
    },
    {
      type: 'category', 
      label: 'Backend',
      items: ['backend/health-checks'],
    },
    {
      type: 'category',
      label: 'Frontend',
      items: ['frontend/health-checks'],
    },
    {
      type: 'category',
      label: 'Controller',
      items: ['controller/overview', 'controller/health-checks'],
    },
  ],
};

export default sidebars;
