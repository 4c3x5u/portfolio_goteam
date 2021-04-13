import React from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import BoardsControlMenu from './Menu/BoardsControlMenu';

import './boardscontrolstoggler.sass';

const BoardsControlsToggler = ({
  isActive, handleActivate, handleCreate, handleDelete, icon,
}) => (
  <Col xs={4} className="ControlsToggler">
    <button
      className="Button"
      onClick={handleActivate}
      aria-label="boards controls toggler"
      type="button"
    >
      <FontAwesomeIcon icon={icon} />

      BOARDS
    </button>

    {isActive && (
      <BoardsControlMenu
        handleCreate={handleCreate}
        handleDelete={handleDelete}
      />
    )}
  </Col>
);

BoardsControlsToggler.propTypes = {
  isActive: PropTypes.bool.isRequired,
  handleActivate: PropTypes.func.isRequired,
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
  icon: PropTypes.string.isRequired,
};

export default BoardsControlsToggler;
