import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faChalkboard } from '@fortawesome/free-solid-svg-icons';

import BoardsControlsMenu from './Menu/BoardsControlsMenu';
import AppContext from '../../../../../AppContext';

import './boardscontrols.sass';

const BoardsControls = ({
  isActive, handleActivate, handleCreate, handleEdit, handleDelete,
}) => {
  const { boards } = useContext(AppContext);

  return (
    <Col xs={4} className="BoardsControls">
      <button
        className="Button"
        onClick={handleActivate}
        aria-label="boards controls toggler"
        type="button"
        disabled={boards.length < 1}
      >
        <FontAwesomeIcon icon={faChalkboard} />

        BOARDS
      </button>

      {isActive && (
        <BoardsControlsMenu
          handleCreate={handleCreate}
          handleDelete={handleDelete}
          handleEdit={handleEdit}
        />
      )}
    </Col>
  );
};

BoardsControls.propTypes = {
  isActive: PropTypes.bool.isRequired,
  handleActivate: PropTypes.func.isRequired,
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
  handleEdit: PropTypes.func.isRequired,
};

export default BoardsControls;
