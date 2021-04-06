import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import './controlstoggler.sass';

const ControlsToggler = ({
  name, isActive, activate, create, icon,
}) => (
  <Col xs={4} className="ControlsToggler">
    <button
      className="Button"
      onClick={activate}
      aria-label={`${name} button`}
      type="button"
    >
      <FontAwesomeIcon icon={icon} />
      {name.toUpperCase()}
    </button>

    {isActive && (
      <div className="Controls">
        <button type="button">
          Control Item #1
        </button>
        <button type="button">
          Control Item #2
        </button>
        <button type="button" onClick={create}>
          <FontAwesomeIcon icon={faPlusCircle} />
        </button>
      </div>
    )}
  </Col>
);

ControlsToggler.propTypes = {
  name: PropTypes.string.isRequired,
  isActive: PropTypes.bool.isRequired,
  activate: PropTypes.func.isRequired,
  create: PropTypes.func.isRequired,
  icon: PropTypes.string.isRequired,
};

export default ControlsToggler;
