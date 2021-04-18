import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import AppContext from '../../../../../../AppContext';
import BoardsControlsMenuItem from './Item/BoardsControlsMenuItem';

import './boardscontrolsmenu.sass';

const BoardsControlsMenu = ({ handleCreate, handleDelete, handleEdit }) => {
  const { boards, activeBoard, loadBoard } = useContext(AppContext);

  return (
    <div className="BoardsControlsMenu">
      {boards.length > 0 && boards.map((board) => (
        <BoardsControlsMenuItem
          id={board.id}
          name={board.name}
          isActive={board.id === activeBoard.id}
          toggleActive={() => loadBoard(board.id)}
          handleDelete={handleDelete}
          handleEdit={handleEdit}
        />
      ))}

      <button className="CreateButton" type="button" onClick={handleCreate}>
        <FontAwesomeIcon icon={faPlusCircle} />
      </button>
    </div>
  );
};

BoardsControlsMenu.propTypes = {
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
  handleEdit: PropTypes.func.isRequired,
};

export default BoardsControlsMenu;
