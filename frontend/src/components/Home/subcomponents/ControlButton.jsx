import React from 'react';
import PropTypes from 'prop-types';
import { Col, Button } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const ControlButton = ({ name, action, icon }) => (
  <Col xs={4} className="ControlCol">
    <Button
      className="Button"
      onClick={action}
      aria-label={`${name} button`}
    >
      <FontAwesomeIcon icon={icon} />
      {name.toUpperCase()}
    </Button>
  </Col>
);

ControlButton.propTypes = {
  name: PropTypes.string.isRequired,
  action: PropTypes.func.isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  icon: PropTypes.object.isRequired,
};

export default ControlButton;
