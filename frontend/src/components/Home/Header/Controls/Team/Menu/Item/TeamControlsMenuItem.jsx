import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { useContextMenu, Item, Menu } from 'react-contexify';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCaretLeft, faCaretRight } from '@fortawesome/free-solid-svg-icons';

import AppContext from '../../../../../../../AppContext';
import UsersAPI from '../../../../../../../api/UsersAPI';

import './teamcontrolsmenuitem.sass';

const TeamControlsMenuItem = ({
  id, username, isActive, handleDelete,
}) => {
  const { activeBoard, loadBoard } = useContext(AppContext);

  const MENU_ID = `item-${id}`;
  const { show } = useContextMenu({ id: MENU_ID });

  const toggleActive = () => (
    UsersAPI
      .post(username, activeBoard.id, !isActive)
      .then(() => loadBoard())
      .catch((err) => console.error(err)) // TODO: Toast
  );

  return (
    <div className="MenuItem">
      <button
        className="ControlButton"
        key={id}
        type="button"
        onClick={toggleActive}
        onContextMenu={show}
      >
        {isActive && (
          <FontAwesomeIcon className="IconLeft" icon={faCaretRight} />
        )}

        {username}

        {isActive && (
          <FontAwesomeIcon className="IconRight" icon={faCaretLeft} />
        )}
      </button>

      <Menu className="ContextMenu" id={MENU_ID}>
        <Item
          className="ContextMenuItem"
          onClick={() => handleDelete({ id, name: username })}
        >
          <span>DELETE</span>
        </Item>
      </Menu>

    </div>
  );
};

TeamControlsMenuItem.propTypes = {
  id: PropTypes.number.isRequired,
  username: PropTypes.string.isRequired,
  isActive: PropTypes.bool.isRequired,
  handleDelete: PropTypes.bool.isRequired,
};

export default TeamControlsMenuItem;
