import React, { useEffect } from 'react';
import PropTypes from 'prop-types';
// import axios from 'axios';
import { Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

// import UserContext from '../../../../../UserContext';
import BoardsControlMenu from './Menu/BoardsControlMenu';

import './boardscontrols.sass';
// import ActiveBoardContext from '../../../ActiveBoardContext';

const BoardsControls = ({
  isActive, handleActivate, handleCreate, handleDelete, icon,
}) => {
  // const { currentUser } = useContext(UserContext);
  // const { activeBoardId, setActiveBoardId } = useContext(ActiveBoardContext);
  // const [boards, setBoards] = useState([{ id: null, name: '' }]);

  useEffect(() => {
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
          // boards={boards}
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
