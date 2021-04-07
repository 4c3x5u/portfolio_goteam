import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faPlusCircle, faCaretRight, faCaretLeft,
} from '@fortawesome/free-solid-svg-icons';

import './controlmenu.sass';

const ControlMenu = ({ create }) => {
  const [items, setItems] = useState([]);

  useEffect(() => (
    // TODO: Make an API call here to get items
    setItems([
      { name: 'An Item', isActive: false },
      { name: 'An Active Item', isActive: true },
    ])
  ), []);

  const toggleItemActive = (item, index) => (
    setItems(items.map((currentItem, i) => (
      i === index
        ? { ...currentItem, isActive: !currentItem.isActive }
        : currentItem
    )))
  );

  return (
    <div className="ControlMenu">
      {items.map((item, index) => (
        <button
          className="ControlButton"
          key={item.name}
          type="button"
          onClick={() => toggleItemActive(item, index)}
        >
          {item.isActive
            && <FontAwesomeIcon className="IconLeft" icon={faCaretRight} />}

          {item.name}

          {item.isActive
            && <FontAwesomeIcon className="IconRight" icon={faCaretLeft} />}
        </button>
      ))}
      <button className="CreateButton" type="button" onClick={create}>
        <FontAwesomeIcon icon={faPlusCircle} />
      </button>
    </div>
  );
};

ControlMenu.propTypes = {
  create: PropTypes.func.isRequired,
};

export default ControlMenu;
