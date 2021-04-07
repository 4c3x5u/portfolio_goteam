import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import ControlMenu from './ControlMenu/ControlMenu';

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

    {isActive
      && <ControlMenu create={create} />}
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
