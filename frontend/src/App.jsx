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
import Spinner from './components/Home/Spinner/Spinner';

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

  const loadBoard = () => {
    const authHeaders = getAuthHeaders();
    if (authHeaders.headers['auth-user'] && authHeaders.headers['auth-token']) {
      axios
        .get(
          `${process.env.REACT_APP_BACKEND_URL}/client-state/?boardId=${
            sessionStorage.getItem('board-id') || activeBoard.id || ''
          }`, authHeaders,
        )
        .then((res) => {
          // Update app state one by one
          setUser(res.data.user);
          if (res.data?.team) { setTeam(res.data.team); }
          setBoards(res.data.boards);
          setActiveBoard(res.data.activeBoard);
          setMembers(res.data.members);
        })
        .catch((err) => {
          setUser({ ...user, isAuthenticated: false });

          // remove username and auth token if verify-token failed
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
        })
        .finally(() => setIsLoading(false));
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
          {isLoading && <Spinner />}
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
