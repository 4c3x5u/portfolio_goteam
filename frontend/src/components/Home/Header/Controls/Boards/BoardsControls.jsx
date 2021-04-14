import React, { useState, useEffect, useContext } from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import UserContext from '../../../../../UserContext';
import BoardsControlMenu from './Menu/BoardsControlMenu';

import './boardscontrols.sass';
import ActiveBoardContext from '../../../ActiveBoardContext';

const BoardsControls = ({
  isActive, handleActivate, handleCreate, handleDelete, icon,
}) => {
  const { currentUser } = useContext(UserContext);
  const { activeBoardId, setActiveBoardId } = useContext(ActiveBoardContext);
  const [boards, setBoards] = useState([{ id: null, name: '' }]);

  useEffect(() => {
    const endpoint = `${process.env.REACT_APP_BACKEND_URL}/boards/?team_id=`;

    axios.get(`${endpoint}${currentUser.teamId}`, {
      headers: {
        'auth-user': sessionStorage.getItem('username'),
        'auth-token': sessionStorage.getItem('auth-token'),
      },
    }).then((res) => {
      setBoards(
        res.data.boards.map((board) => ({
          id: board.id,
          name: board.name,
        })),
      );

      if (!activeBoardId) {
        setActiveBoardId(res.data.boards[0].id);
      }
    }).catch((err) => {
      // TODO: Handle properly
      console.log(`BoardsControls Error: ${err}`);
    });
  }, []);

  return (
    <Col xs={4} className="BoardsControls">
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
          boards={boards}
          handleCreate={handleCreate}
          handleDelete={handleDelete}
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
  icon: PropTypes.string.isRequired,
};

export default BoardsControls;
