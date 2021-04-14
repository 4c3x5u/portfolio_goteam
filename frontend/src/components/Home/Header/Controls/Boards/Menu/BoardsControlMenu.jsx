import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import ActiveBoardContext from '../../../../ActiveBoardContext';
import BoardsControlsMenuItem from './Item/BoardsControlsMenuItem';

import './boardscontrolsmenu.sass';

const BoardsControlMenu = ({ boards, handleCreate, handleDelete }) => {
  const { activeBoardId, setActiveBoardId } = useContext(ActiveBoardContext);

  return (
    <div className="ControlsMenu">
      {boards.map((board) => (
        <BoardsControlsMenuItem
          id={board.id}
          name={board.name}
          isActive={board.id === activeBoardId}
          toggleActive={() => setActiveBoardId(board.id)}
          handleDelete={handleDelete}
        />
      ))}

      <button className="CreateButton" type="button" onClick={handleCreate}>
        <FontAwesomeIcon icon={faPlusCircle} />
      </button>
    </div>
  );
};

BoardsControlMenu.propTypes = {
  boards: PropTypes.arrayOf({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
  }).isRequired,
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
};

export default BoardsControlMenu;
