import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import './register.sass';
import logo from '../../assets/register_title.svg';

const Register = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const handleSubmit = () => 1; // TODO: implement
  return (
    <div id="FormWrapper">
      <Form id="RegisterForm" onSubmit={handleSubmit}>
        <div id="RegisterTitleWrapper">
          <img id="RegisterTitle" alt="logo" src={logo} />
        </div>

        <Form.Group className="RegisterFormGroup">
          <Form.Label className="RegisterFormLabel">
            USERNAME
          </Form.Label>
          <Form.Control
            className="RegisterInput"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </Form.Group>

        <Form.Group className="RegisterFormGroup">
          <Form.Label className="RegisterFormLabel">
            PASSWORD
          </Form.Label>
          <Form.Control
            className="RegisterInput"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Group>

        <div id="ButtonWrapper">
          <Button id="SubmitButton" value="GO!" type="submit">
            GO!
          </Button>
        </div>
      </Form>
    </div>
  );
};

export default Register;
