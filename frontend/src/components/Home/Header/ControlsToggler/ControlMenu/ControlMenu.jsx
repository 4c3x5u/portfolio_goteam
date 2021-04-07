import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlusCircle, faCaretRight, faCaretLeft } from '@fortawesome/free-solid-svg-icons';

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

  return (
    <div className="ControlMenu">
      {items.map((item) => (
        <button key={item.name} type="button">
          {item.isActive
            && <FontAwesomeIcon className="IconLeft" icon={faCaretRight} />}

          {item.name}

          {item.isActive
            && <FontAwesomeIcon className="IconRight" icon={faCaretLeft} />}
        </button>
      ))}
      <button type="button" onClick={create}>
        <FontAwesomeIcon icon={faPlusCircle} />
      </button>
    </div>
  );
};

ControlMenu.propTypes = {
  create: PropTypes.func.isRequired,
};

export default ControlMenu;
