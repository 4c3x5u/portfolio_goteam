/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React, { useContext, useState } from 'react';

import AppContext from '../../AppContext';
import Spinner from './Spinner/Spinner';
import Header from './Header/Header';
import Board from './Board/Board';
import CreateBoard from './CreateBoard/CreateBoard';
import DeleteBoard from './DeleteBoard/DeleteBoard';
import EditBoard from './EditBoard/EditBoard';
import CreateTask from './CreateTask/CreateTask';
import DeleteTask from './DeleteTask/DeleteTask';
import EditTask from './EditTask/EditTask';
import InviteMember from './Invite/Invite';
import Help from './Help/Help';
import Footer from './Footer/Footer';
import DeleteMember from './DeleteMember/DeleteMember';
import window from '../../misc/window';

import './home.sass';

const Home = () => {
  const { isLoading } = useContext(AppContext);
  const [activeWindow, setActiveWindow] = useState(window.NONE);
  const [windowState, setWindowState] = useState(null);

  const handleActivate = (newWindow) => (state) => {
    if (newWindow === activeWindow) {
      setActiveWindow(window.NONE);
    } else {
      if (state) { setWindowState(state); }
      setActiveWindow(newWindow);
    }
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

      case window.EDIT_BOARD:
        return (
          <EditBoard
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
            columnId={windowState.columnId}
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
            username={windowState.username}
            toggleOff={handleActivate(window.NONE)}
          />
        );

      case window.HELP:
        return <Help toggleOff={handleActivate(window.NONE)} />;

      default:
        return <></>;
    }
  };

  return (
    <div
      id="Home"
      onKeyDown={(e) => e.key === 'Escape' && setActiveWindow(window.NONE)}
    >
      {isLoading && <Spinner />}
      <Header activeWindow={activeWindow} handleActivate={handleActivate} />
      <Board handleActivate={handleActivate} />
      {viewActiveWindow()}
      <Footer />
    </div>
  );
};

export default Home;
