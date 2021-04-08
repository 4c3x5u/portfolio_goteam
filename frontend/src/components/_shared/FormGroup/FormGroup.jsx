import React from 'react';
import PropTypes from 'prop-types';
import { Form } from 'react-bootstrap';

import { inputType } from '../../../misc/inputType';

import './formgroup.sass';

const FormGroup = ({
  type, label, value, setValue,
}) => (
  <Form.Group className="FormGroup">
    <Form.Label className="Label">
      {label.toUpperCase()}
    </Form.Label>

    {type === inputType.TEXT && (
      <Form.Control
        className="Input"
        type={type}
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
    )}

    {type === inputType.TEXTAREA && (
      <Form.Control
        className="Input"
        as={type}
        rows={2}
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
    )}

  </Form.Group>
);

FormGroup.propTypes = {
  type: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  setValue: PropTypes.func.isRequired,
};

export default FormGroup;
