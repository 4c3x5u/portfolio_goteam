import React, { useContext, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Form, Button } from 'react-bootstrap';

import AppContext from '../../AppContext';
import AuthAPI from '../../api/AuthAPI';
import FormGroup from '../_shared/FormGroup/FormGroup';
import ValidateUser from '../../validation/ValidateUser';
import inputType from '../../misc/inputType';

import logo from './register.svg';
import './register.sass';

const Register = () => {
  const { loadBoard, setIsLoading, notify } = useContext(AppContext);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');
  const [errors, setErrors] = useState({
    username: '',
    password: '',
    passwordConfirmation: '',
  });
  const { inviteCode } = useParams();

  const handleSubmit = (e) => {
    e.preventDefault();

    const clientErrors = {
      username:
        ValidateUser.username(username),
      password:
        ValidateUser.password(password),
      passwordConfirmation:
        ValidateUser.passwordConfirmation(passwordConfirmation),
    };

    if (
      clientErrors.username
        || clientErrors.password
        || clientErrors.passwordConfirmation
    ) {
      setErrors(clientErrors);
    } else {
      AuthAPI
        .register(username, password, passwordConfirmation, inviteCode)
        .then((res) => {
          sessionStorage.setItem('username', res.data.username);
          sessionStorage.setItem('auth-token', res.data.token);
          setIsLoading(true);
          loadBoard();
        })
        .catch((err) => {
          const serverErrors = {
            username:
              err?.response?.data?.username || '',
            password:
              err?.response?.data?.password || '',
            passwordConfirmation:
              // eslint-disable-next-line camelcase
              err?.response?.data?.password_confirmation || '',
          };

          if (
            serverErrors.username
              || serverErrors.password
              || serverErrors.passwordConfirmation
          ) {
            setErrors(serverErrors);
          } else {
            notify(
              'Unable to register.',
              `${err.message || 'Server Error'}.`,
            );
          }
        });
    }
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
          error={errors.username}
        />

        <FormGroup
          type={inputType.PASSWORD}
          label="password"
          value={password}
          setValue={setPassword}
          error={errors.password}
        />

        <FormGroup
          type={inputType.PASSWORD}
          label="password confirmation"
          value={passwordConfirmation}
          setValue={setPasswordConfirmation}
          error={errors.passwordConfirmation}
        />

        <div className="ButtonWrapper">
          <Button className="Button" value="GO!" type="submit">
            GO!
          </Button>
        </div>

        <div className="Redirect">
          <p>Already have an account?</p>
          <p>
            <Link to="/login">Login here.</Link>
          </p>
        </div>
      </Form>
    </div>
  );
};

export default Register;
