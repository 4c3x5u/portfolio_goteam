import React from 'react';
import PropTypes from 'prop-types';
import { Button } from 'react-bootstrap';

const ControlButton = ({ name, action }) => (
  <li>
    <Button
      className="Button"
      onClick={action}
      aria-label={`${name} button`}
    >
      {name.toUpperCase()}
    </Button>
  </li>
);

ControlButton.propTypes = {
  name: PropTypes.string.isRequired,
  action: PropTypes.func.isRequired,
};

export default ControlButton;
