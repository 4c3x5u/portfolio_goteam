import React, { useContext, useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import { toast } from 'react-toastify';

import AppContext from '../../AppContext';
import AuthAPI from '../../api/AuthAPI';
import FormGroup from '../_shared/FormGroup/FormGroup';
import ValidateUser from '../../validation/ValidateUser';
import inputType from '../../misc/inputType';

import logo from './login.svg';
import './login.sass';

const Login = () => {
  const { loadBoard } = useContext(AppContext);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState({ username: '', password: '' });

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientErrors = {
      username: ValidateUser.username(username),
      password: ValidateUser.password(password),
    };

    AuthAPI
      .login(username, password)
      .then((res) => {
        sessionStorage.setItem('username', res.data.username);
        sessionStorage.setItem('auth-token', res.data.token);
        loadBoard();
      })
      .catch((err) => (
        err?.message
          ? toast(
            <>
              <h5>{err.message}</h5>
              <p>The board has NOT been deleted.</p>
            </>,
          )
          : setErrors({
            username: clientErrors.username || err.response.data.username || '',
            password: clientErrors.password || err.response.data.password || '',
          })
      ));
  };

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
          error={errors.username}
        />

        <FormGroup
          type={inputType.PASSWORD}
          label="password"
          value={password}
          setValue={setPassword}
          error={errors.password}
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
