import React, { useContext, useState } from 'react';
import { Link } from 'react-router-dom';
import { Form, Button } from 'react-bootstrap';

import AppContext from '../../AppContext';
import AuthAPI from '../../api/AuthAPI';
import FormGroup from '../_shared/FormGroup/FormGroup';
import validateRegisterForm from './validateRegisterForm';
import inputType from '../../misc/inputType';

import logo from './register.svg';
import './register.sass';

const Register = () => {
  const { loadBoard } = useContext(AppContext);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');
  const [errors, setErrors] = useState({
    username: '',
    password: '',
    passwordConfirmation: '',
  });

  const handleSubmit = (e) => {
    e.preventDefault();

    setErrors(validateRegisterForm(
      username,
      password,
      passwordConfirmation,
    ));

    if (!errors.username && !errors.password && !errors.passwordConfirmation) {
      AuthAPI
        .register(username, password, passwordConfirmation)
        .then((res) => {
          sessionStorage.setItem('username', res.data.username);
          sessionStorage.setItem('auth-token', res.data.token);
          loadBoard();
        })
        .catch((err) => console.error(err)); // TODO: Toastr
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
