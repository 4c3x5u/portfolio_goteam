import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import ControlMenu from './ControlMenu/ControlMenu';

import './controlstoggler.sass';

const ControlsToggler = ({
  name, isActive, handleActivate, handleCreate, handleDelete, icon,
}) => (
  <Col xs={4} className="ControlsToggler">
    <button
      className="Button"
      onClick={handleActivate}
      aria-label={`${name} button`}
      type="button"
    >
      <FontAwesomeIcon icon={icon} />

      {name.toUpperCase()}
    </button>

    {isActive
      && <ControlMenu handleCreate={handleCreate} handleDelete={handleDelete} />}
  </Col>
);

ControlsToggler.propTypes = {
  name: PropTypes.string.isRequired,
  isActive: PropTypes.bool.isRequired,
  handleActivate: PropTypes.func.isRequired,
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
  icon: PropTypes.string.isRequired,
};

export default ControlsToggler;
