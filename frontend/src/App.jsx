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
import AuthAPI from './api/AuthAPI';
import UsersAPI from './api/UsersAPI';
import BoardsAPI from './api/BoardsAPI';
import Home from './components/Home/Home';
import Login from './components/Login/Login';
import Register from './components/Register/Register';
import activeBoardInit from './misc/activeBoardInit';

import 'bootstrap/dist/css/bootstrap.min.css';
import './app.sass';
import TeamsAPI from './api/TeamsAPI';

const App = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [user, setUser] = useState({
    username: '',
    teamId: null,
    isAdmin: false,
    isAuthenticated: false,
  });
  const [team, setTeam] = useState({ id: null, inviteCode: '' });
  const [members, setMembers] = useState([{
    username: '',
    isActive: false,
    isAdmin: false,
  }]);
  const [boards, setBoards] = useState([{ id: null, name: '' }]);
  const [activeBoard, setActiveBoard] = useState(activeBoardInit);
  const notify = (header, body) => (header || body) && toast(
    <>
      {header && <h4>{header}</h4>}
      {body && <p>{body}</p>}
    </>,
  );

  const loadBoard = async (boardId) => {
    setIsLoading(true);
    try {
      // 1. set the current user
      const userResponse = await AuthAPI.verifyToken();
      delete userResponse.data.msg;
      setUser({ ...userResponse.data, isAuthenticated: true });

      // 2. set current team (only needed for inviting members, which means
      // only needed if the current user is the team admin)
      if (userResponse.data.isAdmin) {
        const teamResponse = await TeamsAPI.get(userResponse.data.teamId);
        setTeam(teamResponse.data);
      }

      // 3. set boards lists
      const teamBoards = await BoardsAPI.get(null, userResponse.data.teamId);
      setBoards(teamBoards.data);

      // 4. set the active board
      const nestedBoard = await BoardsAPI.get(
        boardId || activeBoard.id || teamBoards.data[0].id,
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

      // TODO: Check for cases where you want token verification failure
      //       message to render
      if (!err?.config?.url?.includes('verify-token')) {
        const errMsg = err?.response?.data?.msg;
        if (errMsg) { notify(errMsg); }
        if (err?.message) { notify(err.message); }
      }
    }
    setIsLoading(false);
  };

  useEffect(() => loadBoard(), []);

  return (
    <div className="App">
      <AppContext.Provider
        value={{
          user,
          team,
          members,
          boards,
          activeBoard,
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

      <ToastContainer />
    </div>
  );
};

export default App;
