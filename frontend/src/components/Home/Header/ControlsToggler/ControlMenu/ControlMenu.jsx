import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle } from '@fortawesome/free-solid-svg-icons';

import MenuItem from './MenuItem/MenuItem';

import './controlmenu.sass';

const ControlMenu = ({ handleCreate, handleDelete }) => {
  const [items, setItems] = useState([]);

  useEffect(() => (
    // TODO: API call here to get items
    setItems([
      { id: 0, name: 'An Item', isActive: false },
      { id: 1, name: 'An Active Item', isActive: true },
    ])
  ), []);

  const toggleItemActive = (item, index) => (
    // TODO: API call here to set items active
    setItems(items.map((currentItem, i) => (
      i === index
        ? { ...currentItem, isActive: !currentItem.isActive }
        : currentItem
    )))
  );

  return (
    <div className="ControlMenu">
      {items.map((item, index) => (
        <MenuItem
          id={item.id}
          name={item.name}
          isActive={item.isActive}
          toggleActive={() => toggleItemActive(item, index)}
          handleDelete={handleDelete}
        />
      ))}
      <button className="CreateButton" type="button" onClick={handleCreate}>
        <FontAwesomeIcon icon={faPlusCircle} />
      </button>
    </div>
  );
};

ControlMenu.propTypes = {
  handleCreate: PropTypes.func.isRequired,
  handleDelete: PropTypes.func.isRequired,
};

export default ControlMenu;
