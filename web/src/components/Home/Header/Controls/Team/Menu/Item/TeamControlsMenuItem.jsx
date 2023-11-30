import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { useContextMenu, Item, Menu } from 'react-contexify';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faPlay,
  faChalkboardTeacher,
} from '@fortawesome/free-solid-svg-icons';

import AppContext from '../../../../../../../AppContext';
import UserAPI from '../../../../../../../api/UserAPI';

import './teamcontrolsmenuitem.sass';

const TeamControlsMenuItem = ({
  username, isAdmin, isActive, handleDelete,
}) => {
  const {
    user,
    activeBoard,
    loadBoard,
    notify,
  } = useContext(AppContext);

  const toggleActive = () => (
    UserAPI
      .patch(username, activeBoard.id, !isActive)
      .then(() => loadBoard())
      .catch((err) => {
        notify(
          'Unable to add member to the board.',
          `${err?.message || 'Server Error'}.`,
        );
      })
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

  const viewTooltip = (children) => (
    <OverlayTrigger
      placement="bottom"
      overlay={<Tooltip id="Username">{username}</Tooltip>}
    >
      {children}
    </OverlayTrigger>
  );

  const viewButton = (text) => (
    <button
      className="ControlButton"
      key={username}
      type="button"
      onClick={toggleActive}
      onContextMenu={(e) => (isAdmin ? e.preventDefault() : show(e))}
      disabled={!user.isAdmin}
    >
      {getIcon()}
      {text}
    </button>
  );

  return (
    <div className="TeamControlsMenuItem">
      {username.length <= 20
        ? viewButton(username)
        : viewTooltip(
          viewButton(`${username.substring(0, 17)}...`),
        )}

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
  handleDelete: PropTypes.func.isRequired,
};

export default TeamControlsMenuItem;
