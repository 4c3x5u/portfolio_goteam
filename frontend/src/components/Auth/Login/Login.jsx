import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import FormGroup from '../_shared/FormGroup';
import logo from '../../../assets/loginHeader.svg';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const handleSubmit = () => 1; // TODO: implement
  return (
    <div id="Login">
      <Form className="Form" onSubmit={handleSubmit}>
        <div className="HeaderWrapper">
          <img className="Header" alt="logo" src={logo} />
        </div>

        <FormGroup label="username" value={username} setValue={setUsername} />
        <FormGroup label="password" value={password} setValue={setPassword} />

        <div className="ButtonWrapper">
          <Button className="Button" type="submit" aria-label="submit">
            GO!
          </Button>
        </div>

        <div className="Register">
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
