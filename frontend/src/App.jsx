import React, { useState, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.min.css';
import axios from 'axios';

import AppContext from './AppContext';
import InitialStates from './misc/InitialStates';
import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';
import { getAuthHeaders } from './misc/util';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';

const App = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [user, setUser] = useState(InitialStates.user);
  const [team, setTeam] = useState(InitialStates.team);
  const [members, setMembers] = useState(InitialStates.members);
  const [boards, setBoards] = useState(InitialStates.boards);
  const [activeBoard, setActiveBoard] = useState(InitialStates.activeBoard);
  const notify = (header, body) => (header || body) && toast.error(
    <>
      {header && <h4>{header}</h4>}
      {body && <p>{body}</p>}
    </>,
  );

  const loadBoard = async () => {
    if (
      sessionStorage.getItem('username')
        && sessionStorage.getItem('auth-token')
    ) {
      try {
        const clientState = await axios.get(
          `${process.env.REACT_APP_BACKEND_URL}/client-state/?boardId=${
            sessionStorage.getItem('board-id') || activeBoard.id || ''
          }`, getAuthHeaders(),
        );

        // Set the current user
        setUser(clientState.data.user);

        // Set current team (for invites) if the user is admin
        if (clientState.data?.team) { setTeam(clientState.data.team); }

        // Set boards lists
        try {
          setBoards(clientState.data.boards);
        } catch (e) {
          notify(
            'Unable to load board.',
            'Please ask your team admin to add you to a board.',
          );
        }

        // Set the active board
        setActiveBoard(clientState.data.activeBoard);

        // Set team members
        setMembers(clientState.data.members);
      } catch (err) {
        setUser({ ...user, isAuthenticated: false });

        if (err?.config?.url?.includes('verify-token')) {
          sessionStorage.removeItem('username');
          sessionStorage.removeItem('auth-token');
          setIsLoading(false);
          return;
        }

        notify(
          'Unable to load board.',
          `${err?.message || 'Server Error'}.`,
        );
      }
      setIsLoading(false);
    }
  };

  useEffect(() => loadBoard(), []);

  return (
    <div className="App">
      <AppContext.Provider
        value={{
          user,
          setUser,
          team,
          setTeam,
          members,
          setMembers,
          boards,
          setBoards,
          activeBoard,
          setActiveBoard,
          loadBoard,
          isLoading,
          setIsLoading,
          notify,
        }}
      >
        <Router>
          <Switch>
            <Route exact path="/">
              {user.isAuthenticated
                ? <Home />
                : <Redirect to="/login" />}
            </Route>

            <Route path="/login">
              {!user.isAuthenticated
                ? <Login />
                : <Redirect to="/" />}
            </Route>

            <Route path="/register/:inviteCode?">
              {!user.isAuthenticated
                ? <Register />
                : <Redirect to="/" />}
            </Route>
          </Switch>
        </Router>
      </AppContext.Provider>

      <ToastContainer
        toastClassName="ErrorToast"
        position="bottom-left"
        autoClose={false}
      />
    </div>
  );
};

export default App;
