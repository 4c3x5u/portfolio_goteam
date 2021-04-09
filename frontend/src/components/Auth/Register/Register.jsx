import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import axios from 'axios';

import FormGroup from '../../_shared/FormGroup/FormGroup';
import { inputType } from '../../../misc/inputType';

import logo from './register.svg';
import './register.sass';

const Register = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    axios.post(`${process.env.REACT_APP_BACKEND_URL}/register/`, {
      username,
      password,
      password_confirmation: passwordConfirmation,
    }).then(() => (
      console.log('Register Successful')
    )).catch((err) => (
      console.log(`ERROR=${err.message}`)
    ));
  };

  return (
    <div id="Register">
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

        <FormGroup
          type={inputType.PASSWORD}
          label="password confirmation"
          value={passwordConfirmation}
          setValue={setPasswordConfirmation}
        />

        <div className="ButtonWrapper">
          <Button className="Button" value="GO!" type="submit">
            GO!
          </Button>
        </div>

        <div className="Redirect">
          <p>Already have an account?</p>
          <p>
            <a href="/login">Login here.</a>
          </p>
        </div>
      </Form>
    </div>
  );
};

export default Register;
