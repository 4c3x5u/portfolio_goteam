import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import './register.sass';

const Register = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const handleSubmit = () => 1; // TODO: implement
  return (
    <div id="FormWrapper">
      <Form id="RegisterForm" onSubmit={handleSubmit}>
        <Form.Group controlId="username">
          <Form.Label>Username</Form.Label>
          <Form.Control
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </Form.Group>

        <Form.Group controlId="username">
          <Form.Label>Password</Form.Label>
          <Form.Control
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Group>

        <Button value="GO!" type="submit">GO!</Button>
      </Form>
    </div>
  );
};

export default Register;
