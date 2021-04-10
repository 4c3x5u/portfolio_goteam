import React, { useState, useEffect } from 'react';
import { Redirect } from 'react-router-dom';
import { Form, Button } from 'react-bootstrap';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import verifyToken from '../../../misc/verifyToken';
import inputType from '../../../misc/inputType';

import logo from './login.svg';
import './login.sass';

const Login = () => {
  const [authenticated, setAuthenticated] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    verifyToken()
      .then(() => {
        setAuthenticated(true);
        window.location.reload();
      })
      .catch(() => {
        setAuthenticated(false);
      });
  }, []);

  // TODO: implement
  const handleSubmit = (e) => {
    e.preventDefault();
    window.location.reload();
  };

  if (authenticated) { return <Redirect to="/" />; }

  return (
    <div id="Login">
      <Form className="Form" onSubmit={handleSubmit}>
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup
          type={inputType.TEXT}
          label="username"
          value={username}
          setValue={setUsername}
        />

        <FormGroup
          type={inputType.TEXT}
          label="password"
          value={password}
          setValue={setPassword}
        />

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>

        <div className="Redirect">
          <p>Don&apos;t have an account yet?</p>
          <p>
            <a href="/register">Register now.</a>
          </p>
        </div>
      </Form>
    </div>
  );
};

export default Login;
