import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { useContextMenu, Item, Menu } from 'react-contexify';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay } from '@fortawesome/free-solid-svg-icons';

import AppContext from '../../../../../../../AppContext';

import './boardscontrolsmenuitem.sass';

const BoardsControlsMenuItem = ({
  id, name, isActive, handleDelete, handleEdit,
}) => {
  const { user, loadBoard } = useContext(AppContext);

  const MENU_ID = `item-${id}`;
  const { show } = useContextMenu({ id: MENU_ID });

  const toggleActive = () => {
    loadBoard(id);
  };

  const viewTooltip = (children) => (
    <OverlayTrigger
      placement="bottom"
      overlay={<Tooltip id="BoardName">{name}</Tooltip>}
    >
      {children}
    </OverlayTrigger>
  );

  const viewButton = (text) => (
    <button
      className="ControlButton"
      key={id}
      type="button"
      onClick={toggleActive}
      onContextMenu={user.isAdmin && show}
    >
      {isActive
        && <FontAwesomeIcon className="ActiveIcon" icon={faPlay} />}
      {text}
    </button>
  );

  return (
    <div className="BoardsControlsMenuItem">
      {name.length <= 20
        ? viewButton(name)
        : viewTooltip(
          viewButton(`${name.substring(0, 17)}...`),
        )}

      <Menu className="ContextMenu" id={MENU_ID}>
        <Item
          className="ContextMenuItem"
          onClick={() => handleEdit({ id, name })}
        >
          <span>EDIT</span>
        </Item>
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
  handleDelete: PropTypes.func.isRequired,
  handleEdit: PropTypes.func.isRequired,
};

export default BoardsControlsMenuItem;
