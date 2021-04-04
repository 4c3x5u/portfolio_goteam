import React from 'react';

import Header from './subcomponents/Header';
import ColumnsRow from './subcomponents/ColumnsRow';
import Footer from './subcomponents/Footer';

import './home.sass';

const Home = () => (
  <div id="Home">
    <Header />
    <ColumnsRow />
    <Footer />
  </div>
);

export default Home;
