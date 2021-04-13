import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import TeamControlsMenu from './Menu/TeamControlsMenu';

import './teamcontrolstoggler.sass';

const TeamControlsToggler = ({
  isActive, handleActivate, handleCreate, handleDelete, icon,
}) => (
  <Col xs={4} className="ControlsToggler">
    <button
      className="Button"
      onClick={handleActivate}
      aria-label="team controls toggler"
      type="button"
    >
      <FontAwesomeIcon icon={icon} />
      TEAM
    </button>

    {isActive
      && <TeamControlsMenu handleCreate={handleCreate} handleDelete={handleDelete} />}
  </Col>
);

TeamControlsToggler.propTypes = {
  isActive: PropTypes.bool.isRequired,
  handleActivate: PropTypes.func.isRequired,
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
  icon: PropTypes.string.isRequired,
};

export default TeamControlsToggler;
