import React from 'react';
import PropTypes from 'prop-types';
import { useContextMenu, Item, Menu } from 'react-contexify';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCaretLeft, faCaretRight } from '@fortawesome/free-solid-svg-icons';

import './boardscontrolsmenuitem.sass';

const BoardsControlsMenuItem = ({
  id, name, isActive, toggleActive, handleDelete,
}) => {
  const MENU_ID = `item-${id}`;

  const { show } = useContextMenu({ id: MENU_ID });

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

        {name}

        {isActive && (
          <FontAwesomeIcon className="IconRight" icon={faCaretLeft} />
        )}
      </button>

      <Menu className="ContextMenu" id={MENU_ID}>
        <Item
          className="ContextMenuItem"
          onClick={() => handleDelete({ id, name })}
        >
          <span>DELETE</span>
        </Item>
      </Menu>

    </div>
  );
};

BoardsControlsMenuItem.propTypes = {
  id: PropTypes.number.isRequired,
  name: PropTypes.string.isRequired,
  isActive: PropTypes.bool.isRequired,
  toggleActive: PropTypes.func.isRequired,
  handleDelete: PropTypes.bool.isRequired,
};

export default BoardsControlsMenuItem;
