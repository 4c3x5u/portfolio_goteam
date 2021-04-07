/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useState } from 'react';

import Header from './Header/Header';
import Board from './Board/Board';
import CreateBoard from './CreateBoard/CreateBoard';
import CreateTask from './CreateTask/CreateTask';
import InviteMember from './InviteMember/InviteMember';
import HelpModal from './HelpModal/HelpModal';
import Footer from './Footer/Footer';
import { window } from '../../misc/enums';

import './home.sass';

const Home = () => {
  const [activeWindow, setActiveWindow] = useState(window.NONE);

  const handleActivate = (newWindow) => () => (
    newWindow === activeWindow
      ? setActiveWindow(window.NONE)
      : setActiveWindow(newWindow)
  );

  return (
    <div
      id="Home"
      onKeyDown={(e) => (
        e.key === 'Escape'
          && setActiveWindow(window.NONE)
      )}
    >
      <Header activeWindow={activeWindow} handleActivate={handleActivate} />

      <Board toggleCreateTask={handleActivate(window.CREATE_TASK)} />

      {activeWindow === window.CREATE_BOARD
        && <CreateBoard toggleOff={handleActivate(window.NONE)} />}

      {activeWindow === window.CREATE_TASK
        && <CreateTask toggleOff={handleActivate(window.NONE)} />}

      {activeWindow === window.INVITE_MEMBER
        && <InviteMember toggleOff={handleActivate(window.NONE)} />}

      {activeWindow === window.MODAL
        && <HelpModal toggleOff={handleActivate(window.NONE)} />}

      <Footer />
    </div>
  );
};

export default Home;
