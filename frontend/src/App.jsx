import React, { useState, useEffect } from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
} from 'react-router-dom';
import _ from 'lodash/core';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.min.css';

import AppContext from './AppContext';
import InitialStates from './misc/InitialStates';
import AuthAPI from './api/AuthAPI';
import UsersAPI from './api/UsersAPI';
import BoardsAPI from './api/BoardsAPI';
import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';
import TeamsAPI from './api/TeamsAPI';

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
      setIsLoading(true);
      try {
        // 1. set the current user
        const userResponse = await AuthAPI.verifyToken();
        delete userResponse.data.msg;
        setUser({ ...userResponse.data, isAuthenticated: true });

        // 2. set current team (only needed for inviting members, which means
        //    only needed if the current user is the team admin)
        if (userResponse.data.isAdmin) {
          const teamResponse = await TeamsAPI.get(userResponse.data.teamId);
          setTeam(teamResponse.data);
        }

        // 3. set boards lists
        let teamBoards = [];
        try {
          teamBoards = await BoardsAPI.get(null, userResponse.data.teamId);
        } catch (e) {
          notify(
            'Unable to load board.',
            'Please ask your team admin to add you to a board.',
          );
        }

        setBoards(teamBoards?.data || []);

        // 4. set the active board
        const nestedBoard = await BoardsAPI.get(
          sessionStorage.getItem('board-id')
          || activeBoard.id
          || (teamBoards?.data && teamBoards?.data[0]?.id),
        );
        setActiveBoard(nestedBoard.data);

        // 5. set members list
        const teamMembers = await UsersAPI.get(
          userResponse.data.teamId,
          nestedBoard.data.id,
        );

        setMembers(
          _.sortBy(teamMembers.data, (member) => !member.isAdmin),
        );
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
