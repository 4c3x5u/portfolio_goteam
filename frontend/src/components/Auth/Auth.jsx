import React from 'react';
import { Switch, Route } from 'react-router-dom';
import Register from './subcomponents/Register';
import Login from './subcomponents/Login';
import './auth.sass';

const Auth = () => (
  <Switch>
    <Route exact path="/login" component={Login} />
    <Route exact path="/register" component={Register} />
  </Switch>
);

export default Auth;
