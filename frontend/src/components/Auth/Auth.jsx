import React from 'react';
import { Switch, Route } from 'react-router-dom';
import Register from './Register/Register';
import Login from './Login/Login';
import './auth.sass';

const Auth = () => (
  <Switch>
    <Route exact path="/login" component={Login} />
    <Route exact path="/register" component={Register} />
  </Switch>
);

export default Auth;
