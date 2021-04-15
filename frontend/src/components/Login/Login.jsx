import React, { useContext, useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import axios from 'axios';

import AppContext from '../../AppContext';
import FormGroup from '../_shared/FormGroup/FormGroup';
import inputType from '../../misc/inputType';

import logo from './login.svg';
import './login.sass';

const Login = () => {
  const { loadBoard } = useContext(AppContext);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    axios.post(`${process.env.REACT_APP_BACKEND_URL}/login/`, {
      username,
      password,
    }).then((res) => {
      sessionStorage.setItem('username', res.data.username);
      sessionStorage.setItem('auth-token', res.data.token);
      loadBoard();
    }).catch((err) => {
      // TODO: Add toastr for server-side errors
      console.error(`SERVER-SIDE ERROR: ${JSON.stringify(err)}`);
    });
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
        />

        <FormGroup
          type={inputType.PASSWORD}
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
