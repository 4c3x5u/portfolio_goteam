/* eslint-disable
no-nested-ternary,
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useState } from 'react';

import Header from './Header/Header';
import Board from './Board/Board';
import CreateBoard from './CreateBoard/CreateBoard';
import CreateTask from './CreateTask/CreateTask';
import InviteMember from './Invite/Invite';
import HelpModal from './Help/Help';
import Footer from './Footer/Footer';
import window from '../../misc/window';

import './home.sass';
import EditTask from './EditTask/EditTask';
import DeleteTask from './DeleteTask/DeleteTask';
import DeleteMember from './DeleteMember/DeleteMember';
import DeleteBoard from './DeleteBoard/DeleteBoard';

const Home = () => {
  const [activeWindow, setActiveWindow] = useState(window.NONE);
  const [windowState, setWindowState] = useState(null);

  const handleActivate = (newWindow) => (state) => {
    if (state) { setWindowState(state); }

    if (newWindow === activeWindow) {
      setActiveWindow(window.NONE);
    } else { setActiveWindow(newWindow); }
  };

  const viewActiveWindow = () => {
    switch (activeWindow) {
      case window.CREATE_BOARD:
        return <CreateBoard toggleOff={handleActivate(window.NONE)} />;

      case window.DELETE_BOARD:
        return (
          <DeleteBoard
            id={windowState.id}
            name={windowState.name}
            toggleOff={handleActivate(window.NONE)}
          />
        );

      case window.CREATE_TASK:
        return <CreateTask toggleOff={handleActivate(window.NONE)} />;

      case window.EDIT_TASK:
        return (
          <EditTask
            id={windowState.id}
            title={windowState.title}
            description={windowState.description}
            subtasks={windowState.subtasks}
            toggleOff={handleActivate(window.NONE)}
          />
        );

      case window.DELETE_TASK:
        return (
          <DeleteTask
            id={windowState.id}
            title={windowState.title}
            description={windowState.description}
            subtasks={windowState.subtasks}
            toggleOff={() => setActiveWindow(window.NONE)}
          />
        );

      case window.INVITE_MEMBER:
        return <InviteMember toggleOff={handleActivate(window.NONE)} />;

      case window.DELETE_MEMBER:
        return (
          <DeleteMember
            id={windowState.id}
            username={windowState.name}
            toggleOff={handleActivate(window.NONE)}
          />
        );

      case window.MODAL:
        return <HelpModal toggleOff={handleActivate(window.NONE)} />;

      default:
        return <></>;
    }
  };

  return (
    <div
      id="Home"
      onKeyDown={(e) => e.key === 'Escape' && setActiveWindow(window.NONE)}
    >
      <Header activeWindow={activeWindow} handleActivate={handleActivate} />

      <Board handleActivate={handleActivate} />

      {viewActiveWindow()}

      <Footer />
    </div>
  );
};

export default Home;
