import React from 'react';
import PropTypes from 'prop-types';
import { Form, FormControl } from 'react-bootstrap';

import inputType from '../../../misc/inputType';

import './formgroup.sass';

const FormGroup = ({
  type, label, value, setValue, error, disabled,
}) => (
  <Form.Group className="FormGroup">
    <Form.Label className="Label">
      {label.toUpperCase()}
    </Form.Label>

    {type === inputType.TEXTAREA ? (
      <Form.Control
        className="Input"
        as={type}
        rows={2}
        value={value}
        onChange={(e) => setValue(e.target.value)}
        disabled={disabled}
        isInvalid={!!error}
      />
    ) : (
      <Form.Control
        className="Input"
        type={type}
        value={value}
        onChange={(e) => setValue(e.target.value)}
        disabled={disabled}
        isInvalid={!!error}
      />
    )}

    {error && (
      <FormControl.Feedback type="invalid">
        {error}
      </FormControl.Feedback>
    )}
  </Form.Group>
);

FormGroup.propTypes = {
  type: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  setValue: PropTypes.func,
  error: PropTypes.string,
  disabled: PropTypes.bool,
};

FormGroup.defaultProps = {
  setValue: () => DOMException.INVALID_ACCESS_ERR,
  disabled: false,
  error: null,
};

export default FormGroup;
