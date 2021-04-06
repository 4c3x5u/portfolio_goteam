import React from 'react';
import { Form } from 'react-bootstrap';
import PropTypes from 'prop-types';
import './formgroup.sass';

const FormGroup = ({ label, value, setValue }) => (
  <Form.Group className="FormGroup">
    <Form.Label className="Label">
      {label.toUpperCase()}
    </Form.Label>
    <Form.Control
      className="Input"
      type="text"
      value={value}
      onChange={(e) => setValue(e.target.value)}
    />
  </Form.Group>
);

FormGroup.propTypes = {
  label: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  setValue: PropTypes.func.isRequired,
};

export default FormGroup;
