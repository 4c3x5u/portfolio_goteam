import React, { useState } from 'react';

import Header from './Header/Header';
import ColumnsRow from './ColumnsRow/ColumnsRow';
import CreateBoard from './CreateBoard/CreateBoard';
import Footer from './Footer/Footer';
import { windowEnum } from './windowEnum';

import './home.sass';

const Home = () => {
  const [activeWindow, setActiveWindow] = useState(windowEnum.NONE);

  const handleActivate = (window) => () => (
    window === activeWindow
      ? setActiveWindow(windowEnum.NONE)
      : setActiveWindow(window)
  );

  return (
    <div id="Home">
      <Header
        activeWindow={activeWindow}
        handleActivate={handleActivate}
      />
      <ColumnsRow />

      {activeWindow === windowEnum.CREATE_BOARD
        && <CreateBoard deactivate={handleActivate} />}

      <Footer />
    </div>
  );
};

export default Home;
