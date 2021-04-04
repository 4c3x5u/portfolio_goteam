import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Col } from 'react-bootstrap';

import columnName from './columnName';

import inbox from '../../assets/inbox.svg';
import ready from '../../assets/ready.svg';
import go from '../../assets/go.svg';
import done from '../../assets/done.svg';

const Column = ({ name }) => {
  const [title, setTitle] = useState('');

  useEffect(() => {
    switch (name) {
      case columnName.INBOX: setTitle(inbox); break;
      case columnName.READY: setTitle(ready); break;
      case columnName.GO: setTitle(go); break;
      case columnName.DONE: setTitle(done); break;
      default: throw TypeError('Title must be of type `columnTitle`');
    }
  }, []);

  const column = `Column ${
    name.charAt(0).toUpperCase() + name.slice(1)
  }Column`;

  return (
    <Col className="Col" xs={3}>
      <div className={column}>
        <div className="Header">
          <img src={title} alt={`${name} column header`} />
        </div>
        <div className="Body" />
      </div>
    </Col>
  );
};

Column.propTypes = {
  name: PropTypes.string.isRequired,
};

export default Column;
