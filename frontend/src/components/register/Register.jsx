import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import './register.sass';
import logo from '../../assets/register_title.svg';

const Register = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');
  const handleSubmit = () => 1; // TODO: implement
  return (
    <div id="Register">
      <Form className="Form" onSubmit={handleSubmit}>
        <div className="TitleWrapper">
          <img className="Title" alt="logo" src={logo} />
        </div>

        <Form.Group className="Group">
          <Form.Label className="Label">
            USERNAME
          </Form.Label>
          <Form.Control
            className="Input"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </Form.Group>

        <Form.Group className="Group">
          <Form.Label className="Label">
            PASSWORD
          </Form.Label>
          <Form.Control
            className="Input"
            type="password"
            value={passwordConfirmation}
            onChange={(e) => setPasswordConfirmation(e.target.value)}
          />
        </Form.Group>

        <Form.Group className="Group">
          <Form.Label className="Label">
            PASSWORD CONFIRMATION
          </Form.Label>
          <Form.Control
            className="Input"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Group>

        <div className="ButtonWrapper">
          <Button className="Button" value="GO!" type="submit">
            GO!
          </Button>
        </div>

        <div className="Login">
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
