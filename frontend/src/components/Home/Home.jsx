import React from 'react';

import AppHeader from './subcomponents/AppHeader';
import ColumnsRow from './subcomponents/ColumnsRow';
import AppFooter from './subcomponents/AppFooter';

import './home.sass';

const Home = () => (
  <div id="Home">
    <AppHeader />
    <ColumnsRow />
    <AppFooter />
  </div>
);

export default Home;
