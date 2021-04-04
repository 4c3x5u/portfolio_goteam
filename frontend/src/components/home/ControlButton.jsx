import React from 'react';
import PropTypes from 'prop-types';
import { Button } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const ControlButton = ({ name, action, icon }) => (
  <li>
    <Button
      className="Button"
      onClick={action}
      aria-label={`${name} button`}
    >
      <FontAwesomeIcon icon={icon} />
      {name.toUpperCase()}
    </Button>
  </li>
);

ControlButton.propTypes = {
  name: PropTypes.string.isRequired,
  action: PropTypes.func.isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  icon: PropTypes.object.isRequired,
};

export default ControlButton;
