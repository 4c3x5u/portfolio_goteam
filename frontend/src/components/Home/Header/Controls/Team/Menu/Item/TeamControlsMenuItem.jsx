import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { useContextMenu, Item, Menu } from 'react-contexify';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faPlay,
  faChalkboardTeacher,
} from '@fortawesome/free-solid-svg-icons';

import AppContext from '../../../../../../../AppContext';
import UsersAPI from '../../../../../../../api/UsersAPI';

import './teamcontrolsmenuitem.sass';

const TeamControlsMenuItem = ({
  username, isAdmin, isActive, handleDelete,
}) => {
  const { activeBoard, loadBoard } = useContext(AppContext);

  const toggleActive = () => (
    UsersAPI
      .post(username, activeBoard.id, !isActive)
      .then(() => loadBoard())
      .catch((err) => console.error(err)) // TODO: Toast
  );

  const MENU_ID = `item-${username}`;
  const { show } = useContextMenu({ id: MENU_ID });

  const getIcon = () => {
    if (isAdmin) {
      return <FontAwesomeIcon className="AdminIcon" icon={faChalkboardTeacher} />;
    }
    if (isActive) {
      return <FontAwesomeIcon className="ActiveIcon" icon={faPlay} />;
    }
    return <></>;
  };

  return (
    <div className="MenuItem">
      <button
        className="ControlButton"
        key={username}
        type="button"
        onClick={toggleActive}
        onContextMenu={show}
      >
        {getIcon()}
        {username}
      </button>

      <Menu className="ContextMenu" id={MENU_ID}>
        <Item
          className="ContextMenuItem"
          onClick={() => handleDelete({ username })}
        >
          <span>DELETE</span>
        </Item>
      </Menu>

    </div>
  );
};

TeamControlsMenuItem.propTypes = {
  username: PropTypes.string.isRequired,
  isAdmin: PropTypes.bool.isRequired,
  isActive: PropTypes.bool.isRequired,
  handleDelete: PropTypes.bool.isRequired,
};

export default TeamControlsMenuItem;
